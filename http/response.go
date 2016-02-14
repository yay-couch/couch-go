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
// @uses       regexp, strings, strconv
// @uses       couch.util
// @author     Kerem Güneş <k-gun@mail.com>
package http

import (
   _rex "regexp"
   _str "strings"
   _strc "strconv"
)

import (
   "couch/util"
)

// @object couch.http.Response
type Response struct {
   Stream // extends
   Status     string
   StatusCode uint16
   StatusText string
}

// Response statuses.
// @var map[int]string{}
var (
   STATUS = map[int]string{
      200: "OK",
      201: "Created",
      202: "Accepted",
      304: "Not Modified",
      400: "Bad Request",
      401: "Unauthorized",
      403: "Forbidden",
      404: "Not Found",
      405: "Resource Not Allowed",
      406: "Not Acceptable",
      409: "Conflict",
      412: "Precondition Failed",
      415: "Bad Content Type",
      416: "Requested Range Not Satisfiable",
      417: "Expectation Failed",
      500: "Internal Server Error",
   }
)

// Constructor.
//
// @return *couch.http.Response
func NewResponse() (*Response) {
   this := &Response{
      Stream: *NewStream(TYPE_RESPONSE, "1.0"),
   }

   return this
}

// Set status.
//
// @param  status string
// @return void
func (this *Response) SetStatus(status string) {
   // status line >> HTTP/1.0 200 OK
   re, _ := _rex.Compile("^HTTP/(\\d+\\.\\d+)\\s+(\\d+)\\s+(.+)")
   if re == nil {
      return
   }
   this.Status = _str.TrimSpace(status)

   match := re.FindStringSubmatch(status)
   if len(match) == 4 {
      // update http version
      this.HttpVersion = match[1]
      // set status code/text
      statusCode, _ := _strc.Atoi(match[2])
      statusText   := _str.TrimSpace(match[3])
      this.SetStatusCode(uint16(statusCode))
      this.SetStatusText(string(statusText))
   }
}

// Set status code.
//
// @param  statusCode uint16
// @return void
func (this *Response) SetStatusCode(statusCode uint16) {
   this.StatusCode = statusCode
}

// Set status text.
//
// @param  statusText string
// @return void
func (this *Response) SetStatusText(statusText string) {
   this.StatusText = statusText
}

// Get status.
//
// @return string
func (this *Response) GetStatus() (string) {
   return this.Status
}

// Get status code.
//
// @return uint16
func (this *Response) GetStatusCode() (uint16) {
   return this.StatusCode
}

// Get status text.
//
// @return string
func (this *Response) GetStatusText() (string) {
   return this.StatusText
}

// Set body.
//
// @param  body interface{}
// @return void
// @implemented
func (this *Response) SetBody(body interface{}) {
   if body != nil {
      this.Body = util.String(body)
   }
}

// Get response as string.
//
// @return string
// @implemented
func (this *Response) ToString() (string) {
   return this.toString(util.StringFormat(
      "HTTP/%s %d %s\r\n", this.HttpVersion, this.StatusCode, this.StatusText,
   ))
}
