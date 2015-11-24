package request

import (
    _fmt "fmt"
    _net "net"
    _bio "bufio"
    _str "strings"
)

import _stream "./../stream"
import _query "./../../query"

import u "./../../util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Request struct {
    _stream.Stream // extends :)
    Method         string
    Uri            string
    Config         map[string]interface{}
}

const (
    METHOD_HEAD   = "HEAD"
    METHOD_GET    = "GET"
    METHOD_POST   = "POS"
    METHOD_PUT    = "PUT"
    METHOD_COPY   = "COPY"
    METHOD_DELETE = "DELETE"
)

func Shutup() {}

func New(config map[string]interface{}) *Request {
    stream := _stream.New()
    stream.Type = _stream.TYPE_REQUEST
    stream.HttpVersion = "1.0"

    var this = &Request{
        Stream: *stream,
        Config: config,
    }

    // if config["Username"] != nil {
    //     // this.Headers["Authorization"] = "Basic "+
    //             // Util.Base64.encode(this.client.username +":"+ this.client.password);
    // }

    // this.Headers =
    this.Headers["Host"] = _fmt.Sprintf("%s:%v", config["Host"], config["Port"])
    this.Headers["Connection"] = "close"
    this.Headers["Accept"] = "application/json"
    this.Headers["Content-Type"] = "application/json"
    this.Headers["User-Agent"] = _fmt.Sprintf("%s/v%s (+http://github.com/qeremy/couch-go)",
        config["Couch.NAME"], config["Couch.VERSION"])

    return this
}

func (this *Request) SetMethod(method string) {
    method = _str.ToUpper(method)
    if (method != METHOD_HEAD &&
        method != METHOD_GET &&
        method != METHOD_POST) {
        this.SetHeader("X-HTTP-Method-Override", method)
    }
    this.Method = method
}

func (this *Request) SetUri(uri string, uriParams interface{}) {
    this.Uri = uri
    if uriParams == nil {
        return
    }
    var query = _query.New(uriParams.(map[string]interface{})).ToString()
    if query != "" {
        this.Uri += "?"+ query
    }
}

func (this *Request) Send() string {
    link, err := _net.Dial("tcp", _fmt.Sprintf("%s:%v", this.Config["Host"], this.Config["Port"]))
    if err != nil {
        panic(err)
    }
    defer link.Close()

    var request, response string
    var url = u.ParseUrl(_fmt.Sprintf("%s://%s:%v",
        this.Config["Scheme"], this.Config["Host"], this.Config["Port"]))

    request += _fmt.Sprintf("%s %s HTTP/%s\r\n",
        this.Method, url["Path"], this.HttpVersion)
    for key, value := range this.Headers {
        if !u.IsEmpty(value) {
            request += _fmt.Sprintf("%s: %s\r\n", key, value)
        }
    }
    request += "\r\n"
    request += this.GetBody()

    _fmt.Fprint(link, request)

    var reader = _bio.NewReader(link);
    for {
        var buffer = make([]byte, 1024)
        if read, _ := reader.Read(buffer); read == 0 {
            break // eof
        }
        response += string(buffer)
    }

    return response
}

// @overwrite?
// func (this *Request) SetBody() {}
