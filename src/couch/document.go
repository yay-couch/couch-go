package couch

import (
    "couch/util"
    "couch/uuid"
)

type Document struct {
    Id          *uuid.Uuid
    Rev         string
    Deleted     bool
    Attachments map[string]*DocumentAttachment
    Data        map[string]interface{}
    Database    *Database
}

func NewDocument(database *Database, data ...interface{}) *Document {
    var this = &Document{
        Database: database,
    }
    if data != nil {
        this.SetData(util.ParamList(data...))
    }
    return this
}

func (this *Document) SetDatabase(database *Database) {
    this.Database = database
}
func (this *Document) GetDatabase() *Database {
    return this.Database
}

func (this *Document) Set(data... interface{}) *Document {
    if data == nil {
        panic("Provide at least a key=>value match as param!")
    }
    this.SetData(util.ParamList(data...))
    return this
}
func (this *Document) SetId(id interface{}) {
    if _, ok := id.(*uuid.Uuid); !ok {
        id = uuid.New(id)
    }
    this.Id = id.(*uuid.Uuid)
}
func (this *Document) SetRev(rev string) {
    this.Rev = rev
}
func (this *Document) SetDeleted(deleted bool) {
    this.Deleted = deleted
}
func (this *Document) SetData(data map[string]interface{}) {
    if this.Data == nil {
        this.Data = util.Map()
    }
    for key, value := range data {
        if key == "_id"      { this.SetId(value) }
        if key == "_rev"     { this.SetRev(value.(string)) }
        if key == "_deleted" { this.SetDeleted(value.(bool)) }
        if key == "_attachments" {
            for _, attachment := range value.([]interface{}) {
                this.SetAttachment(attachment)
            }
            continue
        }
        this.Data[key] = value
    }
}
func (this *Document) SetAttachment(attachment interface{}) {
    if _, ok := attachment.(*DocumentAttachment); !ok {
        var file = util.DigString("file", attachment)
        var fileName = util.DigString("fileName", attachment)
        attachment = NewDocumentAttachment(this, file, fileName)
    }
    if _, ok := this.Attachments[attachment.(*DocumentAttachment).FileName]; ok {
        panic("Attachment is alredy exists on this document!")
    }
    if this.Attachments == nil {
        this.Attachments = make(map[string]*DocumentAttachment)
    }
    this.Attachments[attachment.(*DocumentAttachment).FileName] = attachment.(*DocumentAttachment);
}

func (this *Document) Get(key string) interface{} {
    if value, ok := this.Data[key]; ok {
        return value
    }
    return nil
}
func (this *Document) GetId() string {
    if this.Id != nil {
        return this.Id.ToString()
    }
    return ""
}
func (this *Document) GetRev() string {
    return this.Rev
}
func (this *Document) GetDeleted() bool {
    return this.Deleted
}
func (this *Document) GetData() map[string]interface{} {
    return this.Data
}
func (this *Document) GetDataValue(key string) interface{} {
    return util.Dig(key, this.Data)
}

func (this *Document) Ping(statusCode uint16) bool {
    if this.Id == nil {
        panic("_id field is could not be empty!")
    }
    var headers = util.Map()
    if (this.Rev != "") {
        headers["If-None-Match"] = util.Quote(this.Rev);
    }
    return (statusCode == this.Database.Client.
        Head(this.Database.Name +"/"+ this.GetId(), nil, headers).GetStatusCode())
}
func (this *Document) IsExists() bool {
    if this.Id == nil {
        panic("_id field is could not be empty!")
    }
    var headers = util.Map()
    if (this.Rev != "") {
        headers["If-None-Match"] = util.Quote(this.Rev);
    }
    var statusCode = this.Database.Client.
        Head(this.Database.Name +"/"+ this.GetId(), nil, headers).GetStatusCode()
    return (statusCode == 200 || statusCode == 304)
}
func (this *Document) IsNotModified() bool {
    if this.Id == nil || this.Rev == "" {
        panic("_id & _rev fields are could not be empty!")
    }
    var headers = util.Map()
    headers["If-None-Match"] = util.Quote(this.Rev);
    return (304 == this.Database.Client.
        Head(this.Database.Name +"/"+ this.GetId(), nil, headers).GetStatusCode())
}

func (this *Document) Find(query map[string]interface{}) (map[string]interface{}, error) {
    if this.Id == nil {
        panic("_id field is could not be empty!")
    }
    query = util.Param(query)
    if query["rev"] == "" && this.Rev != "" {
        query["rev"] = this.Rev
    }
    data, err := this.Database.Client.Get(this.Database.Name +"/"+ this.GetId(), query, nil).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    for key, value := range data.(map[string]interface{}) {
        _return[key] = value
    }
    return _return, nil
}
func (this *Document) TestFindStruct(data interface{}, query map[string]interface{}) (interface{}, error) {
    if data == nil {
        panic("You should pass your data struct!")
    }
    data, err := this.Database.Client.Get(this.Database.Name +"/"+ this.GetId(), query, nil).
        GetBodyData(data)
    if err != nil {
        return nil, err
    }
    return data, nil
}
func (this *Document) FindRevisions() (map[string]interface{}, error) {
    data, err := this.Find(util.ParamList("revs", true))
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    if data["_revisions"] != nil {
        _return["start"] = util.DigInt("_revisions.start", data)
        _return["ids"]   = util.DigSliceString("_revisions.ids", data)
    }
    return _return, nil
}
func (this *Document) FindRevisionsExtended() ([]map[string]string, error) {
    data, err := this.Find(util.ParamList("revs_info", true))
    if err != nil {
        return nil, err
    }
    var _return = util.MapListString(nil)
    if data["_revs_info"] != nil {
        // @overwrite
        _return = util.MapListString(data["_revs_info"])
        for i, info := range data["_revs_info"].([]interface{}) {
            _return[i] = map[string]string{
                   "rev": util.DigString("rev", info),
                "status": util.DigString("status", info),
            }
        }
    }
    return _return, nil
}
func (this *Document) FindAttachments(attEncInfo bool, attsSince []string) ([]map[string]interface{}, error) {
    var query = util.Param(nil)
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
    var _return = util.MapList(nil)
    if data["_attachments"] != nil {
        for _, attc := range data["_attachments"].(map[string]interface{}) {
            _return = append(_return, attc.(map[string]interface{}))
        }
    }
    return _return, nil
}

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
    var _return = func(data interface{}, err error) (map[string]interface{}, error) {
        if err != nil {
            return nil, err
        }
        var id, rev = util.DigString("id", data), util.DigString("rev", data)
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
        return _return( // insert action
            this.Database.Client.Post(this.Database.Name, query, body, headers).
                GetBodyData(nil))
    } else {
        return _return( // update action
            this.Database.Client.Put(this.Database.Name +"/"+ this.GetId(), query, body, headers).
                GetBodyData(nil))
    }
}

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
    data, err := this.Database.Client.Delete(this.Database.Name +"/"+ id, query, headers).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
         "ok": util.DigBool("ok", data),
         "id": util.DigString("id", data),
        "rev": util.DigString("rev", data),
    }, nil
}

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
    data, err := this.Database.Client.Copy(this.Database.Name +"/"+ id, query, headers).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
         "ok": util.DigBool("ok", data),
         "id": util.DigString("id", data),
        "rev": util.DigString("rev", data),
    }, nil
}

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
    data, err := this.Database.Client.Copy(this.Database.Name +"/"+ id, query, headers).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
         "ok": util.DigBool("ok", data),
         "id": util.DigString("id", data),
        "rev": util.DigString("rev", data),
    }, nil
}

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
    data, err := this.Database.Client.Copy(this.Database.Name +"/"+ id, query, headers).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
         "ok": util.DigBool("ok", data),
         "id": util.DigString("id", data),
        "rev": util.DigString("rev", data),
    }, nil
}
