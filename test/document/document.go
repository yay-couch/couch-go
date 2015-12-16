package server

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

func _document(data map[string]interface{}) *couch.Document {
    return couch.NewDocument(Database, data)
}

/**
 * TestAll
 */
func TestAll() {}

/**
 * TestPing
 */
func TestPing() {
    var doc1 = _document(map[string]interface{}{
        "_id": "0f1eb3ba90772b64aee2f44b3c00055b",
    })
    var doc2 = _document(map[string]interface{}{
        "_id": "0f1eb3ba90772b64aee2f44b3c00055b",
        "_rev": "1-3c92d3e67136c8b206d90ea37a3ee76d",
    })
    util.Dumpf("Document Ping >> %v", doc1.Ping(200))
    util.Dumps("\n---\n")
    util.Dumpf("Document Ping >> %v", doc2.Ping(304))
}

/**
 * TestIsExists
 */
func TestIsExists() {
    var doc = _document(map[string]interface{}{
        "_id": "0f1eb3ba90772b64aee2f44b3c00055b",
    })
    util.Dumpf("Document Is Exists >> %v", doc.IsExists())
}

/**
 * TestIsNotModified
 */
func TestIsNotModified() {
    var doc = _document(map[string]interface{}{
        "_id": "0f1eb3ba90772b64aee2f44b3c00055b",
        "_rev": "1-3c92d3e67136c8b206d90ea37a3ee76d",
    })
    util.Dumpf("Document Is Not Modified >> %v", doc.IsNotModified())
}

/**
 * TestFind
 */
func TestFind() {
    var doc = _document(map[string]interface{}{
        "_id": "0f1eb3ba90772b64aee2f44b3c00055b",
    })
    data, err := doc.Find(nil)
    if err != nil {
        panic(err)
    }
    util.Dumpf("Document Find >> %v", data)
    util.Dumpf("Document Find >> _id: %s", data["_id"])
}
