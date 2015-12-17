package couch

import (
    _path "path"
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
            this.FileName = _path.Base(file)
        }
    }
    return this
}

func (this *DocumentAttachment) ToArray(encode bool) *DocumentAttachment {
    return this
}
