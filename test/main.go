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
    // test_database.TestGetDocument()
    // test_database.TestGetDocumentAll()
    // test_database.TestCreateDocument()
    // test_database.TestCreateDocumentAll()
    // test_database.TestUpdateDocument()
    // test_database.TestUpdateDocumentAll()
    // test_database.TestDeleteDocument()
    // test_database.TestDeleteDocumentAll()
    // test_database.TestGetChanges()
    // test_database.TestCompact()
    // test_database.TestEnsureFullCommit()
    // test_database.TestViewCleanup()
    // test_database.TestViewTemp()
    // test_database.TestGetSecurity()
    // test_database.TestSetSecurity()
    // test_database.TestPurge()
    // test_database.TestGetMissingRevisions()
    test_database.TestGetMissingRevisionsDiff()
}
