package couch

import (
    _fmt "fmt"
    _str "strings"
    _rex "regexp"
)

import (
    "couch/http"
    "couch/util"
)

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

var (
    Scheme        = "http"
    Host          = "localhost"
    Port   uint16 = 5984
    Username      = ""
    Password      = ""
)

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
    Config["Couch.DEBUG"]   = DEBUG // set default

    var config = couch.GetConfig()
    if config != nil {
        for key, value := range config {
            Config[key] = value
        }
    }
    if scheme := config["Scheme"]; scheme != nil {
        this.Scheme = scheme.(string)
    }
    if host := config["Host"]; host != nil {
        this.Host = host.(string)
    }
    if port := config["Port"]; port != nil {
        this.Port = port.(uint16)
    }
    if username := config["Username"]; username != nil {
        this.Username = username.(string)
    }
    if password := config["Password"]; password != nil {
        this.Password = password.(string)
    }

    Config["Scheme"]   = this.Scheme
    Config["Host"]     = this.Host
    Config["Port"]     = this.Port
    Config["Username"] = this.Username
    Config["Password"] = this.Password

    if debug := couch.Config["debug"]; debug != nil {
        Config["Couch.DEBUG"] = debug
    }

    this.Config = Config

    return this
}

func (this *Client) GetRequest() *http.Request {
    return this.Request
}
func (this *Client) GetResponse() *http.Response {
    return this.Response
}

func (this *Client) DoRequest(uri string, uriParams interface{},
        body interface{}, headers interface{}) *http.Response {
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

    uri = _fmt.Sprintf("%s:%v/%s", this.Host, this.Port, _str.Trim(match[2], "/ "))

    this.Request.SetMethod(match[1])
    this.Request.SetUri(uri, uriParams)
    if headers, _ := headers.(map[string]interface{}); headers != nil {
        for key, value := range headers {
            this.Request.SetHeader(key, value)
        }
    }
    this.Request.SetBody(body)

    if result := this.Request.Send(); result != "" {
        tmp := make([]string, 2)
        tmp = _str.SplitN(result, "\r\n\r\n", 2)
        if len(tmp) != 2 {
            panic("No valid response returned from server!")
        }
        if headers := util.ParseHeaders(_str.TrimSpace(tmp[0])); headers != nil {
            if status := headers["0"]; status != "" {
                this.Response.SetStatus(headers["0"])
            }
            for key, value := range headers {
                this.Response.SetHeader(key, value)
            }
        }
        this.Response.SetBody(tmp[1])
    }

    // error?
    if this.Response.GetStatusCode() >= 400 {
        // "" means -use self body-
        this.Response.SetError("")
    }

    return this.Response
}

func (this *Client) Head(uri string, uriParams interface{},
        headers interface{}) *http.Response {
    return this.DoRequest(http.METHOD_HEAD +" /"+ uri, uriParams, nil, headers)
}
func (this *Client) Get(uri string, uriParams interface{},
        headers interface{}) *http.Response {
    return this.DoRequest(http.METHOD_GET +" /"+ uri, uriParams, nil, headers)
}
func (this *Client) Post(uri string, uriParams interface{}, body interface{},
        headers interface{}) *http.Response {
    return this.DoRequest(http.METHOD_POST +" /"+ uri, uriParams, body, headers)
}
func (this *Client) Put(uri string, uriParams interface{}, body interface{},
        headers interface{}) *http.Response {
    return this.DoRequest(http.METHOD_PUT +" /"+ uri, uriParams, body, headers)
}
func (this *Client) Delete(uri string, uriParams interface{},
        headers interface{}) *http.Response {
    return this.DoRequest(http.METHOD_DELETE +" /"+ uri, uriParams, nil, headers)
}
func (this *Client) Copy(uri string, uriParams interface{},
        headers interface{}) *http.Response {
    return this.DoRequest(http.METHOD_COPY +" /"+ uri, uriParams, nil, headers)
}
