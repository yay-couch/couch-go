package document

import u "./../util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Document struct {
    Id          string
    Rev         string
    Deleted     bool
    Attachments []DocumentAttachment
    Data        map[string]interface{}
    // Database    *_database.Database
}

func Shutup() {}

func NewDocument() *Document {
    return &Document{
        //
    }
}

// func (this *Document) SetId(id interface{}) {
//     this.Id = u.String(id)
// }
// func (this *Document) SetRev(rev string) {
//     this.Rev = rev
// }
// func (this *Document) SetRev(deleted bool) {
//     this.Deleted = deleted
// }
// func (this *Document) SetData(data map[string]interface{}) {
//     if this.Data == nil {
//         this.Data = make(map[string]interface{})
//     }
//     for key, value := range data.(map[string]interface{}) {
//         this.Data[key] = value
//     }
// }
