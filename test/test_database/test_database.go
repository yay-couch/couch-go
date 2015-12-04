package test_database

import _couch    "./../../src/couch"
import _client   "./../../src/couch/client"
import _database "./../../src/couch/database"
import _document "./../../src/couch/document"

import u "./../../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

var (
    DEBUG  = true
    DBNAME = "foo"
)

var (
    couch    *_couch.Couch
    client   *_client.Client
    database *_database.Database
)

func init() {
    _document.Shutup()

    couch    = _couch.New(nil, DEBUG)
    client   = _couch.NewClient(couch, nil)
    database = _couch.NewDatabase(client, DBNAME);
}

/**
 * TestAll
 */
func TestAll() {}

/**
 * TestPing
 */
func TestPing() {
    _dumpf("Database Ping >> %v", database.Ping())
}

/**
 * TestInfo
 */
func TestInfo() {
    data, err := database.Info()
    if err != nil {
        panic(err)
    }
    _dumpf("Database Info >> %+v", data)
    _dumpf("Database Info >> db_name: %s", data["db_name"])
    for key, value := range data {
        _dumpf("Database Info >> %s: %v", key, value)
    }
}

/**
 * TestCreate
 */
func TestCreate() {
    _dumpf("Database Create >> %v", database.Create())
    // error?
    // if err := client.GetResponse().GetError(); err != "" {
    //     _dumpf("Response Status: %s", client.GetResponse().GetStatus())
    //     _dumpf("Response Body  : %s", client.GetResponse().GetBody())
    //     panic(err)
    // }
}

/**
 * TestRemove
 */
func TestRemove() {
    _dumpf("Database Remove >> %v", database.Remove())
    // error?
    // if err := client.GetResponse().GetError(); err != "" {
    //     _dumpf("Response Status: %s", client.GetResponse().GetStatus())
    //     _dumpf("Response Body  : %s", client.GetResponse().GetBody())
    //     panic(err)
    // }
}

/**
 * TestReplicate
 */
func TestReplicate() {
    data, err := database.Replicate("foo_replicate", true)
    if err != nil {
        panic(err)
    }
    _dumpf("Database Replicate >> %+v", data)
    _dumpf("Database Replicate >> ok: %v", data["ok"])
    _dumpf("Database Replicate >> history.0: %v", u.Dig("0", data["history"]))
    _dumpf("Database Replicate >> history.0.start_time: %s", u.Dig("0.start_time", data["history"]))
}

/**
 * TestGetDocument
 */
func TestGetDocument() {
    data, err := database.GetDocument("5db345a5f26484352ea5d813180031fb")
    if err != nil {
        panic(err)
    }
    _dumpf("Database Document >> %+v", data)
    _dumpf("Database Document >> id: %s", data["id"])
    _dumpf("Database Document >> key: %s", data["key"])
    _dumpf("Database Document >> value.rev: %s", u.Dig("value.rev", data))
    _dumpf("Database Document >> doc: %+v", data["doc"])
    // _dumpf("Database Document >> doc._id: %s", u.Dig("doc._id", data))
    // _dumpf("Database Document >> doc._rev: %s", u.Dig("doc._rev", data))
    // or
    for key, value := range u.DigMap("doc", data) {
        _dumpf("Database Document >> doc.%s: %v", key, value)
    }
}

/**
 * TestGetDocumentAll
 */
func TestGetDocumentAll() {
    // data, err := database.GetDocumentAll(nil, nil)
    data, err := database.GetDocumentAll(nil, []string{"5db345a5f26484352ea5d813180031fb"})
    if err != nil {
        panic(err)
    }
    // _dumpf("Database Document All >> %+v", data)
    _dumpf("Database Document All >> offset: %d", data["offset"])
    _dumpf("Database Document All >> total_rows: %d", data["total_rows"])
    // // _dumpf("Database Document All >> rows: %+v", data["rows"])
    _dumpf("Database Document All >> rows.0: %+v", u.Dig("rows.0", data))
    _dumpf("Database Document All >> rows.0.id: %s", u.Dig("rows.0.id", data))
    _dumpf("Database Document All >> rows.0.key: %s", u.Dig("rows.0.key", data))
    _dumpf("Database Document All >> rows.0.value.rev: %s", u.Dig("rows.0.value.rev", data))
    _dumpf("Database Document All >> rows.0.doc.name: %s", u.Dig("rows.0.doc.name", data))
}

/**
 * TestCreateDocument
 */
func TestCreateDocument() {
    data, err := database.CreateDocument(map[string]interface{}{
        "name": "kerem", "type": "tmp",
    })
    if err != nil {
        panic(err)
    }
    _dumpf("Database Create Document >> %+v", data)
    _dumpf("Database Create Document >> doc.ok: %v", data["ok"])
    _dumpf("Database Create Document >> doc.id: %s", data["id"])
    _dumpf("Database Create Document >> doc.rev: %s", data["rev"])
    // or
    for key, value := range data {
        _dumpf("Database Create Document >> doc.%s: %v", key, value)
    }
}

/**
 * TestCreateDocumentAll
 */
func TestCreateDocumentAll() {
    data, err := database.CreateDocumentAll([]interface{}{
        0: map[string]interface{}{"name": "kerem", "type": "tmp"},
        1: map[string]interface{}{"name": "murat", "type": "tmp"},
    })
    if err != nil {
        panic(err)
    }
    _dumpf("Database Create Document All >> %+v", data)
    _dumpf("Database Create Document All >> doc.0.ok: %v", u.Dig("0.ok", data))
    _dumpf("Database Create Document All >> doc.0.id: %s", u.Dig("0.id", data))
    _dumpf("Database Create Document All >> doc.0.rev: %s", u.Dig("0.rev", data))
    // or
    for i, doc := range data {
        for key, value := range doc {
            _dumpf("Database Create Document All >> doc.%d.%s: %v", i, key, value)
        }
    }
}

/**
 * TestUpdateDocument
 */
func TestUpdateDocument() {
    data, err := database.UpdateDocument(map[string]interface{}{
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
    _dumpf("Database Update Document >> %+v", data)
    _dumpf("Database Update Document >> doc.ok: %v", data["ok"])
    _dumpf("Database Update Document >> doc.id: %s", data["id"])
    _dumpf("Database Update Document >> doc.rev: %s", data["rev"])
    // or
    for key, value := range data {
        _dumpf("Database Update Document >> doc.%s: %v", key, value)
    }
}

/**
 * TestUpdateDocumentAll
 */
func TestUpdateDocumentAll() {
    data, err := database.UpdateDocumentAll([]interface{}{
        0: map[string]interface{}{"name": "kerem 2", "type": "tmp",
            "_id": "7ee9cdd673b109e030cec8c6f10020f7", "_rev": "1-3c92d3e67136c8b206d90ea37a3ee76d"},
        1: map[string]interface{}{"name": "murat 2", "type": "tmp",
            "_id": "7ee9cdd673b109e030cec8c6f1002cc1", "_rev": "1-09e886345b525e53892815baff169f03"},
    })
    if err != nil {
        panic(err)
    }
    _dumpf("Database Update Document All >> %+v", data)
    _dumpf("Database Update Document All >> doc.0.ok: %v", u.Dig("0.ok", data))
    _dumpf("Database Update Document All >> doc.0.id: %s", u.Dig("0.id", data))
    _dumpf("Database Update Document All >> doc.0.rev: %s", u.Dig("0.rev", data))
    // or
    for i, doc := range data {
        // check ok's
        if doc["ok"] == nil {
            _dumpf("Halt! error: doc.%d > %s reason: %s", i, doc["error"], doc["reason"])
        }
        for key, value := range doc {
            _dumpf("Database Update Document All >> doc.%d.%s: %v", i, key, value)
        }
    }
}

/**
 * TestDeleteDocument
 */
func TestDeleteDocument() {
    data, err := database.DeleteDocument(map[string]interface{}{
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
    _dumpf("Database Delete Document >> %+v", data)
}

/**
 * TestDeleteDocumentAll
 */
func TestDeleteDocumentAll() {
    data, err := database.DeleteDocumentAll([]interface{}{
        0: map[string]interface{}{"_id": "7ee9cdd673b109e030cec8c6f10020f7",
            "_rev": "3-1cb8867864cd6bb80361fed42c719897"},
        1: map[string]interface{}{"_id": "7ee9cdd673b109e030cec8c6f1002cc1",
            "_rev": "2-38a4374c91a4ae8d666d6c1fe13dd916"},
    })
    if err != nil {
        panic(err)
    }
    _dumpf("Database Delete Document All >> %+v", data)
    // or
    for i, doc := range data {
        // check ok's
        if doc["ok"] == nil {
            _dumpf("Halt! error: doc.%d > %s reason: %s", i, doc["error"], doc["reason"])
        }
        for key, value := range doc {
            _dumpf("Database Update Document All >> doc.%d.%s: %v", i, key, value)
        }
    }
}

/**
 * TestGetChanges
 */
func TestGetChanges() {
    data, err := database.GetChanges(nil, nil)
    if err != nil {
        panic(err)
    }
    // _dumpf("Database Changes >> %+v", data)
    _dumpf("Database Changes >> last_seq: %v", data["last_seq"])
    for i, result := range u.DigMapList("results", data) {
        _dumps(i,result)
    }
}

/**
 * TestCompact
 */
func TestCompact() {
    data, err := database.Compact("")
    if err != nil {
        panic(err)
    }
    _dumpf("Database Compact >> ok: %v", data["ok"])
}

/**
 * TestEnsureFullCommit
 */
func TestEnsureFullCommit() {
    data, err := database.EnsureFullCommit()
    if err != nil {
        panic(err)
    }
    _dumpf("Database Ensure Full Commit >> ok: %v", data["ok"])
    _dumpf("Database Ensure Full Commit >> instance_start_time: %s", data["instance_start_time"])
}
