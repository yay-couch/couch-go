package main

import _couch "./../../couch"

import u "./../util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

func main() {
    couch    := _couch.New(nil)
    client   := _couch.NewClient(couch, "", "", nil)
    response := client.DoRequest("GET /", nil, "nil...", nil)
    // _dumpf("response >> %+v", response)

    // var body = response.GetBody()
    // _dumpf("response: len(%d) body(%+v)", len(body), body)

    // type Response struct {
    //     CouchDB string
    //     Uuid    string
    //     Version string
    //     Vendor  map[string]string
    // }
    // var bodyData = u.ParseBody(response.GetBody(), &Response{})
    // _dumps(bodyData)
    // _dumps(bodyData.(*Response).CouchDB)
    // _dumps(bodyData.(*Response).Uuid)
    // _dumps(bodyData.(*Response).Version)
    // _dumps(bodyData.(*Response).Vendor)
    // _dumps(bodyData.(*Response).Vendor["name"])
    // _dumps(bodyData.(*Response).Vendor["version"])
}
