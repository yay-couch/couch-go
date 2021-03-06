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

// @package    couch
// @subpackage couch.http
// @uses       fmt
// @uses       couch.util
// @author     Kerem Güneş <k-gun@mail.com>
package http

import (
   _fmt "fmt"
)

import (
   "couch/util"
)

// @object couch.http.Stream
type Stream struct {
   Type        uint8
   HttpVersion string
   Headers     map[string]interface{}
   Body        interface{}
   Error       string
   ErrorData   map[string]string
   StreamBody
   ToString    func() string
}

// @object couch.http.StreamError
type StreamError struct {
   ErrorKey   string `json:"error"`
   ErrorValue  string `json:"reason"`
}

// @object couch.http.StreamBody
type StreamBody interface {
   SetBody(body interface{})
}

// Stream types.
// @const uint8
const (
   TYPE_REQUEST  = 1
   TYPE_RESPONSE = 2
)

func Shutup() {}

// Constructor.
//
// @param  type_       uint8
// @param  httpVersion string
// @return couch.http.Stream
func NewStream(type_ uint8, httpVersion string) (*Stream) {
   return &Stream{
      Type: type_,
      HttpVersion: httpVersion,
      Headers: util.Map(),
   }
}

// Set header.
//
// @param  key   string
// @param  value interface{}
// @return void
// @panic
func (this *Stream) SetHeader(key string, value interface{}) {
   switch value.(type) {
      case nil: // so nil means remove
         delete(this.Headers, key)
      case int,
          bool,
          string:
         this.Headers[key] = util.String(value)
      default:
         panic("Unsupported value type '"+ _fmt.Sprintf("%T", value) +"' given!")
   }
}

// Get header.
//
// @param  key string
// @return interface{}
func (this *Stream) GetHeader(key string) (interface{}) {
   if value, ok := this.Headers[key]; ok {
      return value
   }

   return nil
}

// Get all headers.
//
// @return map[string]interface{}
func (this *Stream) GetHeaderAll() (map[string]interface{}) {
   return this.Headers
}

// Get body.
//
// @return string
func (this *Stream) GetBody() (string) {
   if this.Body == nil {
      return ""
   }

   return this.Body.(string)
}

// Get body data (parsed).
//
// @param  to interface{}
// @return interface{}, error
func (this *Stream) GetBodyData(to interface{}) (interface{}, error) {
   // error?
   if this.Error != "" {
      data, err := util.ParseBody(this.Body.(string), &StreamError{})
      if err != nil {
         return nil, err
      }

      return nil, _fmt.Errorf("Stream Error >> error: \"%s\", reason: \"%s\"",
         data.(*StreamError).ErrorKey,
         data.(*StreamError).ErrorValue,
      )
   }

   data, err := util.ParseBody(this.Body.(string), to)
   if err != nil {
      return nil, err
   }

   return data, nil
}

// Set error.
//
// @param  body string
// @return void
func (this *Stream) SetError(body string) {
   if body == "" {
      body = this.Body.(string)
   }

   data, err := util.ParseBody(body, &StreamError{})
   if data != nil && err == nil {
      errorKey   := data.(*StreamError).ErrorKey
      errorValue := data.(*StreamError).ErrorValue

      this.Error = _fmt.Sprintf("Stream Error >> error: \"%s\", reason: \"%s\"",
         errorKey,
         errorValue,
      )

      this.ErrorData = util.MapString()
      this.ErrorData["error"]  = errorKey
      this.ErrorData["reason"] = errorValue
   }
}

// Get error.
//
// @return string
func (this *Stream) GetError() (string) {
   return this.Error
}

// Get error value.
//
// @return string
func (this *Stream) GetErrorValue(key string) (string) {
   return this.ErrorData[key]
}

// Get stream as string
//
// @param  firstLine string
// @return string
// @protected
func (this *Stream) toString(firstLine string) (string) {
   str := firstLine

   // add headers
   if this.Headers != nil {
      for key, value := range this.Headers {
         if key == "0" { // response only
            continue
         }
         if (value != nil) { // remove?
            str += util.StringFormat("%s: %s\r\n", key, value)
         }
      }
   }

   // add headers/body seperator
   str += "\r\n"

   // add body
   if this.Body != nil {
      switch this.Body.(type) {
         case string:
            str += this.Body.(string)
         default:
            body, err := util.UnparseBody(this.Body)
            if err == nil {
               str += body
            }
      }
   }

   return str
}
