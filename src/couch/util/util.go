package util

import (
    _fmt "fmt"
    _str "strings"
    _strc "strconv"
    _json "encoding/json"
    _rex "regexp"
)

type Nil struct {}

func Foo(a interface{}) {

}

func Type(args ...interface{}) string {
    return _str.Trim(TypeReal(args[0]), " *<>{}[]")
}
func TypeReal(args ...interface{}) string {
    return _fmt.Sprintf("%T", args[0])
}

func ToString(args ...interface{}) string {
    return _fmt.Sprintf("%v", args[0]);
}
func String(input interface{}) string {
    return _fmt.Sprintf("%v", input);
}
func Int(input interface{}) int {
    result, err := _strc.Atoi(input.(string))
    if err != nil {
        return int(result)
    }
    return 0
}
func Number(input interface{}, inputType string) interface{} {
    result, err := _strc.Atoi(input.(string))
    if err != nil {
        return nil
    }
    switch inputType {
        // signed
        case    "int": return int(result)
        case   "int8": return int8(result)
        case  "int16": return int16(result)
        case  "int32": return int32(result)
        case  "int64": return int64(result)
        // unsigned
        case   "uint": return uint(result)
        case  "uint8": return uint8(result)
        case "uint16": return uint16(result)
        case "uint32": return uint32(result)
        case "uint64": return uint64(result)
    }
    return 0
}


func IsEmpty(input interface{}) bool {
    if input == nil || input == "" || input == 0 {
        return true
    }
    return false
}
func IsEmptySet(input interface{}, inputDefault interface{}/*, inputType string*/) interface{} {
    if IsEmpty(input) {
        input = inputDefault
        // switch inputType {
        //     case "string":
        //         input = ToString(inputDefault)
        //     default:
        //         panic("Unimplemeted type '"+ inputType +"' given!")
        // }
    }
    return input
}

func Dump(args ...interface{}) {
    _fmt.Println(args...)
}
func Dumps(args ...interface{}) {
    var format string
    for _, arg := range args {
        _ = arg // silence..
        format += "%+v "
    }
    _fmt.Printf("%s\n", _fmt.Sprintf(format, args...))
}
func Dumpf(format string, args ...interface{}) {
    if format == "" {
        for _, arg := range args {
            _ = arg // silence..
            format += "%+v "
        }
    }
    _fmt.Printf("%s\n", _fmt.Sprintf(format, args...))
}

func Quote(input string) {
    input = _str.Replace(input, "\"", "%22", -1)
}

// parsers
func ParseUrl(url string) map[string]string {
    if url == "" {
        panic("No URL given!")
    }
    var result = make(map[string]string)
    var pattern = "(?:(?P<Scheme>https?)://(?P<Host>[^:/]+))?" +
                  "(?:\\:(?P<Port>\\d+))?(?P<Path>/[^?#]*)?"   +
                  "(?:\\?(?P<Query>[^#]+))?"                   +
                  "(?:\\??#(?P<Fragment>.*))?"
    re, _ := _rex.Compile(pattern)
    if re == nil {
        return result
    }
    var match = re.FindStringSubmatch(url)
    for i, name := range re.SubexpNames() {
        if i != 0 {
            result[name] = match[i]
        }
    }
    return result
}

func ParseHeaders(headers string) map[string]string {
    var result = make(map[string]string)
    if tmps := _str.Split(headers, "\r\n"); tmps != nil {
        for _, tmp := range tmps {
            var tmp = _str.SplitN(tmp, ":", 2)
            // request | response check?
            if len(tmp) == 1 {
                // status line >> HTTP/1.0 200 OK
                result["0"] = tmp[0]
                continue
            }
            var key, value =
                _str.TrimSpace(tmp[0]),
                _str.TrimSpace(tmp[1]);
            result[key] = value
        }
    }
    return result
}

func ParseBody(in string, out interface{}) (interface{}, error) {
    // simply prevent useless unmarshal error
    if in == "" {
        in = `null`
    }
    err := _json.Unmarshal([]byte(in), &out)
    if err != nil {
        return nil, _fmt.Errorf("JSON error: %s!", err)
    }
    return out, nil
}

func Dig(key string, object interface{}) interface{} {
    var keys = _str.Split(key, ".")
    key  = ArrayShiftString(&keys)
    if len(keys) != 0 {
        // nö!
        // if value, ok := object.(map[string]interface{})[key]; ok {
        //     return Dig(_str.Join(keys, "."), value)
        // }

        // @overwrite
        var keys = _str.Join(keys, ".")

        // add more if needs
        switch object.(type) {
            case map[string]int:
                return Dig(keys, object.(map[string]int)[key])
            case map[string]string:
                return Dig(keys, object.(map[string]string)[key])
            case map[string]interface{}:
                return Dig(keys, object.(map[string]interface{})[key])
            default:
                // panic?
        }
    } else {
        // add more if needs
        switch object.(type) {
            case map[string]int:
                return object.(map[string]int)[key]
            case map[string]string:
                return object.(map[string]string)[key]
            case map[string]interface{}:
                return object.(map[string]interface{})[key]
            default:
                // panic?
        }
    }

    return nil
}

// @todo add more if needs
// tip kontrolu yapabilirsin ArrayShift >> String seysnden kurtulmak icun..
func ArrayShiftString(slice *[]string) string {
    var value = (*slice)[0]
    *slice = (*slice)[1 : len(*slice)]
    return value
}
