package main

// import _couch  "./../src/couch"

// import test_uuid "./test_uuid"
// import test_client "./test_client"
// import test_server "./test_server"
import test_database "./test_database"


import u "./../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

func main() {
    // test_database.TestPing()
    // test_database.TestInfo()
    // test_database.TestCreate()
    // test_database.TestRemove()
    // test_database.TestReplicate()
    test_database.TestGetDocument()
}
