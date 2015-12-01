package stream

import (
    _fmt "fmt"
)

import u "./../../util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Stream struct {
    Type        uint8
    HttpVersion string
    Headers     map[string]interface{}
    Body        interface{}
    Error       bool
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

func New() *Stream {
    return &Stream{
        Headers: make(map[string]interface{}),
    }
}

func (this *Stream) SetHeader(key string, value interface{}) {
    _checkStreamHeaders(this)
    switch value.(type) {
        case nil:
            delete(this.Headers, key)
        case int,
             string:
            this.Headers[key] = u.ToString(value);
        default:
            panic("Unsupported value type '"+ u.Type(value) +"' given!");
    }
}
func (this *Stream) GetHeader(key string) interface{} {
    _checkStreamHeaders(this)
    if value, ok := this.Headers[key]; ok {
        return value
    }
    return nil
}
func (this *Stream) GetHeaderAll() map[string]interface{} {
    _checkStreamHeaders(this)
    return this.Headers
}

func (this *Stream) GetBody() string {
    if this.Body == nil {
        return ""
    }
    return this.Body.(string)
}

func (this *Stream) GetData(to interface{}) (interface{}, error) {
    if this.Error == true {
        data, err := u.ParseBody(this.Body.(string), &StreamError{})
        if err != nil {
            return nil, err
        }
        return nil, _fmt.Errorf("Stream Error: %s, %s",
            data.(*StreamError).ErrorKey,
            data.(*StreamError).ErrorValue,
        )
    }

    data, err := u.ParseBody(this.Body.(string), to)
    if err != nil {
        return nil, err
    }
    return data, nil
}

func (this *Stream) ToString() string {
    // @todo
    return ""
}

func _checkStreamHeaders(stream *Stream) {
    if stream.Headers == nil {
        stream.Headers = make(map[string]interface{})
    }
}
