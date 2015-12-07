package couch

import (
    _fmt "fmt"
    _str "strings"
    _rex "regexp"
)

import (
    "./http"
    "./util"
)

type Client struct {
    Scheme    string
    Host      string
    Port      uint16
    Username, Password string
    Request   *http.Request
    Response  *http.Response
    Config    map[string]interface{}
}

func New(config interface{}, username, password string) *Client {
    var this = &Client{}

    // config in config, just for import cycle..
    if config != nil {
        this.Config = make(map[string]interface{})
        switch config.(type) {
            case string:
                var url = util.ParseUrl(config.(string));
                var scheme, host, port =
                    url["Scheme"], url["Host"], util.Number(url["Port"], "uint16").(uint16)
                this.Scheme = scheme; this.Config["Scheme"] = scheme
                this.Host   = host;   this.Config["Host"] = host
                this.Port   = port;   this.Config["Port"] = port
            case map[string]interface{}:
                // copy
                for key, value := range config.(map[string]interface{}) {
                    this.Config[key] = value
                }
                this.Scheme = this.Config["Scheme"].(string)
                this.Host   = this.Config["Host"].(string)
                this.Port   = this.Config["Port"].(uint16)
        }
    }

    this.Username = username
    this.Password = password

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
