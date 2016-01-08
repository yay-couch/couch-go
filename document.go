// Copyright 2015 Kerem Güneş
//    <http://qeremy.com>
//
// Apache License, Version 2.0
//    <http://www.apache.org/licenses/LICENSE-2.0>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// @package couch
// @uses    couch.util, couch.uuid
// @author  Kerem Güneş <qeremy[at]gmail[dot]com>
package couch

import (
    "couch/util"
    "couch/uuid"
)

// @object couch.Document
type Document struct {
    Id          *uuid.Uuid
    Rev         string
    Deleted     bool
    Attachments map[string]*DocumentAttachment
    Data        map[string]interface{}
    Database    *Database
}

// Constructor.
//
// @param  database *couch.Database
// @param  data... interface{}
// @return (*couch.Document)
func NewDocument(database *Database, data... interface{}) (*Document) {
    var this = &Document{
        Data : util.Map(),
        Database: database,
    }

    if data != nil {
        this.SetData(util.ParamList(data...))
    }

    return this
}

// Set database.
//
// @param  database *couch.Database
// @return (void)
func (this *Document) SetDatabase(database *Database) {
    this.Database = database
}

// Get database
//
// @return (database *couch.Database)
func (this *Document) GetDatabase() (*Database) {
    return this.Database
}

// Setter.
//
// @param  data... interface{}
// @return (*couch.Document)
// @panics
func (this *Document) Set(data... interface{}) (*Document) {
    if data == nil {
        panic("Provide at least a key=>value match as param!")
    }
    this.SetData(util.ParamList(data...))

    return this
}

// Set ID.
//
// @param id interface{}
func (this *Document) SetId(id interface{}) {
    if _, ok := id.(*uuid.Uuid); !ok {
        id = uuid.New(id)
    }
    this.Id = id.(*uuid.Uuid)
}

// Set rev.
//
// @param  rev string
// @return (void)
func (this *Document) SetRev(rev string) {
    this.Rev = rev
}

// Set deleted.
//
// @param  deleted bool
// @return (void)
func (this *Document) SetDeleted(deleted bool) {
    this.Deleted = deleted
}

// Set attachment.
//
// @param  attachment interface{}
// @return (void)
// @panics
func (this *Document) SetAttachment(attachment interface{}) {
    if _, ok := attachment.(*DocumentAttachment); !ok {
        var file = util.DigString("file", attachment)
        var fileName = util.DigString("fileName", attachment)
        attachment = NewDocumentAttachment(this, file, fileName)
    }

    // file name must be uniq
    if _, ok := this.Attachments[attachment.(*DocumentAttachment).FileName]; ok {
        panic("Attachment is alredy exists on this document!")
    }

    if this.Attachments == nil {
        this.Attachments = make(map[string]*DocumentAttachment)
    }

    this.Attachments[attachment.(*DocumentAttachment).FileName] =
        attachment.(*DocumentAttachment);
}

// Set data.
//
// @param  data map[string]interface{}
// @return (void)
func (this *Document) SetData(data map[string]interface{}) {
    for key, value := range data {
        // set special properties
        if key == "_id"      { this.SetId(value) }
        if key == "_rev"     { this.SetRev(value.(string)) }
        if key == "_deleted" { this.SetDeleted(value.(bool)) }

        // set "_attachments" and pass
        if key == "_attachments" {
            for _, attachment := range value.([]interface{}) {
                this.SetAttachment(attachment)
            }
            continue
        }

        this.Data[key] = value
    }
}

// Getter.
//
// @param  key string
// @return (interface{})
func (this *Document) Get(key string) (interface{}) {
    if value, ok := this.Data[key]; ok {
        return value
    }

    return nil
}

// Get ID.
//
// @return (string)
func (this *Document) GetId() (string) {
    if this.Id != nil {
        return this.Id.ToString()
    }
    return ""
}

// Get rev.
//
// @return (string)
func (this *Document) GetRev() (string) {
    return this.Rev
}

// Get deleted.
//
// @return (bool)
func (this *Document) GetDeleted() (bool) {
    return this.Deleted
}

// Get attachment.
//
// @return (interface{})
func (this *Document) GetAttachment(fileName string) (interface{}) {
    if attachment, ok := this.Attachments[fileName]; ok {
        return attachment
    }
    return nil
}

// Get data.
//
// @return (map[string]interface{})
func (this *Document) GetData() (map[string]interface{}) {
    return this.Data
}

// Ping.
//
// @param  statusCode uint16
// @return (bool)
// @panics
func (this *Document) Ping(statusCode uint16) (bool) {
    var id = this.GetId()
    if id == "" {
        panic("_id field is could not be empty!")
    }

    var headers = util.Map()
    if (this.Rev != "") {
        headers["If-None-Match"] = util.Quote(this.Rev);
    }

    return (statusCode == this.Database.Client.
        Head(this.Database.Name +"/"+ util.UrlEncode(id), nil, headers).GetStatusCode())
}

// Check is exists.
//
// @return (bool)
// @panics
func (this *Document) IsExists() (bool) {
    var id = this.GetId()
    if id == "" {
        panic("_id field is could not be empty!")
    }

    var headers = util.Map()
    if (this.Rev != "") {
        headers["If-None-Match"] = util.Quote(this.Rev);
    }

    var statusCode = this.Database.Client.
        Head(this.Database.Name +"/"+ util.UrlEncode(id), nil, headers).GetStatusCode()

    return (statusCode == 200 || statusCode == 304)
}

// Check is not modified.
//
// @return (bool)
// @panics
func (this *Document) IsNotModified() (bool) {
    var id, rev = this.GetId(), this.GetRev()
    if id == "" || rev == "" {
        panic("_id & _rev fields are could not be empty!")
    }

    var headers = util.Map()
    headers["If-None-Match"] = util.Quote(rev);

    return (304 == this.Database.Client.
        Head(this.Database.Name +"/"+ util.UrlEncode(id), nil, headers).GetStatusCode())
}

// Find.
//
// @param  query map[string]interface{}
// @return (map[string]interface{}, error)
// @panics
func (this *Document) Find(query map[string]interface{}) (map[string]interface{}, error) {
    var id = this.GetId()
    if id == "" {
        panic("_id field is could not be empty!")
    }

    query = util.Param(query)
    if query["rev"] == "" && this.Rev != "" {
        query["rev"] = this.Rev
    }

    data, err := this.Database.Client.Get(this.Database.Name +"/"+ util.UrlEncode(id), query, nil).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }

    var ret = util.Map()
    for key, value := range data.(map[string]interface{}) {
        ret[key] = value
    }

    return ret, nil
}

// Find as struct.
//
// @param  data interface{}
// @param  query map[string]interface{}
// @return (interface{}, error)
func (this *Document) FindStruct(
    data interface{}, query map[string]interface{}) (interface{}, error) {
    var id = this.GetId()
    if id == "" {
        panic("_id field is could not be empty!")
    }
    if data == nil {
        panic("You should pass your data struct!")
    }

    query = util.Param(query)
    if query["rev"] == "" && this.Rev != "" {
        query["rev"] = this.Rev
    }

    data, err := this.Database.Client.Get(this.Database.Name +"/"+ util.UrlEncode(id), query, nil).
        GetBodyData(data)
    if err != nil {
        return nil, err
    }

    return data, nil
}

// Find revisions.
//
// @return (map[string]interface{}, error)
func (this *Document) FindRevisions() (map[string]interface{}, error) {
    data, err := this.Find(util.ParamList("revs", true))
    if err != nil {
        return nil, err
    }

    var ret = util.Map()
    if data["_revisions"] != nil {
        ret["start"] = util.DigInt("_revisions.start", data)
        ret["ids"]   = util.DigSliceString("_revisions.ids", data)
    }

    return ret, nil
}

// Find revisions extended.
//
// @return ([]map[string]interface{}, error)
func (this *Document) FindRevisionsExtended() ([]map[string]string, error) {
    data, err := this.Find(util.ParamList("revs_info", true))
    if err != nil {
        return nil, err
    }

    var ret = util.MapListString(nil)
    if data["_revs_info"] != nil {
        ret = util.MapListString(data["_revs_info"]) // @overwrite
        for i, info := range data["_revs_info"].([]interface{}) {
            ret[i] = map[string]string{
                   "rev": util.DigString("rev", info),
                "status": util.DigString("status", info),
            }
        }
    }

    return ret, nil
}

// Find attachments
//
// @param  attEncInfo bool
// @param  attsSince []string
// return  []map[string]interface{}, error
func (this *Document) FindAttachments(
    attEncInfo bool, attsSince []string) ([]map[string]interface{}, error) {
    var query = util.Map()
    query["attachments"] = true
    query["att_encoding_info"] = attEncInfo

    if attsSince != nil {
        var attsSinceArray = util.MapSliceString(attsSince)
        for _, attsSinceValue := range attsSince {
            attsSinceArray = append(attsSinceArray, util.QuoteEncode(attsSinceValue))
        }
    }

    data, err := this.Find(query)
    if err != nil {
        return nil, err
    }

    var ret = util.MapList(nil)
    if data["_attachments"] != nil {
        for _, attc := range data["_attachments"].(map[string]interface{}) {
            ret = append(ret, attc.(map[string]interface{}))
        }
    }

    return ret, nil
}

// Save.
//
// @param  args... bool
// @return (map[string]interface{}, error)
func (this *Document) Save(args... bool) (map[string]interface{}, error) {
    var query, headers, body = util.Map(), util.Map(), this.GetData()
    if args != nil {
        if args[0] == true {
            query["batch"] = "ok"
        }
        if args[1] == true {
            headers["X-Couch-Full-Commit"] = "true"
        }
    }

    if this.Rev != "" {
        headers["If-Match"] = this.Rev
    }

    if this.Attachments != nil {
        body["_attachments"] = util.Map()
        for name, attachment := range this.Attachments {
            body["_attachments"].(map[string]interface{})[name] = attachment.ToArray(true)
        }
    }

    // make a reusable lambda
    var _func = func(data interface{}, err error) (map[string]interface{}, error) {
        if err != nil {
            return nil, err
        }

        var id, rev = util.DigString("id", data),
                      util.DigString("rev", data)
        // set id & rev for next save() instant calls
        if id != "" && this.Id == nil {
            this.SetId(id)
        }
        if rev != "" {
            this.SetRev(rev)
        }

        return map[string]interface{}{
             "ok": util.DigBool("ok", data),
             "id": id,
            "rev": rev,
        }, nil
    }

    if (this.Id == nil) {
        return _func( // insert action
            this.Database.Client.Post(this.Database.Name, query, body, headers).
                GetBodyData(nil))
    } else {
        return _func( // update action
            this.Database.Client.Put(this.Database.Name +"/"+ this.GetId(), query, body, headers).
                GetBodyData(nil))
    }
}

// Remove.
//
// @param  args... bool
// @return (map[string]interface{}, error)
// @panics
func (this *Document) Remove(args... bool) (map[string]interface{}, error) {
    var id, rev = this.GetId(), this.GetRev()
    if id == "" || rev == "" {
        panic("Both _id & _rev fields could not be empty!")
    }

    var query, headers = util.Map(), util.Map()
    headers["If-Match"] = rev

    if args != nil {
        if args[0] == true {
            query["batch"] = "ok"
        }
        if args[1] == true {
            headers["X-Couch-Full-Commit"] = "true"
        }
    }

    data, err := this.Database.Client.Delete(this.Database.Name +"/"+ util.UrlEncode(id), query, headers).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
         "ok": util.DigBool("ok", data),
         "id": util.DigString("id", data),
        "rev": util.DigString("rev", data),
    }, nil
}

// Copy.
//
// @param  dest string
// @param  args... bool
// @return (map[string]interface{}, error)
// @panics
func (this *Document) Copy(dest string, args... bool) (map[string]interface{}, error) {
    var id = this.GetId()
    if id == "" {
        panic("_id field could not be empty!");
    }
    if dest == "" {
        panic("Destination could not be empty!");
    }

    var query, headers = util.Map(), util.Map()
    headers["Destination"] = dest

    if args != nil {
        if args[0] == true {
            query["batch"] = "ok"
        }
        if args[1] == true {
            headers["X-Couch-Full-Commit"] = "true"
        }
    }
    data, err := this.Database.Client.Copy(this.Database.Name +"/"+ util.UrlEncode(id), query, headers).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
         "ok": util.DigBool("ok", data),
         "id": util.DigString("id", data),
        "rev": util.DigString("rev", data),
    }, nil
}

// Copy from.
//
// @param  dest string
// @param  args... bool
// @return (map[string]interface{}, error)
// @panics
func (this *Document) CopyFrom(dest string, args... bool) (map[string]interface{}, error) {
    var id, rev = this.GetId(), this.GetRev()
    if id == "" || rev == "" {
        panic("Both _id & _rev fields could not be empty!");
    }
    if dest == "" {
        panic("Destination could not be empty!");
    }

    var query, headers = util.Map(), util.Map()
    headers["If-Match"] = rev
    headers["Destination"] = dest

    if args != nil {
        if args[0] == true {
            query["batch"] = "ok"
        }
        if args[1] == true {
            headers["X-Couch-Full-Commit"] = "true"
        }
    }

    data, err := this.Database.Client.Copy(this.Database.Name +"/"+ util.UrlEncode(id), query, headers).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
         "ok": util.DigBool("ok", data),
         "id": util.DigString("id", data),
        "rev": util.DigString("rev", data),
    }, nil
}

// Copy to.
//
// @param  dest    string
// @param  destRev string
// @param  args... bool
// @return (map[string]interface{}, error)
// @panics
func (this *Document) CopyTo(dest, destRev string, args... bool) (map[string]interface{}, error) {
    var id, rev = this.GetId(), this.GetRev()
    if id == "" || rev == "" {
        panic("Both _id & _rev fields could not be empty!");
    }
    if dest == "" || destRev == "" {
        panic("Destination & destination revision could not be empty!");
    }

    var query, headers = util.Map(), util.Map()
    headers["If-Match"] = rev
    headers["Destination"] = util.StringFormat("%s?rev=%s", dest, destRev);

    if args != nil {
        if args[0] == true {
            query["batch"] = "ok"
        }
        if args[1] == true {
            headers["X-Couch-Full-Commit"] = "true"
        }
    }

    data, err := this.Database.Client.Copy(this.Database.Name +"/"+ util.UrlEncode(id), query, headers).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
         "ok": util.DigBool("ok", data),
         "id": util.DigString("id", data),
        "rev": util.DigString("rev", data),
    }, nil
}
