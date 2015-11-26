package uuid

import (
    _fmt "fmt"
    _time "time"
    _rand "crypto/rand"
    _strc "strconv"
)

const (
    RFC       = -1
    HEX_8     = 8
    HEX_32    = 32
    HEX_40    = 40
    TIMESTAMP = 0
)

type Uuid struct {
    Value interface{}
}

func Shutup() {
    _ = _time.ANSIC
    _ = _rand.Reader
}

func New(value interface{}) *Uuid {
    if value == nil {
        value = Generate(HEX_32)
    }

    var this = &Uuid{}

    this.SetValue(value)

    return this
}

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
func (this *Uuid) GetValue() interface{} {
    if this.Value == nil {
        return nil
    }
    return this.Value.(string)
}
func (this *Uuid) ToString() string {
    if this.Value == nil {
        return ""
    }
    return this.Value.(string)
}

func Generate(limit int) string {
    if limit == TIMESTAMP {
        return _fmt.Sprintf("%v", _time.Now().Unix())
    }

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
            panic("Unimplemented limit given, only -1|0 or 8|32|40 available!")
    }

    _, err := _rand.Read(bytes)
    if err != nil {
        panic(err)
    }

    if !isRfc {
        return _fmt.Sprintf("%x", bytes)
    } else {
        // UUID/v4 >> https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_.28random.29
        bytes[6] = (bytes[6] | 0x40) & 0x4f
        bytes[8] = (bytes[8] | 0x80) & 0xbf
        return _fmt.Sprintf("%x-%x-%x-%x-%x",
            bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:])
    }
}
