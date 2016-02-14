// Copyright 2015 Kerem Güneş
//   <k-gun@mail.com>
//
// Apache License, Version 2.0
//   <http://www.apache.org/licenses/LICENSE-2.0>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// @package    couch
// @subpackage couch.util
// @uses       fmt, net/url, os, os/exec, regexp, strings, strconv,
//             encoding/json, encoding/base64, path/filepath
// @author     Kerem Güneş <k-gun@mail.com>
package util

import (
   _fmt "fmt"
   _url "net/url"
   _os "os"
   _ose "os/exec"
   _re "regexp"
   _str "strings"
   _strc "strconv"
   _json "encoding/json"
   _b64 "encoding/base64"
   _pathf "path/filepath"
)

func Shutup() {}

// Get short type.
//
// @param  args.. interface{}
// @return string
func Type(args... interface{}) (string) {
   return _str.Trim(TypeReal(args[0]), " *<>{}[]")
}

// Get real type.
//
// @param  args.. interface{}
// @return string
func TypeReal(args... interface{}) (string) {
   return _fmt.Sprintf("%T", args[0])
}

// Int converter..
//
// @param  input interface{}
// @return int
func Int(input interface{}) (int) {
   if number := Number(input, "int"); number != nil {
      return number.(int)
   }
   return 0
}

// UInt converter..
//
// @param  input interface{}
// @return uint
func UInt(input interface{}) (uint) {
   if number := Number(input, "uint"); number != nil {
      return number.(uint)
   }
   return 0
}

// Number converter..
//
// @param  input    interface{}
// @param  inputType string
// @return interface{}
func Number(input interface{}, inputType string) (interface{}) {
   if input != nil {
      number, err := _strc.Atoi(String(input))
      if err == nil {
         switch inputType {
            // signed
            case   "int": return int(number)
            case   "int8": return int8(number)
            case  "int16": return int16(number)
            case  "int32": return int32(number)
            case  "int64": return int64(number)
            // unsigned
            case   "uint": return uint(number)
            case  "uint8": return uint8(number)
            case "uint16": return uint16(number)
            case "uint32": return uint32(number)
            case "uint64": return uint64(number)
            // float
            case "float32": return float32(number)
            case "float64": return float64(number)
         }
      }
   }
   return nil
}

// Check empty.
//
// @param  input interface{}
// @return bool
func IsEmpty(input interface{}) (bool) {
   return (input == nil || input == "" || input == 0)
}

// Check empty & set default value.
//
// @param  input      interface{}
// @param  inputDefault interface{}
// @return interface{}
func IsEmptySet(input, inputDefault interface{}) (interface{}) {
   if IsEmpty(input) {
      input = inputDefault
   }
   return input
}

// Dump.
//
// @param  args... interface{}
// @return void
func Dump(args... interface{}) {
   _fmt.Println(args...)
}

// Dump as string.
//
// @param  args... interface{}
// @return void
func Dumps(args... interface{}) {
   var format string
   for _, arg := range args {
      _ = arg // silence..
      format += "%+v "
   }
   _fmt.Printf("%s\n", _fmt.Sprintf(format, args...))
}

// Dump as formatted string.
//
// @param  format  string
// @param  args... interface{}
// @return void
func Dumpf(format string, args... interface{}) {
   if format == "" {
      for _, arg := range args {
         _ = arg // silence..
         format += "%+v "
      }
   }
   _fmt.Printf("%s\n", _fmt.Sprintf(format, args...))
}

// Quote.
//
// @param  input string
// @return string
func Quote(input string) (string) {
   return _strc.Quote(input)
}

// Quote with encoding.
//
// @param  input string
// @return string
func QuoteEncode(input string) (string) {
   return _str.Replace(input, "\"", "%22", -1)
}

// Check param, init if nil
//
// @param  param map[string]interface{}
// @return map[string]interface{}
func Param(param map[string]interface{}) (map[string]interface{}) {
   if param == nil {
      param = Map()
   }
   return param
}

// Create param list like "key=>value" pairs
//
// @usage
//   x := util.ParamList(
//      "name", "kerem",
//      "old", 20,
//      // ...
//   )
// @param  argv... interface{}
// @return map[string]interface{}
// @panics
func ParamList(argv... interface{}) (map[string]interface{}) {
   argc := len(argv)
   if argc % 2 != 0 {
      panic("Wrong param count (key1, value1, key2, value2)!")
   }

   paramList := Map()
   // tricky?
   for i := 1; i < argc; i += 2 {
      if key, ok := argv[i - 1].(string); ok {
         paramList[key] = argv[i]
         continue
      }
      panic("Each param key must be string!")
   }

   return paramList
}

// Parse URL.
//
// @param  url string
// @return map[string]string
// @panics
func ParseUrl(url string) (map[string]string) {
   if url == "" {
      panic("No URL given!")
   }

   ret := MapString()
   ptr := "(?:(?P<Scheme>https?)://(?P<Host>[^:/]+))?" +
          "(?:\\:(?P<Port>\\d+))?(?P<Path>/[^?#]*)?"   +
          "(?:\\?(?P<Query>[^#]+))?"                   +
          "(?:\\??#(?P<Fragment>.*))?"

   re, _ := _re.Compile(ptr)
   if re == nil {
      return ret
   }

   match := re.FindStringSubmatch(url)
   for i, name := range re.SubexpNames() {
      if i != 0 { // pass re input
         ret[name] = match[i]
      }
   }

   return ret
}

// Parse query.
//
// @param  query string
// @return map[string]string
func ParseQuery(query string) (map[string]string) {
   ret := MapString()
   if tmps := _str.Split(query, "&"); tmps != nil {
      for _, tmp := range tmps {
         if t := _str.SplitN(tmp, "=", 2); len(t) == 2 {
            ret[t[0]] = t[1]
         }
      }
   }

   return ret
}

// Parse raw headers.
//
// @param  headers string
// @return map[string]string
func ParseHeaders(headers string) (map[string]string) {
   ret := MapString()
   if tmps := _str.Split(headers, "\r\n"); tmps != nil {
      // status line (HTTP/1.0 200 OK)
      ret["0"] = _shiftSliceString(&tmps)

      for _, tmp := range tmps {
         if t := _str.SplitN(tmp, ":", 2); len(t) == 2 {
            ret[_str.TrimSpace(t[0])] = _str.TrimSpace(t[1])
         }
      }
   }

   return ret
}

// Parse  body (JSON decode).
//
// @param  in  string
// @param  out interface{}
// @return interface{}, error
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

// Un-Parse  body (JSON encode).
//
// @param  in interface{}
// @return string, error
func UnparseBody(in interface{}) (string, error) {
   out, err := _json.Marshal(in)
   if err != nil {
      return "", _fmt.Errorf("JSON error: %s!", err)
   }

   return string(out), nil
}

// URL encode.
//
// @param  input string
// @return string
func UrlEncode(input string) (string) {
   return _url.QueryEscape(input)
}

// URL decode.
//
// @param  input string
// @return string
func UrlDecode(input string) (string) {
   input, err := _url.QueryUnescape(input)
   if err != nil {
      return ""
   }

   return input
}

// String converter.
//
// @param  input interface{}
// @return string
// @panics
func String(input interface{}) (string) {
   switch input.(type) {
      case int,
          bool,
          string:
         return _fmt.Sprintf("%v", input)
      default:
         inputType := _fmt.Sprintf("%T", input)
         // check numerics
         if StringSearch(inputType, "u?int(\\d+)?|float(32|64)") {
            return _fmt.Sprintf("%v", input)
         }
         panic("Unsupported input type '"+ inputType +"' given!")
   }
}

// String format.
//
// @param  format string
// @param  args... interface{}
// @return string
func StringFormat(format string, args... interface{}) (string) {
   return _fmt.Sprintf(format, args...)
}

// String search.
//
// @param  format string
// @param  search string
// @return bool
func StringSearch(input, search string) (bool) {
   re, _ := _re.Compile(search)
   if re == nil {
      return false
   }

   return ("" != re.FindString(input))
}

// Trim.
//
// @param  input string
// @param  chars string
// @return string
func Trim(input, chars string) (string) {
   return _str.Trim(input, chars)
}

// Get directory name.
//
// @param  path string
// @return string
func Dirname(path string) (string) {
   dirname, err := _pathf.Abs(path)
   if err != nil {
      return ""
   }

   return _pathf.Dir(dirname)
}

// Get base name.
//
// @param  path string
// @return string
func Basename(path string) (string) {
   return _pathf.Base(path)
}

// Check file exists.
//
// @param  file string
// @return bool
func FileExists(file string) (bool) {
   if _, err := _os.Stat(file); err == nil {
      return true
   }
   return false
}

// Get file size.
//
// @param file string
// @retur (int64)
func FileSize(file string) (int64) {
   if stat, err := _os.Stat(file); err == nil {
      return stat.Size()
   }
   return -1
}

// Get file info.
//
// @param file string
// @retur (map[string]string, error)
func FileInfo(file string) (map[string]string, error) {
   if !FileExists(file) {
      return nil, _fmt.Errorf("Given file does not exist! file: '%s'", file)
   }

   info := map[string]string{
      "mime": "",
      "charset": "",
      "name": "",
      "extension": "",
   }

   // add name/extension
   info["name"] = Basename(file)
   info["extension"] = _str.TrimLeft(_pathf.Ext(file), ".")

   // use built-in "file -i /file.txt" command
   out, err := _ose.Command("file", "-i", String(info["name"])).Output()
   if err != nil {
      return nil, _fmt.Errorf("EXEC error! %s", err)
   }

   tmp := _str.Split(_str.TrimSpace(string(out)), " ")
   if len(tmp) == 3 {
      mime := _str.TrimSpace(tmp[1])
      if i := _str.LastIndex(mime, ";"); i > -1 {
         mime = mime[:i]
      }

      // add mime/charset
      info["mime"] = mime
      info["charset"] = _str.Split(_str.TrimSpace(tmp[2]), "=")[1]
   }

   return info, nil
}

// Get file contents.
//
// @param  file string
// @return string, error
func FileGetContents(file string) (string, error) {
   fp, err := _os.Open(file)
   if err != nil {
      return "", _fmt.Errorf("FILE error! %s", err)
   }

   data := make([]byte, FileSize(file))
   if _, err := fp.Read(data); err != nil {
      return "", _fmt.Errorf("FILE error! %s", err)
   }

   return string(data), nil
}

// Base64 encode.
//
// @param  input string
// @return string, error
func Base64Encode(input string) (string) {
   return _b64.StdEncoding.EncodeToString([]byte(input))
}

// Base64 decode.
//
// @param  input string
// @return string
func Base64Decode(input string) (string) {
   data, err := _b64.StdEncoding.DecodeString(input)
   if err == nil {
      return string(data)
   }
   return ""
}

// Yes, Dig the digger!
//
// @param  key   string
// @param  object interface{}
// @return interface{}
func Dig(key string, object interface{}) (interface{}) {
   if object == nil {
      return nil
   }

   var keys = _str.Split(key, ".")
   key = _shiftSliceString(&keys)

   if len(keys) == 0 {
      // notation: util.Dig("x", map)
      switch object.(type) {
         case map[string]int:
            return object.(map[string]int)[key]
         case map[string]string:
            return object.(map[string]string)[key]
         case map[string]interface{}:
            return object.(map[string]interface{})[key]
         case []string:
            // check numeric key
            key, err := _strc.Atoi(key)
            if err == nil {
               return object.([]string)[key]
            }
         case []interface{}:
            // check numeric key
            key, err := _strc.Atoi(key)
            if err == nil {
               return object.([]interface{})[key]
            }
         case []map[string]interface{}:
            // check numeric key
            key, err := _strc.Atoi(key)
            if err == nil {
               return object.([]map[string]interface{})[key]
            }
         default:
            // @tmp for debugging or add more if needs?
            panic("Unimplemented type: "+ TypeReal(object))
      }
   } else {
      // notation: util.Dig("x.y.z", map)
      var keys = _str.Join(keys, ".") // @overwrite
      switch object.(type) {
         case map[string]int:
            return Dig(keys, object.(map[string]int)[key])
         case map[string]string:
            return Dig(keys, object.(map[string]string)[key])
         case map[string]interface{}:
            return Dig(keys, object.(map[string]interface{})[key])
         case []string:
            // check numeric key
            key, err := _strc.Atoi(key)
            if err == nil {
               return object.([]string)[key]
            }
         case []interface{}:
            // check numeric key
            key, err := _strc.Atoi(key)
            if err == nil {
               return object.([]interface{})[key]
            }
         case []map[string]interface{}:
            // check numeric key
            key, err := _strc.Atoi(key)
            if err == nil {
               return Dig(keys, object.([]map[string]interface{})[key])
            }
         default:
            // @tmp for debugging or add more if needs?
            panic("Unimplemented type: "+ TypeReal(object))
      }
   }

   return nil
}

// Dig for int value.
//
// @param  key   string
// @param  object interface{}
// @return int
func DigInt(key string, object interface{}) (int) {
   if value := Dig(key, object); value != nil {
      switch value := value.(type) {
         case int:
            return value
         case float32:
            return int(value)
         case float64:
            return int(value)
         case string:
            return Int(value)
      }
   }
   return 0
}

// Dig for uint value.
//
// @param  key   string
// @param  object interface{}
// @return uint
func DigUInt(key string, object interface{}) (uint) {
   return uint(DigInt(key, object))
}

// Dig for float value.
//
// @param  key   string
// @param  object interface{}
// @return float64
func DigFloat(key string, object interface{}) (float64) {
   if value := Dig(key, object); value != nil {
      switch value := value.(type) {
         case float64:
            return value
         case float32:
            return float64(value)
         case int:
            return float64(value)
      }
   }
   return 0.00
}

// Dig for string value.
//
// @param  key   string
// @param  object interface{}
// @return string
func DigString(key string, object interface{}) (string) {
   if value := Dig(key, object); value != nil {
      return value.(string)
   }
   return ""
}

// Dig for bool value.
//
// @param  key   string
// @param  object interface{}
// @return bool
func DigBool(key string, object interface{}) (bool) {
   if value := Dig(key, object); value != nil {
      return true
   }
   return false
}

// Dig for map value.
//
// @param  key   string
// @param  object interface{}
// @return map[string]interface{}
func DigMap(key string, object interface{}) (map[string]interface{}) {
   return Dig(key, object).(map[string]interface{})
}

// Dig for map list value.
//
// @param  key   string
// @param  object interface{}
// @return []map[string]interface{}
func DigMapList(key string, object interface{}) ([]map[string]interface{}) {
   return Dig(key, object).([]map[string]interface{})
}

// Dig for []int value.
//
// @param  key   string
// @param  object interface{}
// @return []int
func DigSliceInt(key string, object interface{}) ([]int) {
   slice := MapSliceInt(nil)
   for _, value := range Dig(key, object).([]interface{}) {
      slice = append(slice, value.(int))
   }
   return slice
}

// Dig for []string value.
//
// @param  key   string
// @param  object interface{}
// @return []string
func DigSliceString(key string, object interface{}) ([]string) {
   slice := MapSliceString(nil)
   for _, value := range Dig(key, object).([]interface{}) {
      slice = append(slice, value.(string))
   }
   return slice
}

// Map maker.
//
// @return map[string]interface{}
func Map() (map[string]interface{}) {
   return make(map[string]interface{})
}

// Int map maker.
//
// @return map[int]string
func MapInt() (map[int]string) {
   return make(map[int]string)
}

// String map maker.
//
// @return map[string]string
func MapString() (map[string]string) {
   return make(map[string]string)
}

// Int map map maker.
//
// @return map[int]map[string]interface{}
func MapMapInt() (map[int]map[string]interface{}) {
   return make(map[int]map[string]interface{})
}

// String map map maker.
//
// @return map[string]map[string]interface{}
func MapMapString() (map[string]map[string]interface{}) {
   return make(map[string]map[string]interface{})
}

// Map list maker.
//
// @param  length interface{}
// @return []map[string]interface{}
func MapList(length interface{}) ([]map[string]interface{}) {
   len := _length(length)
   if len != -1 {
      return make([]map[string]interface{}, len)
   }
   return []map[string]interface{}{}
}

// Map int list maker.
//
// @param  length interface{}
// @return []map[int]string
func MapListInt(length interface{}) ([]map[int]string) {
   len := _length(length)
   if len != -1 {
      return make([]map[int]string, len)
   }
   return []map[int]string{}
}

// Map string list maker.
//
// @param  length interface{}
// @return []map[string]string
func MapListString(length interface{}) ([]map[string]string) {
   len := _length(length)
   if len != -1 {
      return make([]map[string]string, len)
   }
   return []map[string]string{}
}

// Map int slice maker.
//
// @param  length interface{}
// @return []int
func MapSliceInt(length interface{}) ([]int) {
   len := _length(length)
   if len != -1 {
      return make([]int, len)
   }
   return []int{}
}

// Map string slice maker.
//
// @param  length interface{}
// @return []string
func MapSliceString(length interface{}) ([]string) {
   len := _length(length)
   if len != -1 {
      return make([]string, len)
   }
   return []string{}
}

// Detect length.
//
// @param  length interface{}
// @return int
// @private
func _length(length interface{}) (int) {
   switch length.(type) {
      case int:
         return length.(int)
      case []int:
         return len(length.([]int))
      case []string:
         return len(length.([]string))
      case []interface{}:
         return len(length.([]interface{}))
      // case:
         // @todo add more cases if needs
   }
   return -1
}

// Shift int slice.
//
// @param  slice *[]int
// @return int
// @private
func _shiftSliceInt(slice *[]int) (int) {
   ret := (*slice)[0]
   *slice = (*slice)[1 : len(*slice)]
   return ret
}

// Shift string slice.
//
// @param  slice *[]string
// @return string
// @private
func _shiftSliceString(slice *[]string) (string) {
   ret := (*slice)[0]
   *slice = (*slice)[1 : len(*slice)]
   return ret
}
