// Copyright 2015 Kerem Güneş
//    <http://qeremy.com>
//
// Apache License, Version 2.0
//    <http://www.apache.org/licenses/LICENSE-2.0>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// @package    couch
// @subpackage couch.query
// @uses       fmt, strings
// @uses       couch.util
// @author     Kerem Güneş <qeremy[at]gmail[dot]com>
package query

import (
    _fmt "fmt"
    _str "strings"
    _url "net/url"
)

import (
    "couch/util"
)

// @object couch.query.Query
type Query struct {
    Data        map[string]interface{}
    DataString  string
}

func Shutup() {}

// Constructor.
//
// @param  data map[string]interface{}
// @return (*couch.query.Query)
func New(data map[string]interface{}) *Query {
    if data == nil {
        data = util.Map()
    }

    return &Query{
        Data: data,
        DataString: "",
    }
}

// Set a param.
//
// @param  key    string
// @param  value interface{}
// @return (couch.query.Query)
func (this *Query) Set(key string, value interface{}) *Query {
    this.Data[key] = value

    return this
}

// Get a param.
//
// @param  key    string
// @return interface{}
func (this *Query) Get(key string) interface{} {
    if value, ok := this.Data[key]; ok {
        return value
    }

    return nil
}

// Skip as a shortcut.
//
// @param  value int
// @return (couch.query.Query)
func (this *Query) Skip(value int) *Query {
    this.Data["skip"] = value

    return this
}

// Limit as a shortcut.
//
// @param  value int
// @return (couch.query.Query)
func (this *Query) Limit(value int) *Query {
    this.Data["limit"] = value

    return this
}

// Get query data as map.
//
// @return (map[string]interface{})
func (this *Query) ToData() map[string]interface{} {
    return this.Data
}

// Get query data as string.
//
// @return (string)
func (this *Query) ToString() string {
    if this.DataString != "" {
        return this.DataString
    }

    for key, value := range this.Data {
        if util.TypeReal(value) == "[]string" {
            value = _fmt.Sprintf("[\"%s\"]", _str.Join(value.([]string), "\",\""))
        }
        this.DataString += _fmt.Sprintf(
            "%s=%s&", util.UrlEncode(key), util.UrlEncode(util.String(value)))
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
