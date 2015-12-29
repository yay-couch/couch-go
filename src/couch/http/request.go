// Copyright 2015 Kerem Güneş
//    <http://qeremy.com>
//
// Apache License, Version 2.0
//    <http://www.apache.org/licenses/LICENSE-2.0>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// @package couch
// @uses    fmt, net, bufio, strings
// @uses    couch.util, couch.query
// @author  Kerem Güneş <qeremy[at]gmail[dot]com>
package http

import (
    _fmt "fmt"
    _net "net"
    _bio "bufio"
    _str "strings"
)

import (
    "couch/util"
    "couch/query"
)

// @object couch.http.Request
type Request struct {
    Stream // extends :)
    Method         string
    Uri            string
    Config         map[string]interface{}
}

// Request methods.
// @const string
const (
    METHOD_HEAD   = "HEAD"
    METHOD_GET    = "GET"
    METHOD_POST   = "POST"
    METHOD_PUT    = "PUT"
    METHOD_COPY   = "COPY"
    METHOD_DELETE = "DELETE"
)

// Constructor.
//
// @param  config map[string]interface{}
// @return (*couch.http.Request)
func NewRequest(config map[string]interface{}) *Request {
    stream := NewStream()
    stream.Type = TYPE_REQUEST
    stream.HttpVersion = "1.0"

    var this = &Request{
        Stream: *stream,
        Config: config,
    }

    // add default headers
    this.Headers["Host"] = _fmt.Sprintf("%s:%v", config["Host"], config["Port"])
    this.Headers["Connection"] = "close"
    this.Headers["Accept"] = "application/json"
    this.Headers["Content-Type"] = "application/json"
    this.Headers["User-Agent"] = _fmt.Sprintf("%s/v%s (+http://github.com/qeremy/couch-go)",
        config["Couch.NAME"], config["Couch.VERSION"])

    // add auth header
    if config["Username"] != "" && config["Password"] != "" {
        this.Headers["Authorization"] = "Basic "+
            util.Base64Encode(_fmt.Sprintf("%s:%s", config["Username"], config["Username"]))
    }

    return this
}

// Set method.
//
// @param  method string
// @return (void)
func (this *Request) SetMethod(method string) {
    method = _str.ToUpper(method)
    // add method override header
    if (method != METHOD_HEAD &&
        method != METHOD_GET &&
        method != METHOD_POST) {
        this.SetHeader("X-HTTP-Method-Override", method)
    }

    this.Method = method
}

// Set URI with params.
//
// @param  uri       string
// @param  uriParams interface{}
// @return (void)
func (this *Request) SetUri(uri string, uriParams interface{}) {
    this.Uri = uri
    if uriParams == nil {
        return
    }

    // append params if provided
    var query = query.New(uriParams.(map[string]interface{})).ToString()
    if query != "" {
        this.Uri += "?"+ query
    }
}

// Send!
//
// @return (string)
// @panics
func (this *Request) Send() string {
    link, err := _net.Dial("tcp", _fmt.Sprintf("%s:%v", this.Config["Host"], this.Config["Port"]))
    if err != nil {
        panic(err)
    }
    defer link.Close()

    var send, recv string
    var url = util.ParseUrl(_fmt.Sprintf("%s://%s", this.Config["Scheme"], this.Uri))

    // add first line & headers
    send += _fmt.Sprintf("%s %s?%s HTTP/%s\r\n",
        this.Method, url["Path"], url["Query"], this.HttpVersion)
    for key, value := range this.Headers {
        if !util.IsEmpty(value) {
            send += _fmt.Sprintf("%s: %s\r\n", key, value)
        }
    }
    send += "\r\n"
    send += this.GetBody()
    _fmt.Fprint(link, send)

    var reader = _bio.NewReader(link)

    status, err := reader.ReadString('\n')
    if status == "" {
        print("HTTP error: no response returned from server!\n")
        print("---------------------------------------------\n")
        print(send)
        print("---------------------------------------------\n")
        panic(err)
    }
    recv += status

    for {
        var buffer = make([]byte, 1024)
        if read, _ := reader.Read(buffer); read == 0 {
            break // eof
        }
        // yes, we've allocated to much..
        recv += _str.Trim(string(buffer), "\x00")
    }

    // @debug
    if this.Config["Couch.DEBUG"] == true {
        util.Dump(send)
        util.Dump(recv)
    }

    return recv
}

// Set body.
//
// @param  body interface{}
// @return (void)
// @panics
// @implements
func (this *Request) SetBody(body interface{}) {
    if body != nil &&
       // these methods not allowed for body
       this.Method != METHOD_HEAD &&
       this.Method != METHOD_GET {
        switch body := body.(type) {
            case string:
                // @overwrite
                if this.GetHeader("Content-Type") == "application/json" {
                    // embrace with quotes for valid JSON body
                    body = util.Quote(body)
                }
                this.Body = body
            default:
                var bodyType = _fmt.Sprintf("%T", body)
                if util.StringSearch(bodyType, "^u?int(\\d+)?|float(32|64)$") {
                    // @overwrite
                    this.Body = util.String(body)
                } else {
                    if this.GetHeader("Content-Type") == "application/json" {
                        // @overwrite
                        body, err := util.UnparseBody(body)
                        if err != nil {
                            panic(err)
                        }
                        this.Body = body
                    }
                }
        }

        // auto-set content length headers
        this.SetHeader("Content-Length", len(this.Body.(string)))
    }
}

// Get request as string.
//
// @return (string)
// @implements
func (this *Request) ToString() string {
    var ret = util.StringFormat("%s %s HTTP/%s\r\n", this.Method, this.Uri, this.HttpVersion)
    if this.Headers != nil {
        for key, value := range this.Headers {
            if (value != nil) {
                ret += util.StringFormat("%s: %s\r\n", key, value)
            }
        }
    }
    ret += "\r\n"

    if this.Body != nil {
        switch this.Body.(type) {
            case string:
                ret += this.Body.(string)
            default:
                body, err := util.UnparseBody(this.Body)
                if err == nil {
                    ret += body
                }
        }
    }

    return ret
}
