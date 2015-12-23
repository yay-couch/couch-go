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
    ToString    func() string
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
        return nil, _fmt.Errorf("Stream Error\n   >> error: \"%s\", reason: \"%s\"",
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
        this.ErrorData = make(map[string]string)
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
