package couch

import (
    "./util"
)

type Document struct {
    Id          string
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
    this.Id = id.(string)
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
        this.Data[key] = value
    }
}

func (this *Document) GetId() string {
    return this.Id
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
    if this.Id == "" {
        panic("_id field is could not be empty!")
    }
    var headers = util.Map()
    if (this.Rev != "") {
        headers["If-None-Match"] = util.Quote(this.Rev);
    }
    return (statusCode == this.Database.Client.
        Head(this.Database.Name +"/"+ this.Id, nil, headers).GetStatusCode())
}
