package query

import (
    _fmt "fmt"
    _str "strings"
    _url "net/url"
)

import u "./../util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Query struct {
    Data        map[string]interface{}
    DataString  string
}

func Shutup() {}

func New(data map[string]interface{}) *Query {
    return &Query{
        Data: data,
        DataString: "",
    }
}

func (this *Query) ToString() string {
    if this.DataString != "" {
        return this.DataString
    }

    for key, value := range this.Data {
        if u.TypeReal(value) == "[]string" {
            value = _fmt.Sprintf("[\"%s\"]", _str.Join(value.([]string), "\",\""))
        }
        this.DataString += _fmt.Sprintf(
            "%s=%s&", _url.QueryEscape(key), _url.QueryEscape(u.String(value)))
    }

    if this.DataString != "" {
        // drop last "&"
        this.DataString = this.DataString[0 : len(this.DataString) - 1]
        this.DataString = _str.NewReplacer(
            "%5B", "[",
            "%5D", "]",
            "%2C", ",",
        ).Replace(this.DataString)
    }

    return this.DataString
}
