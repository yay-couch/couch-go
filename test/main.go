package main

// import _couch "./../src/couch"

// import "encoding/hex"
import (
    // "fmt"
    // "time"
    // "math"
    // "math/rand"
    // "crypto/rand"
    // "strconv"
)

import _uuid "./../src/couch/uuid"

import u "./../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

func main() {
    u1 := _uuid.Generate(_uuid.HEX_8)
    _dumpf("%s", u1)
    _dumpf("%d", len(u1))

    // u2 := _uuid.New("123")
    u2 := _uuid.New(123)
    _dump(u2.Value)
    _dump(u2.GetValue())
    _dump(u2.ToString())

    // couch    := _couch.New(nil)
    // client   := _couch.NewClient(couch, nil)
    // // client   := _couch.NewClient(couch, map[string]interface{}{
    // //     "Host": "127.0.0.1",
    // // })
    // response := client.DoRequest("GET /", nil, "", nil)
    // _dumpf("response >> %+v", response)

    // var responseBody = response.GetBody()
    // _dumpf("response: len(%d) body(%+v)", len(responseBody), responseBody)

    // type Response struct {
    //     CouchDB string
    //     Uuid    string
    //     Version string
    //     Vendor  map[string]string
    // }
    // body, err := u.ParseBody(response.GetBody(), &Response{})
    // if err != nil {
    //     _dumps(err)
    //     return
    // }
    // _dumps(body)
    // _dumps(body.(*Response).CouchDB)
    // _dumps(body.(*Response).Uuid)
    // _dumps(body.(*Response).Version)
    // _dumps(body.(*Response).Vendor)
    // _dumps(body.(*Response).Vendor["name"])
    // _dumps(body.(*Response).Vendor["version"])
}
