package database

import (
    "./../../src/couch"
    "./../../src/couch/util"
)

var (
    DEBUG  = true
    DBNAME = "foo"
)

var (
    Couch    *couch.Couch
    Client   *couch.Client
    Database *couch.Database
)

func init() {
    Couch    = couch.New(nil, DEBUG)
    Client   = couch.NewClient(Couch)
    Database = couch.NewDatabase(Client, DBNAME);
}

/**
 * TestAll
 */
func TestAll() {}

/**
 * TestPing
 */
func TestPing() {
    util.Dumpf("Database Ping >> %v", Database.Ping())
}

/**
 * TestInfo
 */
func TestInfo() {
    data, err := Database.Info()
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Info >> %+v", data)
    util.Dumpf("Database Info >> db_name: %s", data["db_name"])
    for key, value := range data {
        util.Dumpf("Database Info >> %s: %v", key, value)
    }
}

/**
 * TestCreate
 */
func TestCreate() {
    util.Dumpf("Database Create >> %v", Database.Create())
    // error?
    // if err := client.GetResponse().GetError(); err != "" {
    //     util.Dumpf("Response Status: %s", client.GetResponse().GetStatus())
    //     util.Dumpf("Response Body  : %s", client.GetResponse().GetBody())
    //     panic(err)
    // }
}

/**
 * TestRemove
 */
func TestRemove() {
    util.Dumpf("Database Remove >> %v", Database.Remove())
    // error?
    // if err := client.GetResponse().GetError(); err != "" {
    //     util.Dumpf("Response Status: %s", client.GetResponse().GetStatus())
    //     util.Dumpf("Response Body  : %s", client.GetResponse().GetBody())
    //     panic(err)
    // }
}

/**
 * TestReplicate
 */
func TestReplicate() {
    data, err := Database.Replicate("foo_replicate", true)
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Replicate >> %+v", data)
    util.Dumpf("Database Replicate >> ok: %v", data["ok"])
    util.Dumpf("Database Replicate >> history.0: %v", util.Dig("0", data["history"]))
    util.Dumpf("Database Replicate >> history.0.start_time: %s", util.Dig("0.start_time", data["history"]))
}

/**
 * TestGetDocument
 */
func TestGetDocument() {
    data, err := Database.GetDocument("5db345a5f26484352ea5d813180031fb")
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Document >> %+v", data)
    util.Dumpf("Database Document >> id: %s", data["id"])
    util.Dumpf("Database Document >> key: %s", data["key"])
    util.Dumpf("Database Document >> value.rev: %s", util.Dig("value.rev", data))
    util.Dumpf("Database Document >> doc: %+v", data["doc"])
    // util.Dumpf("Database Document >> doc._id: %s", util.Dig("doc._id", data))
    // util.Dumpf("Database Document >> doc._rev: %s", util.Dig("doc._rev", data))
    // or
    for key, value := range util.DigMap("doc", data) {
        util.Dumpf("Database Document >> doc.%s: %v", key, value)
    }
}

/**
 * TestGetDocumentAll
 */
func TestGetDocumentAll() {
    data, err := Database.GetDocumentAll(nil, nil)
    // data, err := Database.GetDocumentAll(nil, []string{"5db345a5f26484352ea5d813180031fb"})
    if err != nil {
        panic(err)
    }
    // util.Dumpf("Database Document All >> %+v", data)
    util.Dumpf("Database Document All >> offset: %d", data["offset"])
    util.Dumpf("Database Document All >> total_rows: %d", data["total_rows"])
    // // util.Dumpf("Database Document All >> rows: %+v", data["rows"])
    util.Dumpf("Database Document All >> rows.0: %+v", util.Dig("rows.0", data))
    util.Dumpf("Database Document All >> rows.0.id: %s", util.Dig("rows.0.id", data))
    util.Dumpf("Database Document All >> rows.0.key: %s", util.Dig("rows.0.key", data))
    util.Dumpf("Database Document All >> rows.0.value.rev: %s", util.Dig("rows.0.value.rev", data))
    util.Dumpf("Database Document All >> rows.0.doc.name: %s", util.Dig("rows.0.doc.name", data))
}

/**
 * TestCreateDocument
 */
func TestCreateDocument() {
    data, err := Database.CreateDocument(map[string]interface{}{
        "name": "kerem", "type": "tmp",
    })
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Create Document >> %+v", data)
    util.Dumpf("Database Create Document >> doc.ok: %v", data["ok"])
    util.Dumpf("Database Create Document >> doc.id: %s", data["id"])
    util.Dumpf("Database Create Document >> doc.rev: %s", data["rev"])
    // or
    for key, value := range data {
        util.Dumpf("Database Create Document >> doc.%s: %v", key, value)
    }
}

/**
 * TestCreateDocumentAll
 */
func TestCreateDocumentAll() {
    data, err := Database.CreateDocumentAll([]interface{}{
        0: map[string]interface{}{"name": "kerem", "type": "tmp"},
        1: map[string]interface{}{"name": "murat", "type": "tmp"},
    })
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Create Document All >> %+v", data)
    util.Dumpf("Database Create Document All >> doc.0.ok: %v", util.Dig("0.ok", data))
    util.Dumpf("Database Create Document All >> doc.0.id: %s", util.Dig("0.id", data))
    util.Dumpf("Database Create Document All >> doc.0.rev: %s", util.Dig("0.rev", data))
    // or
    for i, doc := range data {
        for key, value := range doc {
            util.Dumpf("Database Create Document All >> doc.%d.%s: %v", i, key, value)
        }
    }
}

/**
 * TestUpdateDocument
 */
func TestUpdateDocument() {
    data, err := Database.UpdateDocument(map[string]interface{}{
        "name": "kerem 3", "type": "tmp",
            "_id": "7ee9cdd673b109e030cec8c6f10020f7",
            "_rev": "2-BÖ!", // give correct rev!
    })
    if err != nil {
        panic(err)
    }
    // check ok's
    if data["ok"] == nil {
        panic("Halt! error: "+ data["error"].(string) +", reason: "+ data["reason"].(string))
    }
    util.Dumpf("Database Update Document >> %+v", data)
    util.Dumpf("Database Update Document >> doc.ok: %v", data["ok"])
    util.Dumpf("Database Update Document >> doc.id: %s", data["id"])
    util.Dumpf("Database Update Document >> doc.rev: %s", data["rev"])
    // or
    for key, value := range data {
        util.Dumpf("Database Update Document >> doc.%s: %v", key, value)
    }
}

/**
 * TestUpdateDocumentAll
 */
func TestUpdateDocumentAll() {
    data, err := Database.UpdateDocumentAll([]interface{}{
        0: map[string]interface{}{"name": "kerem 2", "type": "tmp",
            "_id": "7ee9cdd673b109e030cec8c6f10020f7", "_rev": "1-3c92d3e67136c8b206d90ea37a3ee76d"},
        1: map[string]interface{}{"name": "murat 2", "type": "tmp",
            "_id": "7ee9cdd673b109e030cec8c6f1002cc1", "_rev": "1-09e886345b525e53892815baff169f03"},
    })
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Update Document All >> %+v", data)
    util.Dumpf("Database Update Document All >> doc.0.ok: %v", util.Dig("0.ok", data))
    util.Dumpf("Database Update Document All >> doc.0.id: %s", util.Dig("0.id", data))
    util.Dumpf("Database Update Document All >> doc.0.rev: %s", util.Dig("0.rev", data))
    // or
    for i, doc := range data {
        // check ok's
        if doc["ok"] == nil {
            util.Dumpf("Halt! error: doc.%d > %s reason: %s", i, doc["error"], doc["reason"])
        }
        for key, value := range doc {
            util.Dumpf("Database Update Document All >> doc.%d.%s: %v", i, key, value)
        }
    }
}

/**
 * TestDeleteDocument
 */
func TestDeleteDocument() {
    data, err := Database.DeleteDocument(map[string]interface{}{
        "_id": "7ee9cdd673b109e030cec8c6f100322b",
        "_rev": "1-3c92d3e67136c8b206d90ea37a3ee76d",
    })
    if err != nil {
        panic(err)
    }
    // check ok's
    if data["ok"] == nil {
        panic("Halt! error: "+ data["error"].(string) +", reason: "+ data["reason"].(string))
    }
    util.Dumpf("Database Delete Document >> %+v", data)
}

/**
 * TestDeleteDocumentAll
 */
func TestDeleteDocumentAll() {
    data, err := Database.DeleteDocumentAll([]interface{}{
        0: map[string]interface{}{"_id": "7ee9cdd673b109e030cec8c6f10020f7",
            "_rev": "3-1cb8867864cd6bb80361fed42c719897"},
        1: map[string]interface{}{"_id": "7ee9cdd673b109e030cec8c6f1002cc1",
            "_rev": "2-38a4374c91a4ae8d666d6c1fe13dd916"},
    })
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Delete Document All >> %+v", data)
    // or
    for i, doc := range data {
        // check ok's
        if doc["ok"] == nil {
            util.Dumpf("Halt! error: doc.%d > %s reason: %s", i, doc["error"], doc["reason"])
        }
        for key, value := range doc {
            util.Dumpf("Database Update Document All >> doc.%d.%s: %v", i, key, value)
        }
    }
}

/**
 * TestGetChanges
 */
func TestGetChanges() {
    data, err := Database.GetChanges(nil, nil)
    if err != nil {
        panic(err)
    }
    // util.Dumpf("Database Changes >> %+v", data)
    util.Dumpf("Database Changes >> last_seq: %v", data["last_seq"])
    for i, result := range util.DigMapList("results", data) {
        util.Dumps(i,result)
    }
}

/**
 * TestCompact
 */
func TestCompact() {
    data, err := Database.Compact("")
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Compact >> ok: %v", data["ok"])
}

/**
 * TestEnsureFullCommit
 */
func TestEnsureFullCommit() {
    data, err := Database.EnsureFullCommit()
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Ensure Full Commit >> ok: %v", data["ok"])
    util.Dumpf("Database Ensure Full Commit >> instance_start_time: %s", data["instance_start_time"])
}

/**
 * TestViewCleanup
 */
func TestViewCleanup() {
    data, err := Database.ViewCleanup()
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database View Cleanup >> ok: %v", data["ok"])
}

/**
 * TestViewTemp
 */
func TestViewTemp() {
    var map_ = "function(doc){if(doc.type=='tmp') emit(null,doc)}"
    data, err := Database.ViewTemp(map_, "")
    if err != nil {
        panic(err)
    }
    // util.Dumpf("Database View Temp >> %+v", data)
    util.Dumpf("Database View Temp >> offset: %d", data["offset"])
    util.Dumpf("Database View Temp >> total_rows: %d", data["total_rows"])
    // // util.Dumpf("Database View Temp >> rows: %+v", data["rows"])
    util.Dumpf("Database View Temp >> rows.0: %+v", util.Dig("rows.0", data))
    util.Dumpf("Database View Temp >> rows.0.id: %s", util.Dig("rows.0.id", data))
    util.Dumpf("Database View Temp >> rows.0.key: %s", util.Dig("rows.0.key", data))
    util.Dumpf("Database View Temp >> rows.0.value: %+v", util.Dig("rows.0.value", data))
    util.Dumpf("Database View Temp >> rows.0.value._id: %s", util.Dig("rows.0.value._id", data))
    util.Dumpf("Database View Temp >> rows.0.value._rev: %s", util.Dig("rows.0.value._rev", data))
    util.Dumpf("Database View Temp >> rows.0.value.type: %s", util.Dig("rows.0.value.type", data))
    util.Dumpf("Database View Temp >> rows.0.value.name: %s", util.Dig("rows.0.value.name", data))
}

/**
 * TestGetSecurity
 */
func TestGetSecurity() {
    data, err := Database.GetSecurity()
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Get Security >> %+v", data)
    util.Dumpf("Database Get Security >> admins %+v", data["admins"])
    util.Dumpf("Database Get Security >> admins.names.0 %s", util.Dig("admins.names.0", data))
    util.Dumpf("Database Get Security >> admins.roles.0 %s", util.Dig("admins.roles.0", data))
}

/**
 * TestSetSecurity
 */
func TestSetSecurity() {
    var admins = map[string]interface{}{
        "names": []string{"superuser"},
        "roles": []string{"admins"},
    }
    var members = map[string]interface{}{
        "names": []string{"user1","user2"},
        "roles": []string{"developers"},
    }
    data, err := Database.SetSecurity(admins, members)
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Set Security >> ok: %v", data["ok"])
}

/**
 * TestPurge
 */
func TestPurge() {
    data, err := Database.Purge(map[string]interface{}{
        "667b0208441066a0954717b50c0008a9": []string{
            "5-dd1a3738fcbd759ed744f7971fe94332",
        },
    })
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Purge >> ok: %v", data)
    util.Dumpf("Database Purge >> ok: %d", data["purge_seq"])
}

/**
 * TestGetMissingRevisions
 */
func TestGetMissingRevisions() {
    data, err := Database.GetMissingRevisions(map[string]interface{}{
        "667b0208441066a0954717b50c0008a9": []string{
            "5-dd1a3738fcbd759ed744f7971fe94332",
        },
    })
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Missing Revisions >> %v", data)
    util.Dumpf("Database Missing Revisions >> missing_revs: %v", data["missing_revs"])
    util.Dumpf("Database Missing Revisions >> missing_revs.667b0208441066a0954717b50c0008a9: %s",
        util.Dig("missing_revs.667b0208441066a0954717b50c0008a9", data))
    util.Dumpf("Database Missing Revisions >> missing_revs.0: %s",
        util.Dig("missing_revs.667b0208441066a0954717b50c0008a9.0", data))
}

/**
 * TestGetMissingRevisionsDiff
 */
func TestGetMissingRevisionsDiff() {
    data, err := Database.GetMissingRevisionsDiff(map[string]interface{}{
        "667b0208441066a0954717b50c0008a9": []string{
            "5-dd1a3738fcbd759ed744f7971fe94332",
        },
    })
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Missing Revisions Diff >> %v", data)
    util.Dumpf("Database Missing Revisions Diff >> 667b0208441066a0954717b50c0008a9: %v",
        util.Dig("667b0208441066a0954717b50c0008a9", data))
    util.Dumpf("Database Missing Revisions Diff >> 667b0208441066a0954717b50c0008a9.missing: %v",
        util.Dig("667b0208441066a0954717b50c0008a9.missing", data))
    util.Dumpf("Database Missing Revisions Diff >> 667b0208441066a0954717b50c0008a9.missing.0: %s",
        util.Dig("667b0208441066a0954717b50c0008a9.missing.0", data))
}

/**
 * TestGetRevisionLimit
 */
func TestGetRevisionLimit() {
    data, err := Database.GetRevisionLimit()
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Get Revision Limit >> %d", data)
}

/**
 * TestSetRevisionLimit
 */
func TestSetRevisionLimit() {
    data, err := Database.SetRevisionLimit(1000)
    if err != nil {
        panic(err)
    }
    util.Dumpf("Database Set Revision Limit >> ok: %v", data["ok"])
}
