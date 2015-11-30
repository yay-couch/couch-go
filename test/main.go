package main

// import _couch  "./../src/couch"

// import test_uuid "./test_uuid"
// import test_client "./test_client"
import test_server "./test_server"


import u "./../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

func main() {
    // test_uuid.TestAll()
    // test_client.TestAll()

    // test_server.TestInfo()
    // test_server.TestVersion()
    // test_server.TestGetActiveTasks()
    // test_server.TestGetAllDatabases()
    // test_server.TestGetDatabaseUpdates()
    test_server.TestGetLogs()
    // test_server.TestGetStats()
}
