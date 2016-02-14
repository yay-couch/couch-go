// Copyright 2015 Kerem Güneş
//   <k-gun@mail.com>
//
// Apache License, Version 2.0
//   <http://www.apache.org/licenses/LICENSE-2.0>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// @package couch
// @uses    couch.util
// @author  Kerem Güneş <k-gun@mail.com>
package couch

import (
   "couch/util"
)

// @object DocumentAttachment
type DocumentAttachment struct {
   Document    *Document
   File        string
   FileName    string
   Data        string
   DataLength  int64
   ContentType string
   Digest      string
}

// Constructor.
//
// @param  document *couch.Document
// @param  string file
// @param  string fileName
// @return *couch.DocumentAttachment
func NewDocumentAttachment(document *Document, file, fileName string) (*DocumentAttachment) {
   this := &DocumentAttachment{
      Document: document,
   }

   if file != "" {
      this.File = file
      if fileName != "" {
         // set filename if provided
         this.FileName = fileName
      } else {
         // extract filename
         this.FileName = util.Basename(file)
      }
   }

   return this
}

// Set document.
//
// @param  document *couch.Document
// @return void
func (this *DocumentAttachment) SetDocument(document *Document) {
   this.Document = document
}

// Get document.
//
// @return document *couch.Document
func (this *DocumentAttachment) GetDocument() (*Document) {
   return this.Document
}

// Get attachment data as string array.
//
// @param  encode bool
// @return map[string]string
func (this *DocumentAttachment) ToArray(encode bool) (map[string]string) {
   // read file contents
   this.ReadFile(encode)

   array := util.MapString()
   array["data"] = this.Data
   array["content_type"] = this.ContentType

   return array
}

// Get attachment data as JSON string.
//
// @param  encode bool
// @return string
func (this *DocumentAttachment) ToJson(encode bool) (string) {
   json, _ := util.UnparseBody(this.ToArray(encode))
   return json
}

// Ping.
//
// @param  statusCodes... uint16
// @return bool
// @panics
func (this *DocumentAttachment) Ping(statusCodes... uint16) (bool) {
   if this.Document == nil {
      panic("Attachment document is not defined!")
   }

   docId := this.Document.GetId()
   docRev := this.Document.GetRev()
   if docId == "" {
      panic("Attachment document _id is required!")
   }
   if this.FileName == "" {
      panic("Attachment file name is required!")
   }

   query, headers := util.Map(), util.Map()
   if docRev != "" {
      query["rev"] = docRev
   }
   if this.Digest != "" {
      headers["If-None-Match"] = util.Quote(this.Digest)
   }

   database := this.Document.GetDatabase()
   response := database.Client.Head(util.StringFormat("%s/%s/%s",
      database.Name, docId, util.UrlEncode(this.FileName)), query, headers)

   // try to match given status codes
   for _, statusCode := range statusCodes {
      if response.GetStatusCode() == statusCode {
         return true
      }
   }

   return false
}

// Find.
//
// @return map[string]interface{}
// @panics
func (this *DocumentAttachment) Find() (map[string]interface{}) {
   if this.Document == nil {
      panic("Attachment document is not defined!")
   }

   docId := this.Document.GetId()
   docRev := this.Document.GetRev()
   if docId == "" {
      panic("Attachment document _id is required!")
   }
   if this.FileName == "" {
      panic("Attachment file name is required!")
   }

   query, headers := util.Map(), util.Map()
   if docRev != "" {
      query["rev"] = docRev
   }
   if this.Digest != "" {
      headers["If-None-Match"] = util.Quote(this.Digest)
   }
   headers["Accept"] = "*/*"
   headers["Content-Type"] = nil // nil=remove

   database := this.Document.GetDatabase()
   response := database.Client.Get(util.StringFormat("%s/%s/%s",
      database.Name, docId, util.UrlEncode(this.FileName)), query, headers)
   statusCode := response.GetStatusCode()

   ret := util.Map()
   // try to match excepted status code
   if  statusCode == 200 || statusCode == 304 {
      ret["content"] = response.GetBody()
      ret["content_type"] = response.GetHeader("Content-Type")
      ret["content_length"] = util.UInt(response.GetHeader("Content-Length"))
      // set digest
      md5 := response.GetHeader("Content-MD5")
      if md5 == nil {
         md5 = response.GetHeader("ETag")
      }
      ret["digest"] = "md5-"+ util.Trim(md5.(string), "\"")
   }

   return ret
}

// Save.
//
// @return map[string]interface{}, error
// @panics
func (this *DocumentAttachment) Save() (map[string]interface{}, error) {
   if this.Document == nil {
      panic("Attachment document is not defined!")
   }

   docId := this.Document.GetId()
   docRev := this.Document.GetRev()
   if docId == "" {
      panic("Attachment document _id is required!")
   }
   if docRev == "" {
      panic("Attachment document _rev is required!")
   }
   if this.FileName == "" {
      panic("Attachment file name is required!")
   }

   // read file contents
   this.ReadFile(false)

   headers := util.Map()
   headers["If-Match"] = docRev
   headers["Content-Type"] = this.ContentType

   data, err := this.Document.Database.Client.Put(util.StringFormat(
         "%s/%s/%s", this.Document.Database.Name, docId, util.UrlEncode(this.FileName),
      ), nil, this.Data, headers,
   ).GetBodyData(nil)
   if err != nil {
      return nil, err
   }

   return map[string]interface{}{
       "ok": util.DigBool("ok", data),
       "id": util.DigString("id", data),
      "rev": util.DigString("rev", data),
   }, nil
}

// Remove.
//
// @return map[string]interface{}, error
// @panics
func (this *DocumentAttachment) Remove(args... bool) (map[string]interface{}, error) {
   if this.Document == nil {
      panic("Attachment document is not defined!")
   }

   docId := this.Document.GetId()
   docRev := this.Document.GetRev()
   if docId == "" {
      panic("Attachment document _id is required!")
   }
   if docRev == "" {
      panic("Attachment document _rev is required!")
   }
   if this.FileName == "" {
      panic("Attachment file name is required!")
   }

   query, headers := util.Map(), util.Map()
   if args != nil {
      if args[0] {
         query["batch"] = "ok"
      }
      if args[1] {
         headers["X-Couch-Full-Commit"] = "true"
      }
   }
   headers["If-Match"] = docRev

   data, err := this.Document.Database.Client.Delete(util.StringFormat(
         "%s/%s/%s", this.Document.Database.Name, docId, util.UrlEncode(this.FileName),
      ), nil, headers,
   ).GetBodyData(nil)
   if err != nil {
      return nil, err
   }

   return map[string]interface{}{
       "ok": util.DigBool("ok", data),
       "id": util.DigString("id", data),
      "rev": util.DigString("rev", data),
   }, nil
}

// Read file.
//
// @return void
// @panics
func (this *DocumentAttachment) ReadFile(encode bool) {
   if this.File == "" {
      panic("Attachment file is empty!")
   }

   // get file info
   info, err := util.FileInfo(this.File)
   if err != nil {
      panic(err)
   }
   this.ContentType = util.String(info["mime"])

   // get file contents
   data, err := util.FileGetContents(this.File)
   if err != nil {
      panic(err)
   }

   this.Data = data
   // convert to base64 if encode=true
   if encode {
      this.Data = util.Base64Encode(data)
   }

   this.DataLength = int64(len(data))
}
