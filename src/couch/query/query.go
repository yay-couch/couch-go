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
    return &Query{data, ""}
}

func (query *Query) ToString() string {
    if query.DataString != "" {
        return query.DataString
    }

    for key, value := range query.Data {
        if u.TypeReal(value) == "[]string" {
            value = _fmt.Sprintf("[\"%s\"]", _str.Join(value.([]string), "\",\""))
        }
        query.DataString += _fmt.Sprintf(
            "%s=%s&", _url.QueryEscape(key), _url.QueryEscape(u.String(value)))
    }

    if query.DataString != "" {
        // drop last "&"
        query.DataString = query.DataString[0 : len(query.DataString) - 1]
        query.DataString = _str.NewReplacer(
            "%5B", "[",
            "%5D", "]",
            "%2C", ",",
        ).Replace(query.DataString)
    }

    return query.DataString
}
