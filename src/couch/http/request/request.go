package request

import (
    _fmt "fmt"
    _net "net"
    _bio "bufio"
    _str "strings"
    _b64 "encoding/base64"
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
    METHOD_POST   = "POST"
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

    if config["Username"] != "" && config["Password"] != "" {
        this.Headers["Authorization"] = "Basic "+
            _b64.StdEncoding.EncodeToString([]byte(
                config["Username"].(string) +":"+ config["Username"].(string)))
    }

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
    // link, _ := _net.Dial("tcp", "localhost:5984")
    link, err := _net.Dial("tcp",
        _fmt.Sprintf("%s:%v", this.Config["Host"], this.Config["Port"]))
    if err != nil {
        panic(err)
    }
    defer link.Close()

    var request, response string
    var url = u.ParseUrl(_fmt.Sprintf("%s://%s", this.Config["Scheme"], this.Uri))
    request += _fmt.Sprintf("%s %s?%s HTTP/%s\r\n", this.Method, url["Path"], url["Query"], this.HttpVersion)
    for key, value := range this.Headers {
        if !u.IsEmpty(value) {
            request += _fmt.Sprintf("%s: %s\r\n", key, value)
        }
    }
    request += "\r\n"
    request += this.GetBody()

    _fmt.Fprint(link, request)

    var reader = _bio.NewReader(link);

    status, err := reader.ReadString('\n')
    if status == "" {
        print("HTTP error: no response returned from server!\n")
        print("---------------------------------------------\n")
        print(request)
        print("---------------------------------------------\n")
        panic(err)
    }
    response += status

    for {
        var buffer = make([]byte, 1024)
        if read, _ := reader.Read(buffer); read == 0 {
            break // eof
        }
        response += _str.Trim(string(buffer), "\x00")
    }

    // @debug
    if this.Config["Couch.DEBUG"] == true {
        _dump(request)
        _dump(response)
    }

    return response
}

// @overwrite
func (this *Request) SetBody(body interface{}) {
    if body != nil &&
       this.Method != METHOD_HEAD &&
       this.Method != METHOD_GET {
        switch body.(type) {
            case string:
                // @overwrite
                var body = u.String(body)
                // trim null bytes & \r\n
                body = _str.Trim(body, "\x00")
                body = _str.TrimSpace(body)
                if this.GetHeader("Content-Type") == "application/json" {
                    // embrace with quotes for valid JSON body
                    body = u.Quote(body)
                }
                this.Body = body
                this.SetHeader("Content-Length", len(body))
            case map[string]interface{}:
                if this.GetHeader("Content-Type") == "application/json" {
                    // @overwrite
                    var body, err = u.UnparseBody(body)
                    if err != nil {
                        panic(err)
                    }
                    // pass empty body
                    if body != "{}" && body != "[]" {
                        this.Body = body
                        this.SetHeader("Content-Length", len(body))
                    }
                }
            default:
                var bodyType = _fmt.Sprintf("%T", body)
                if u.StringSearch(bodyType, "u?int(\\d+)?|float(32|64)") {
                    // @overwrite
                    var body = u.String(body)
                    this.Body = body
                    this.SetHeader("Content-Length", len(body))
                } else {
                    panic("Unsupported body type '"+ bodyType +"' given!");
                }
        }
    }
}

