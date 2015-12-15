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
