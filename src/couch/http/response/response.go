package response

import (
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
        this.HttpVersion = match[1]
        responseCode, _ := _strc.Atoi(match[2])
        if responseCode != 0 {
            this.StatusCode = uint16(responseCode)
        }
        this.StatusText = _str.TrimSpace(match[3])
    }
}

// @overwrite?
// func (this *Request) SetBody() {}
