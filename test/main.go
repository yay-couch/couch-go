package main

import test_uuid "./test_uuid"
// import test_client "./test_client"

import u "./../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

func main() {
    test_uuid.TestAll()
    // test_client.TestAll()
}
