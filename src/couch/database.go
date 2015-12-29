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
// @uses    couch.util
// @author  Kerem Güneş <qeremy[at]gmail[dot]com>
package couch

import (
    "couch/util"
)

// Database object.
type Database struct {
    Client *Client
    Name   string
}

// DatabaseDocument object.
type DatabaseDocument struct {
    Id     string
    Key    string
    Value  map[string]interface{}
    Doc    map[string]interface{}
}

// DatabaseDocumentList object.
type DatabaseDocumentList struct {
    Offset     uint
    TotalRows  uint `json:"total_rows"`
    UpdateSeq  uint `json:"update_seq"`
    Rows       []DatabaseDocument
}

// Constructor.
//
// @param  client *couch.Client
// @param  client string
// @return *couch.Database
func NewDatabase(client *Client, name string) *Database {
    return &Database{
        Client: client,
          Name: name,
    }
}

// Ping database.
//
// @return bool
func (this *Database) Ping() bool {
    return (200 == this.Client.Head(this.Name, nil, nil).GetStatusCode())
}

// Get database info.
//
// @return map[string]interface{}, error
func (this *Database) Info() (map[string]interface{}, error) {
    data, err := this.Client.Get(this.Name, nil, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }

    var ret = util.Map()
    for key, value := range data.(map[string]interface{}) {
        ret[key] = value
    }

    return ret, nil
}

// Create database.
//
// @return bool
func (this *Database) Create() bool {
    return (201 == this.Client.Put(this.Name, nil, nil, nil).GetStatusCode())
}

// Remove database.
//
// @return bool
func (this *Database) Remove() bool {
    return (200 == this.Client.Delete(this.Name, nil, nil).GetStatusCode())
}


// Create database.
//
// @param  target       string
// @param  targetCreate bool
// @return bool
func (this *Database) Replicate(
    target string, targetCreate bool) (map[string]interface{}, error) {
    // prepare body
    var body = util.ParamList(
        "source", this.Name,
        "target", target,
        "create_target", targetCreate,
    )

    data, err := this.Client.Post("/_replicate", nil, body, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }

    var ret = util.Map()
    for key, value := range data.(map[string]interface{}) {
        // grap, set & pass history field
        if key == "history" {
            ret[key] = util.MapList(value)
            for i, history := range value.([]interface{}) {
                ret[key].([]map[string]interface{})[i] = util.Map()
                for kkey, vvalue := range history.(map[string]interface{}) {
                    ret[key].([]map[string]interface{})[i][kkey] = vvalue
                }
            }
            continue
        }
        ret[key] = value
    }

    return ret, nil
}

// Get document.
//
// @param  key string
// @return map[string]interface{}, error
func (this *Database) GetDocument(key string) (map[string]interface{}, error) {
    // prepare query
    var query = util.ParamList(
        "include_docs", true,
        "key"         , util.Quote(key),
    )

    data, err := this.Client.Get(this.Name +"/_all_docs", query, nil).
        GetBodyData(&DatabaseDocumentList{})
    if err != nil {
        return nil, err
    }

    var ret = util.Map()

    for _, doc := range data.(*DatabaseDocumentList).Rows {
        ret["id"]    = doc.Id
        ret["key"]   = doc.Key
        ret["value"] = map[string]string{"rev": doc.Value["rev"].(string)}
        ret["doc"]   = map[string]interface{}{}
        // fill doc field
        for key, value := range doc.Doc {
            ret["doc"].(map[string]interface{})[key] = value
        }
    }

    return ret, nil
}

// Get documents.
//
// @param  query map[string]interface{}
// @param  keys  []string
// @return map[string]interface{}, error
func (this *Database) GetDocumentAll(
    query map[string]interface{}, keys []string) (map[string]interface{}, error) {
    query = util.Param(query)
    if query["include_docs"] == nil {
        query["include_docs"] = true
    }

    // short?
    type ddl DatabaseDocumentList

    // make a reusable lambda
    var _func = func(data interface{}, err error) (map[string]interface{}, error) {
        if err != nil {
            return nil, err
        }

        var ret = util.Map()
        ret["offset"]     = data.(*ddl).Offset
        ret["total_rows"] = data.(*ddl).TotalRows

        var rows = data.(*ddl).Rows
        ret["rows"]       = util.MapList(len(rows))

        // append docs
        for i, row := range rows {
            ret["rows"].([]map[string]interface{})[i] = map[string]interface{}{
                   "id": row.Id,
                  "key": row.Key,
                "value": map[string]string{"rev": row.Value["rev"].(string)},
                  "doc": row.Doc,
            }
        }

        return ret, nil
    }

    if keys == nil {
        // get all
        return _func(
            this.Client.Get(this.Name +"/_all_docs", query, nil).GetBodyData(&ddl{}))
    } else {
        var body = util.ParamList("keys", keys)
        // get all only matched keys
        return _func(
            this.Client.Post(this.Name +"/_all_docs", query, body, nil).GetBodyData(&ddl{}))
    }
}

// Create document.
//
// @param  document map[string]interface{}
// @return map[string]interface{}, error
func (this *Database) CreateDocument(
    document interface{}) (map[string]interface{}, error) {
    data, err := this.CreateDocumentAll([]interface{}{document})
    if err != nil {
        return nil, err
    }

    if data := data[0]; data != nil {
        return data, nil
    }

    return nil, nil
}

// Create documents.
//
// @param  document []map[string]interface{}
// @return []map[string]interface{}, error
func (this *Database) CreateDocumentAll(
    documents []interface{}) ([]map[string]interface{}, error) {
    var docs = util.MapList(documents)
    for i, doc := range documents {
        if docs[i] == nil { docs[i] = util.Map() }
        // filter documents
        for key, value := range doc.(map[string]interface{}) {
            // this is create method, no update allowed
            if key == "_id" || key == "_id" || key == "_deleted" {
                continue
            }
            docs[i][key] = value
        }
    }

    var body = util.ParamList("docs", docs)
    data, err := this.Client.Post(this.Name +"/_bulk_docs", nil, body, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }

    var ret = util.MapList(data)
    for i, doc := range data.([]interface{}) {
        if ret[i] == nil { ret[i] = util.Map() }

        for key, value := range doc.(map[string]interface{}) {
            ret[i][key] = value
        }
    }

    return ret, nil
}

// Update document.
//
// @param  document map[string]interface{}
// @return map[string]interface{}, error
// @panics
func (this *Database) UpdateDocument(
    document interface{}) (map[string]interface{}, error) {
    data, err := this.UpdateDocumentAll([]interface{}{document})
    if err != nil {
        return nil, err
    }

    if data := data[0]; data != nil {
        return data, nil
    }

    return nil, nil
}

// Update documents.
//
// @param  document []map[string]interface{}
// @return []map[string]interface{}, error
// @panics
func (this *Database) UpdateDocumentAll(
    documents []interface{}) ([]map[string]interface{}, error) {
    var docs = util.MapList(documents)
    for i, doc := range documents {
        if docs[i] == nil {
            docs[i] = util.Map()
        }
        for key, value := range doc.(map[string]interface{}) {
            docs[i][key] = value
        }
        // these are required params
        if docs[i]["_id"] == nil || docs[i]["_rev"] == nil {
            panic("Both _id & _rev fields are required!")
        }
    }

    var body = util.ParamList("docs", docs)
    data, err := this.Client.Post(this.Name +"/_bulk_docs", nil, body, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }

    var ret = util.MapList(data)
    for i, doc := range data.([]interface{}) {
        if ret[i] == nil { ret[i] = util.Map() }

        for key, value := range doc.(map[string]interface{}) {
            ret[i][key] = value
        }
    }

    return ret, nil
}

// Delete documents.
//
// @param  document map[string]interface{}
// @return map[string]interface{}, error
// @panics
func (this *Database) DeleteDocument(
    document interface{}) (map[string]interface{}, error) {
    data, err := this.DeleteDocumentAll([]interface{}{document})
    if err != nil {
        return nil, err
    }

    if data := data[0]; data != nil {
        return data, nil
    }

    return nil, nil
}

// Delete documents.
//
// @param  document []map[string]interface{}
// @return []map[string]interface{}, error
// @panics
func (this *Database) DeleteDocumentAll(documents []interface{}) (
        []map[string]interface{}, error) {
    for i, _ := range documents {
        // just add "_deleted" param into document
        documents[i].(map[string]interface{})["_deleted"] = true
    }

    return this.UpdateDocumentAll(documents)
}

// Get changes.
//
// @param  query map[string]interface{}
// @return []map[string]interface{}, error
func (this *Database) GetChanges(
    query map[string]interface{}, docIds []string) (map[string]interface{}, error) {
    query = util.Param(query)
    if docIds != nil {
        query["filter"] = "_doc_ids"
    }

    var body = util.ParamList("doc_ids", docIds)
    data, err := this.Client.Post(this.Name +"/_changes", query, body, nil).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }

    var ret = util.Map()
    ret["last_seq"] = util.Dig("last_seq", data)
    ret["results"]  = util.MapList(0) // set empty as default
    if results := data.(map[string]interface{})["results"].([]interface{});
       results != nil {
        ret["results"] = util.MapList(results) // @overwrite
        for i, result := range results {
            ret["results"].([]map[string]interface{})[i] = map[string]interface{}{
                     "id": util.Dig("id", result),
                    "seq": util.Dig("seq", result),
                "deleted": util.Dig("deleted", result),
                "changes": util.Dig("changes", result),
            }
        }
    }

    return ret, nil
}

// Compact.
//
// @param  ddoc string
// @return bool, error
func (this *Database) Compact(ddoc string) (bool, error) {
    data, err := this.Client.Post(this.Name +"/_compact/"+ ddoc, nil, nil, nil).GetBodyData(nil)
    if err != nil {
        return false, err
    }

    return util.DigBool("ok", data), nil
}

func (this *Database) EnsureFullCommit() (bool, uint, error) {
    data, err := this.Client.Post(this.Name +"/_ensure_full_commit", nil, nil, nil).GetBodyData(nil)
    if err != nil {
        return false, 0, err
    }

    return util.DigBool("ok", data),
           util.DigUInt("instance_start_time", data),
           nil
}

func (this *Database) ViewCleanup() (bool, error) {
    data, err := this.Client.Post(this.Name +"/_view_cleanup", nil, nil, nil).
        GetBodyData(nil)
    if err != nil {
        return false, err
    }
    return util.DigBool("ok", data), nil
}

func (this *Database) ViewTemp(_map string, reduce interface{}) (
        map[string]interface{}, error) {
    var body = util.ParamList(
        "map", _map,
        // prevent "missing function" error
        "reduce", util.IsEmptySet(reduce, nil),
    )
    data, err := this.Client.Post(this.Name +"/_temp_view", nil, body, nil).
        GetBodyData(&DatabaseDocumentList{})
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    var _returnRows = data.(*DatabaseDocumentList).Rows
    _return["offset"]     = data.(*DatabaseDocumentList).Offset
    _return["total_rows"] = data.(*DatabaseDocumentList).TotalRows
    _return["rows"]       = util.MapList(len(_returnRows))
    for i, row := range _returnRows {
        _return["rows"].([]map[string]interface{})[i] = map[string]interface{}{
               "id": row.Id,
              "key": row.Key,
            "value": row.Value,
        }
    }
    return _return, nil
}

func (this *Database) GetSecurity() (map[string]interface{}, error) {
    data, err := this.Client.Get(this.Name +"/_security", nil, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    return data.(map[string]interface{}), nil
}

func (this *Database) SetSecurity(admins, members map[string]interface{}) (bool, error) {
    if admins["names"].([]string) == nil || admins["roles"].([]string)  == nil ||
       members["names"].([]string) == nil || members["roles"].([]string) == nil {
        panic("Specify admins and/or members with names=>roles fields!")
    }
    var body = util.ParamList("admins", admins, "members", members)
    data, err := this.Client.Put(this.Name +"/_security", nil, body, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    return util.DigBool("ok", data), nil
}

func (this *Database) Purge(object map[string]interface{}) (map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_purge", nil, object, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    _return["purge_seq"] = util.DigInt("purge_seq", data)
    _return["purged"]  = util.Map()
    for id, revs := range data.(map[string]interface{})["purged"].
        (map[string]interface{}) {
        _return["purged"].(map[string]interface{})[id] = revs
    }
    return _return, nil
}

func (this *Database) GetMissingRevisions(object map[string]interface{}) (
        map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_missing_revs", nil, object, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    _return["missing_revs"] = util.Map()
    for id, revs := range data.(map[string]interface{})["missing_revs"].
        (map[string]interface{}) {
        _return["missing_revs"].(map[string]interface{})[id] = revs
    }
    return _return, nil
}

func (this *Database) GetMissingRevisionsDiff(object map[string]interface{}) (
        map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_revs_diff", nil, object, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    for id, _ := range data.(map[string]interface{}) {
        _return[id] = map[string]interface{}{
            "missing": util.Dig(id +".missing", data),
        }
    }
    return _return, nil
}

func (this *Database) GetRevisionLimit() (int, error) {
    data, err := this.Client.Get(this.Name +"/_revs_limit", nil, nil).GetBodyData(nil)
    if err != nil {
        return -1, err
    }
    return int(data.(float64)), nil
}

func (this *Database) SetRevisionLimit(limit int) (bool, error) {
    data, err := this.Client.Put(this.Name +"/_revs_limit", nil, limit, nil).GetBodyData(limit)
    if err != nil {
        return false, err
    }
    return util.DigBool("ok", data), nil
}
