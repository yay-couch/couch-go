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
// @subpackage couch.uuid
// @uses       fmt, time, crypto/rand, strconv
// @author     Kerem Güneş <qeremy[at]gmail[dot]com>
package uuid

import (
    _fmt "fmt"
    _time "time"
    _rand "crypto/rand"
    _strc "strconv"
)

// @object couch.uuid.Uuid
type Uuid struct {
    Value interface{}
}

// UUID types (limits)
// @const int
const (
    RFC            = -1
    HEX_8          = 8
    HEX_32         = 32
    HEX_40         = 40
    TIMESTAMP      = 0
    TIMESTAMP_NANO = 1
)

func init() {
    _ = _time.ANSIC
    _ = _rand.Reader
}

func Shutup() {}

// Constructor.
//
// @param  value interface{}
// @return (*couch.uuid.Uuid)
func New(value interface{}) (*Uuid) {
    var this = &Uuid{}

    // auto-generate?
    if value == true {
        value = Generate(HEX_32)
    }
    this.SetValue(value)

    return this
}

// Set value.
//
// @param  value interface{}
// @return (void)
// @panics
func (this *Uuid) SetValue(value interface{}) {
    if value != nil {
        switch value.(type) {
            case int:
                this.Value = _strc.Itoa(value.(int))
            case string:
                this.Value = value
            default:
                panic("Only int or string values are accepted!")
        }
    }
}

// Get value.
//
// @return (interface{})
func (this *Uuid) GetValue() (interface{}) {
    if this.Value == nil {
        return nil
    }
    return this.Value.(string)
}

// Get value as string.
//
// @return (string)
func (this *Uuid) ToString() (string) {
    if this.Value == nil {
        return ""
    }

    return this.Value.(string)
}

// Generator.
//
// @param  limit int
// @return (string)
// @panics
func Generate(limit int) (string) {
    // unix epoch
    if limit == TIMESTAMP {
        return _fmt.Sprintf("%v", _time.Now().Unix())
    }
    // unix nano-epoch
    if limit == TIMESTAMP_NANO {
        return _fmt.Sprintf("%v", _time.Now().UnixNano())
    }

    // locals
    var (
        isRfc bool
        bytes []byte
    )

    switch limit {
        case RFC:
            isRfc = true
            bytes = make([]byte, 16)
        case HEX_8,
             HEX_32,
             HEX_40:
            bytes = make([]byte, limit / 2)
        default:
            panic("Unimplemented limit given, only -1,1,0 or 8,32,40 available!")
    }

    _, err := _rand.Read(bytes)
    if err != nil {
        panic(err)
    }

    if !isRfc {
        return _fmt.Sprintf("%x", bytes)
    } else {
        // use UUID/v4
        // https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_.28random.29
        bytes[6] = (bytes[6] | 0x40) & 0x4f
        bytes[8] = (bytes[8] | 0x80) & 0xbf
        return _fmt.Sprintf("%x-%x-%x-%x-%x",
            bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:])
    }
}
