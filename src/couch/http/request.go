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

type Request struct {
    Stream // extends :)
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

func NewRequest(config map[string]interface{}) *Request {
    stream := NewStream()
    stream.Type = TYPE_REQUEST
    stream.HttpVersion = "1.0"

    var this = &Request{
        Stream: *stream,
        Config: config,
    }

    if config["Username"] != "" && config["Password"] != "" {
        this.Headers["Authorization"] = "Basic "+
            util.Base64Encode(_fmt.Sprintf("%s:%s", config["Username"], config["Username"]))
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
    var query = query.New(uriParams.(map[string]interface{})).ToString()
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
    var url = util.ParseUrl(_fmt.Sprintf("%s://%s", this.Config["Scheme"], this.Uri))
    request += _fmt.Sprintf("%s %s?%s HTTP/%s\r\n",
        this.Method, url["Path"], url["Query"], this.HttpVersion)
    for key, value := range this.Headers {
        if !util.IsEmpty(value) {
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
        util.Dump(request)
        util.Dump(response)
    }

    return response
}

// @implement
func (this *Request) SetBody(body interface{}) {
    if body != nil &&
       this.Method != METHOD_HEAD &&
       this.Method != METHOD_GET {
        switch body.(type) {
            case string:
                // @overwrite
                var body = util.String(body)
                if this.GetHeader("Content-Type") == "application/json" {
                    // embrace with quotes for valid JSON body
                    body = util.Quote(body)
                }
                this.Body = body
            default:
                var bodyType = _fmt.Sprintf("%T", body)
                if util.StringSearch(bodyType, "^u?int(\\d+)?|float(32|64)$") {
                    // @overwrite
                    var body = util.String(body)
                    this.Body = body
                } else {
                    if this.GetHeader("Content-Type") == "application/json" {
                        // @overwrite
                        body, err := util.UnparseBody(body)
                        if err != nil {
                            panic(err)
                        }
                        this.Body = body
                    }
                    // panic("Unsupported body type '"+ bodyType +"' given!");
                }
        }
        this.SetHeader("Content-Length", len(this.Body.(string)))
    }
}

// @implement
func (this *Request) ToString() string {
    var ret string
    ret = util.StringFormat("%s %s HTTP/%s\r\n", this.Method, this.Uri, this.HttpVersion)
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
