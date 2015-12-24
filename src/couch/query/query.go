package query

import (
    _fmt "fmt"
    _str "strings"
    _url "net/url"
)

import (
    "couch/util"
)

type Query struct {
    Data        map[string]interface{}
    DataString  string
}

func Shutup() {}

func New(data map[string]interface{}) *Query {
    if data == nil {
        data = util.Map()
    }
    return &Query{
        Data: data,
        DataString: "",
    }
}

func (this *Query) Set(key string, value interface{}) *Query {
    this.Data[key] = value
    return this
}
func (this *Query) Get(key string) interface{} {
    if value, ok := this.Data[key]; ok {
        return value
    }
    return nil
}

func (this *Query) Skip(value int) *Query {
    this.Data["skip"] = value
    return this
}
func (this *Query) Limit(value int) *Query {
    this.Data["limit"] = value
    return this
}

func (this *Query) ToData() map[string]interface{} {
    return this.Data
}

func (this *Query) ToString() string {
    if this.DataString != "" {
        return this.DataString
    }

    for key, value := range this.Data {
        if util.TypeReal(value) == "[]string" {
            value = _fmt.Sprintf("[\"%s\"]", _str.Join(value.([]string), "\",\""))
        }
        this.DataString += _fmt.Sprintf(
            "%s=%s&", _url.QueryEscape(key), _url.QueryEscape(util.String(value)))
    }

    if this.DataString != "" {
        // drop last "&"
        this.DataString = this.DataString[0 : len(this.DataString) - 1]
        // purify some encoded stuff
        this.DataString = _str.NewReplacer(
            "%5B", "[",
            "%5D", "]",
            "%2C", ",",
        ).Replace(this.DataString)
    }

    return this.DataString
}
