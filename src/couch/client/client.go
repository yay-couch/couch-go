package client

import (
    _fmt "fmt"
    _str "strings"
    _rex "regexp"
)

// import _stream "./../http/stream"
import _request "./../http/request"
import _response "./../http/response"

import u "./../util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Client struct {
    Scheme    string
    Host      string
    Port      uint16
    Username, Password string
    Request   *_request.Request
    Response  *_response.Response
    Config    map[string]interface{}
}

func Shutup() {}

func New(config interface{}, username, password string) *Client {
    var this = &Client{}

    // config in config, just for import cycle..
    if config != nil {
        this.Config = make(map[string]interface{})
        switch config.(type) {
            case string:
                var url = u.ParseUrl(config.(string));
                var scheme, host, port =
                    url["Scheme"], url["Host"], u.Number(url["Port"], "uint16").(uint16)
                this.Scheme = scheme; this.Config["Scheme"] = scheme
                this.Host   = host;   this.Config["Host"] = host
                this.Port   = port;   this.Config["Port"] = port
            case map[string]interface{}:
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

func (this *Client) GetRequest() *_request.Request {
    return this.Request
}
func (this *Client) GetResponse() *_response.Response {
    return this.Response
}

func (this *Client) DoRequest(uri string, uriParams interface{},
        body interface{}, headers interface{}) *_response.Response {
    re, _ := _rex.Compile("^([A-Z]+)\\s+(/.*)")
    if re == nil {
        panic("Usage: <REQUEST METHOD> <REQUEST URI>!")
    }
    var match = re.FindStringSubmatch(uri)
    if len(match) < 3 {
        panic("Usage: <REQUEST METHOD> <REQUEST URI>!")
    }

    this.Request = _request.New(this.Config)
    this.Response = _response.New()

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
        if headers := u.ParseHeaders(_str.TrimSpace(tmp[0])); headers != nil {
            if status := headers["0"]; status != "" {
                this.Response.SetStatus(headers["0"])
            }
            for key, value := range headers {
                this.Response.SetHeader(key, value)
            }
        }
        var body = _str.TrimSpace(tmp[1])
        this.Response.SetBody(body)
    }
    return this.Response
}

func (this *Client) Head(uri string, uriParams interface{}, headers interface{}) *_response.Response {
    return this.DoRequest(_request.METHOD_HEAD +" /"+ uri, uriParams, nil, headers)
}
