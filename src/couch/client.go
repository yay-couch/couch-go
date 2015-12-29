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
// @uses    fmt, strings, regexp
// @uses    couch.util, couch.http
// @author  Kerem Güneş <qeremy[at]gmail[dot]com>
package couch

import (
    _fmt "fmt"
    _str "strings"
    _rex "regexp"
)

import (
    "couch/util"
    "couch/http"
)

// Client object.
type Client struct {
    Scheme    string
    Host      string
    Port      uint16
    Username  string
    Password  string
    Request   *http.Request
    Response  *http.Response
    Config    map[string]interface{}
}

// Default config options
var (
    Scheme        = "http"
    Host          = "localhost"
    Port   uint16 = 5984
    Username      = ""
    Password      = ""
)

// Constructor.
//
// @param  couch *couch.http.Couch
// @return *couch.Client
func NewClient(couch *Couch) *Client {
    var this = &Client{
        Scheme: Scheme,
          Host: Host,
          Port: Port,
      Username: Username,
      Password: Password,
    }

    var Config = util.Map()
    Config["Couch.NAME"]    = NAME
    Config["Couch.VERSION"] = VERSION
    // set default
    Config["Couch.DEBUG"]   = DEBUG

    // copy Couch configs
    var config = couch.GetConfig()
    if config != nil {
        for key, value := range config {
            Config[key] = value
        }
    }
    // add scheme if provided
    if scheme := config["Scheme"]; scheme != nil {
        this.Scheme = scheme.(string)
    }
    // add host if provided
    if host := config["Host"]; host != nil {
        this.Host = host.(string)
    }
    // add port if provided
    if port := config["Port"]; port != nil {
        this.Port = port.(uint16)
    }
    // add username if provided
    if username := config["Username"]; username != nil {
        this.Username = username.(string)
    }
    // add password if provided
    if password := config["Password"]; password != nil {
        this.Password = password.(string)
    }

    Config["Scheme"]   = this.Scheme
    Config["Host"]     = this.Host
    Config["Port"]     = this.Port
    Config["Username"] = this.Username
    Config["Password"] = this.Password

    // add debug if provided
    if debug := couch.Config["debug"]; debug != nil {
        Config["Couch.DEBUG"] = debug
    }

    this.Config = Config

    return this
}

// Get request object.
//
// @return *couch.http.Request
func (this *Client) GetRequest() *http.Request {
    return this.Request
}

// Get response object.
//
// @return *couch.http.Response
func (this *Client) GetResponse() *http.Response {
    return this.Response
}

// Perform a request.
//
// @param  uri       string
// @param  uriParams map[string]interface{}
// @param  body      interface{}
// @param  headers   map[string]interface{}
// @return *couch.http.Response
// @panics
func (this *Client) DoRequest(
    uri string, uriParams interface{},
    body interface{},
    headers interface{},
) *http.Response {
    // notation: GET /foo
    re, _ := _rex.Compile("^([A-Z]+)\\s+(/.*)")
    if re == nil {
        panic("Usage: <REQUEST METHOD> <REQUEST URI>!")
    }
    var match = re.FindStringSubmatch(uri)
    if len(match) < 3 {
        panic("Usage: <REQUEST METHOD> <REQUEST URI>!")
    }

    this.Request = http.NewRequest(this.Config)
    this.Response = http.NewResponse()

    // merge host, port and uri
    uri = _fmt.Sprintf("%s:%v/%s",
        this.Host, this.Port, _str.Trim(match[2], "/ "))

    // set request method & uri
    this.Request.SetMethod(match[1])
    this.Request.SetUri(uri, uriParams)

    // set request headers
    if headers, _ := headers.(map[string]interface{}); headers != nil {
        for key, value := range headers {
            this.Request.SetHeader(key, value)
        }
    }

    // set request body
    this.Request.SetBody(body)

    // perform request
    if result := this.Request.Send(); result != "" {
        tmp := make([]string, 2)
        tmp = _str.SplitN(result, "\r\n\r\n", 2)
        if len(tmp) != 2 {
            panic("No valid response returned from server!")
        }

        // set response status & headers
        if headers := util.ParseHeaders(_str.TrimSpace(tmp[0])); headers != nil {
            if status := headers["0"]; status != "" {
                this.Response.SetStatus(headers["0"])
            }
            for key, value := range headers {
                this.Response.SetHeader(key, value)
            }
        }

        // set response body
        this.Response.SetBody(tmp[1])
    }

    // error?
    if this.Response.GetStatusCode() >= 400 {
        // "" means -use self body-
        this.Response.SetError("")
    }

    return this.Response
}

// Perform a HEAD request.
//
//
// @param  uri       string
// @param  uriParams map[string]interface{}
// @param  headers   map[string]interface{}
// @return *couch.http.Response
// @panics
func (this *Client) Head(
    uri string, uriParams interface{},
    headers interface{},
) *http.Response {
    return this.DoRequest(http.METHOD_HEAD +" /"+ uri, uriParams, nil, headers)
}

// Perform a GET request.
//
// @param  uri       string
// @param  uriParams map[string]interface{}
// @param  headers   map[string]interface{}
// @return *couch.http.Response
// @panics
func (this *Client) Get(
    uri string, uriParams interface{},
    headers interface{},
) *http.Response {
    return this.DoRequest(http.METHOD_GET +" /"+ uri, uriParams, nil, headers)
}

// Perform a POST request.
//
// @param  uri       string
// @param  uriParams map[string]interface{}
// @param  body      interface{}
// @param  headers   map[string]interface{}
// @return *couch.http.Response
// @panics
func (this *Client) Post(
    uri string, uriParams interface{},
    body interface{},
    headers interface{},
) *http.Response {
    return this.DoRequest(http.METHOD_POST +" /"+ uri, uriParams, body, headers)
}

// Perform a PUT request.
//
// @param  uri       string
// @param  uriParams map[string]interface{}
// @param  body      interface{}
// @param  headers   map[string]interface{}
// @return *couch.http.Response
// @panics
func (this *Client) Put(
    uri string, uriParams interface{},
    body interface{},
    headers interface{},
) *http.Response {
    return this.DoRequest(http.METHOD_PUT +" /"+ uri, uriParams, body, headers)
}

// Perform a DELETE request.
//
// @param  uri       string
// @param  uriParams map[string]interface{}
// @param  headers   map[string]interface{}
// @return *couch.http.Response
// @panics
func (this *Client) Delete(
    uri string, uriParams interface{},
    headers interface{},
) *http.Response {
    return this.DoRequest(http.METHOD_DELETE +" /"+ uri, uriParams, nil, headers)
}

// Perform a COPY request.
//
// @param  uri       string
// @param  uriParams map[string]interface{}
// @param  headers   map[string]interface{}
// @return *couch.http.Response
// @panics
func (this *Client) Copy(
    uri string, uriParams interface{},
    headers interface{},
) *http.Response {
    return this.DoRequest(http.METHOD_COPY +" /"+ uri, uriParams, nil, headers)
}
