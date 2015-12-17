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
    var this = &DocumentAttachment{}
    if document != nil {
        this.Document = document
    }
    if file != "" {
        this.File = file
        if fileName != "" {
            this.FileName = fileName
        } else {
            this.FileName = _path.Base(fileName)
        }
    }
    return this
}
