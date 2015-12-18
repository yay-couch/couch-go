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

func (this *DocumentAttachment) Ping(statusCodes ...uint16) bool {
    if this.Document == nil {
        panic("Attachment document is not defined!")
    }
    var docId = this.Document.GetId()
    var docRev = this.Document.GetRev()
    if docId == "" {
        panic("Attachment document _id is required!")
    }
    if this.FileName == "" {
        panic("Attachment file name is required!")
    }
    var query = util.Param(nil)
    if docRev != "" {
        query["rev"] = docRev
    }
    var headers = util.Param(nil)
    if this.Digest != "" {
        headers["If-None-Match"] = util.Quote(this.Digest)
    }
    var database = this.Document.GetDatabase()
    var response = database.Client.Head(util.StringFormat("%s/%s/%s",
        database.Name, docId, util.UrlEncode(this.FileName)), query, headers)
    for _, statusCode := range statusCodes {
        if response.GetStatusCode() == statusCode {
            return true
        }
    }
    return false
}

func (this *DocumentAttachment) Find() map[string]interface{} {
    if this.Document == nil {
        panic("Attachment document is not defined!")
    }
    var docId = this.Document.GetId()
    var docRev = this.Document.GetRev()
    if docId == "" {
        panic("Attachment document _id is required!")
    }
    if this.FileName == "" {
        panic("Attachment file name is required!")
    }
    var query = util.Param(nil)
    if docRev != "" {
        query["rev"] = docRev
    }
    var headers = util.Param(nil)
    if this.Digest != "" {
        headers["If-None-Match"] = util.Quote(this.Digest)
    }
    headers["Accept"] = "*/*"
    headers["Content-Type"] = nil

    var _return = util.Map()
    var database = this.Document.GetDatabase()
    var response = database.Client.Get(util.StringFormat("%s/%s/%s",
        database.Name, docId, util.UrlEncode(this.FileName)), query, headers)
    var statusCode = response.GetStatusCode()
    if  statusCode == 200 || statusCode == 304 {
        _return["content"] = response.GetBody()
        _return["content_type"] = response.GetHeader("Content-Type")
        _return["content_length"] = util.UInt(response.GetHeader("Content-Length"))
        var md5 = response.GetHeader("Content-MD5")
        if md5 == nil {
            md5 = response.GetHeader("ETag")
        }
        _return["digest"] = "md5-"+ util.Trim(md5.(string), "\"")
    }
    return _return
}
