package couch

import (
    "./util"
)

type DocumentAttachment struct {
    Document    *Document
    File        string
    FileName    string
    Data        string
    DataLength  int64
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
func (this *DocumentAttachment) ToArray(encode bool) map[string]string {
    this.ReadFile(encode)
    var array = util.MapString()
    array["data"] = this.Data
    array["content_type"] = this.ContentType
    return array
}
func (this *DocumentAttachment) ToJson(encode bool) string {
    json, _ := util.UnparseBody(this.ToArray(encode))
    return json
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
    var query = util.Map()
    if docRev != "" {
        query["rev"] = docRev
    }
    var headers = util.Map()
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

// @todo return DocumentAttachment?
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
    var query = util.Map()
    if docRev != "" {
        query["rev"] = docRev
    }
    var headers = util.Map()
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

// @todo return DocumentAttachment?
func (this *DocumentAttachment) Save() (map[string]interface{}, error) {
    if this.Document == nil {
        panic("Attachment document is not defined!")
    }
    var docId = this.Document.GetId()
    var docRev = this.Document.GetRev()
    if docId == "" {
        panic("Attachment document _id is required!")
    }
    if docRev == "" {
        panic("Attachment document _rev is required!")
    }
    if this.FileName == "" {
        panic("Attachment file name is required!")
    }
    this.ReadFile(false)
    var headers = util.Map()
    headers["If-Match"] = docRev
    headers["Content-Type"] = this.ContentType
    data, err := this.Document.Database.Client.Put(util.StringFormat("%s/%s/%s",
        this.Document.Database.Name, docId, util.UrlEncode(this.FileName)), nil,
        this.Data, headers).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "ok": util.DigBool("ok", data),
        "id": util.DigString("id", data),
        "rev": util.DigString("rev", data),
    }, nil
}

func (this *DocumentAttachment) ReadFile(encode bool) {
    if this.File == "" {
        panic("Attachment file is empty!")
    }
    info, err := util.FileInfo(this.File)
    if err != nil {
        panic(err)
    }
    this.ContentType = util.String(info["mime"])
    data, err := util.FileGetContents(this.File)
    if err != nil {
        panic(err)
    }
    this.Data = data
    if encode {
        this.Data = util.Base64Encode(data)
    }
    this.DataLength = int64(len(data))
}
