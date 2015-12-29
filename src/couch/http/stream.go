package http

import (
    _fmt "fmt"
)

import (
    "couch/util"
)

type Stream struct {
    Type        uint8
    HttpVersion string
    Headers     map[string]interface{}
    Body        interface{}
    Error       string
    ErrorData   map[string]string
    StreamBody
}

type StreamError struct {
    ErrorKey    string `json:"error"`
    ErrorValue  string `json:"reason"`
}

type StreamBody interface {
    SetBody(body interface{})
}

const (
    TYPE_REQUEST  = 1
    TYPE_RESPONSE = 2
)

func Shutup() {}

func NewStream() *Stream {
    return &Stream{
        Headers: util.Map(),
    }
}

func (this *Stream) SetHeader(key string, value interface{}) {
    switch value.(type) {
        case nil:
            delete(this.Headers, key)
        case int,
             bool,
             string:
            this.Headers[key] = util.String(value)
        default:
            panic("Unsupported value type '"+ _fmt.Sprintf("%T", value) +"' given!")
    }
}
func (this *Stream) GetHeader(key string) interface{} {
    if value, ok := this.Headers[key]; ok {
        return value
    }
    return nil
}
func (this *Stream) GetHeaderAll() map[string]interface{} {
    return this.Headers
}

func (this *Stream) GetBody() string {
    if this.Body == nil {
        return ""
    }
    return this.Body.(string)
}

func (this *Stream) GetBodyData(to interface{}) (interface{}, error) {
    if this.Error != "" {
        data, err := util.ParseBody(this.Body.(string), &StreamError{})
        if err != nil {
            return nil, err
        }
        return nil, _fmt.Errorf("Stream Error >> error: \"%s\", reason: \"%s\"",
            data.(*StreamError).ErrorKey,
            data.(*StreamError).ErrorValue,
        )
    }

    data, err := util.ParseBody(this.Body.(string), to)
    if err != nil {
        return nil, err
    }
    return data, nil
}

func (this *Stream) SetError(body string) {
    if body == "" {
        body = this.Body.(string)
    }
    data, err := util.ParseBody(body, &StreamError{})
    if data != nil && err == nil {
        var errorKey   = data.(*StreamError).ErrorKey
        var errorValue = data.(*StreamError).ErrorValue
        this.Error = _fmt.Sprintf("Stream Error >> error: \"%s\", reason: \"%s\"",
            errorKey,
            errorValue,
        )
        this.ErrorData = util.MapString()
        this.ErrorData["error"]  = errorKey
        this.ErrorData["reason"] = errorValue
    }
}

func (this *Stream) GetError() string {
    return this.Error
}
func (this *Stream) GetErrorValue(key string) string {
    return this.ErrorData[key]
}

// Get stream as string.
//
// @return (string)
func (this *Request) ToString() string {
    var ret string

    // add first line by type
    if this.Type == TYPE_REQUEST {
        ret = util.StringFormat("%s %s HTTP/%s\r\n", this.Method, this.Uri, this.HttpVersion)
    } else if this.Type == TYPE_RESPONSE {
        ret = util.StringFormat("HTTP/%s %d %s\r\n", this.HttpVersion, this.StatusCode, this.StatusText)
    }

    // add headers
    if this.Headers != nil { for key, value := range this.Headers {
        if key == "0" { // response only
            continue
        }
        if (value != nil) { // remove?
            ret += util.StringFormat("%s: %s\r\n", key, value)
        }
    }

    // add seperator
    ret += "\r\n"

    // add body
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
