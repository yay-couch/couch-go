package response

import (
    _fmt "fmt"
    _rex "regexp"
    _str "strings"
    _strc "strconv"
)

import _stream "./../stream"

import u "./../../util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Response struct {
    _stream.Stream // extends :)
    Status       string
    StatusCode   uint16
    StatusText   string
}

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

func Shutup() {}

func New() *Response {
    stream := _stream.New()
    stream.Type = _stream.TYPE_RESPONSE
    stream.HttpVersion = "1.1"
    var this = &Response{
        Stream: *stream,
    }
    return this
}

func (this *Response) SetStatus(status string) {
    // status line >> HTTP/1.0 200 OK
    re, _ := _rex.Compile("^HTTP/(\\d+\\.\\d+)\\s+(\\d+)\\s+(.+)")
    if re == nil {
        return
    }
    this.Status = _str.TrimSpace(status)
    var match = re.FindStringSubmatch(status)
    if len(match) == 4 {
        // update http version
        this.HttpVersion = match[1]
        responseCode, _ := _strc.Atoi(match[2])
        responseText    := _str.TrimSpace(match[3])
        this.SetStatusCode(uint16(responseCode))
        this.SetStatusText(string(responseText))
    }
}
func (this *Response) SetStatusCode(code uint16) {
    this.StatusCode = code
}
func (this *Response) SetStatusText(text string) {
    this.StatusText = text
}

func (this *Response) GetStatus() string {
    return this.Status
}
func (this *Response) GetStatusCode() uint16 {
    return this.StatusCode
}
func (this *Response) GetStatusText() string {
    return this.StatusText
}

// @overwrite
func (this *Response) SetBody(body interface{}) {
    if body != nil {
        // @overwrite
        var body = _fmt.Sprintf("%s", body)
        // trim null bytes & \r\n
        body = _str.Trim(body, "\x00")
        body = _str.TrimSpace(body)
        this.Body = body
    }
}
