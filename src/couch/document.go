package couch

import (
    "./util"
    "./uuid"
)

type Document struct {
    Id          *uuid.Uuid
    Rev         string
    Deleted     bool
    Attachments []DocumentAttachment
    Data        map[string]interface{}
    Database    *Database
}

func NewDocument(database *Database, data map[string]interface{}) *Document {
    var this = &Document{
        Database: database,
    }
    if data != nil {
        this.SetData(data)
    }
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
            // @todo
        }
        this.Data[key] = value
    }
}

func (this *Document) GetId() string {
    return this.Id.ToString()
}
func (this *Document) GetRev() string {
    return this.Rev
}
func (this *Document) GetDeleted() bool {
    return this.Deleted
}
func (this *Document) GetData(key interface{}) interface{} {
    if key != nil {
        return util.Dig(key.(string), this.Data)
    }
    return this.Data
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
        Head(this.Database.Name +"/"+ this.Id.ToString(), nil, headers).GetStatusCode())
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
        Head(this.Database.Name +"/"+ this.Id.ToString(), nil, headers).GetStatusCode()
    return (statusCode == 200 || statusCode == 304)
}
func (this *Document) IsNotModified() bool {
    if this.Id == nil || this.Rev == "" {
        panic("_id & _rev fields are could not be empty!")
    }
    var headers = util.Map()
    headers["If-None-Match"] = util.Quote(this.Rev);
    return (304 == this.Database.Client.
        Head(this.Database.Name +"/"+ this.Id.ToString(), nil, headers).GetStatusCode())
}
func (this *Document) Find(query map[string]interface{}) (map[string]interface{}, error) {
    if this.Id == nil {
        panic("_id field is could not be empty!")
    }
    query = util.Param(query)
    if query["rev"] == "" && this.Rev != "" {
        query["rev"] = this.Rev
    }
    data, err := this.Database.Client.Get(
        this.Database.Name +"/"+ this.Id.ToString(), query, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    for key, value := range data.(map[string]interface{}) {
        _return[key] = value
    }
    return _return, nil
}
func (this *Document) FindRevisions() (map[string]interface{}, error) {
    data, err := this.Find(util.ParamList("revs", true))
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    if data["_revisions"] != nil {
        _return["ids"] = util.MapStringSlice(nil)
        for _, id := range data["_revisions"].(map[string]interface{})["ids"].([]interface{}) {
            _return["ids"] = append(_return["ids"].([]string), id.(string))
        }
        _return["start"] = uint(data["_revisions"].(map[string]interface{})["start"].(float64))
    }
    return _return, nil
}
func (this *Document) FindRevisionsExtended() ([]map[string]string, error) {
    data, err := this.Find(util.ParamList("revs_info", true))
    if err != nil {
        return nil, err
    }
    var _return = util.MapStringArray(nil)
    if data["_revs_info"] != nil {
        // @overwrite
        _return = util.MapStringArray(data["_revs_info"])
        for i, info := range data["_revs_info"].([]interface{}) {
            _return[i] = map[string]string{
                "rev": util.DigString("rev", info),
                "status": util.DigString("status", info),
            }
        }
    }
    return _return, nil
}
