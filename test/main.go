package main

// import _couch "./../src/couch"

// import "encoding/hex"
import (
    "fmt"
    "time"
    "math/rand"
)

import u "./../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

func main() {
    rand.Seed( time.Now().UnixNano())
    // var bytes int

    s := "mö!"
    // for i:= 0 ; i < 8 ; i++{
    //     bytes = rand.Intn(10)+1
    //     _dumpf("%v", bytes)
    //     s += u.String(bytes)
    // }
    fmt.Printf(">> %x\n", s)

    // _dumpf(">> %s", hex.EncodeToString([]byte(s)))
    // _dumpf(">> %x", s)
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
