package couch

import (
    "./util"
)

type DocumentAttachment struct {
    Document    *Document
    File        string
    FileName    string
    Data        string
    DataLength  uint
    ContentType string
    Digest      string
}

func NewDocumentAttachment(document *Document, file, fileName string) *DocumentAttachment {
    var this = &DocumentAttachment{
        Document: document,
    }
    if file != "" {
        this.File = file
        if fileName != "" {
            this.FileName = fileName
        } else {
            this.FileName = util.Basename(file)
        }
    }
    return this
}
func (this *DocumentAttachment) SetDocument(document *Document) {
    this.Document = document
}
func (this *DocumentAttachment) GetDocument() *Document {
    return this.Document
}
func (this *DocumentAttachment) ToArray(encode bool) *DocumentAttachment {
    return this
}
func (this *DocumentAttachment) ToJson() *DocumentAttachment {
    return this
}
