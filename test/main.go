package main

// import test_uuid "./test_uuid"
// import test_client "./test_client"

import _couch  "./../src/couch"

import u "./../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

func main() {
    couch  := _couch.New(nil)
    client := _couch.NewClient(couch, nil)
    server := _couch.NewServer(client)
    // _dumps(server.Client)

    // _dump(server.Ping())

    var info = server.Info()
    _dumps(info)
    // _dumps(info["couchdb"])
    // _dumps(info["vendor"].(map[string]string)["name"])
    // _dumps(info["vendor.name"])

    // var value = u.Extract("couchdb", info)
    // var value = u.Extract("vendor.name", info)
    var value = u.Extract("vendor.name", info)
    _dumps(value == "Ubuntu")


    // test_uuid.TestAll()
    // test_client.TestAll()
}
