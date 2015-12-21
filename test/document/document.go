package test

import (
    "./../../src/couch"
    "./../../src/couch/util"
)

var (
    DEBUG  = true
    DBNAME = "foo2"
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

func _document(data ...interface{}) *couch.Document {
    return couch.NewDocument(Database, data...)
}

/**
 * TestAll
 */
func TestAll() {}

/**
 * TestPing
 */
func TestPing() {
    var doc1 = _document(
        "_id", "0f1eb3ba90772b64aee2f44b3c00055b",
    )
    var doc2 = _document(
        "_id", "0f1eb3ba90772b64aee2f44b3c00055b",
        "_rev", "1-3c92d3e67136c8b206d90ea37a3ee76d",
    )
    util.Dumpf("Document Ping >> %v", doc1.Ping(200))
    util.Dumps("\n---\n")
    util.Dumpf("Document Ping >> %v", doc2.Ping(304))
}

/**
 * TestIsExists
 */
func TestIsExists() {
    var doc = _document(
        "_id", "0f1eb3ba90772b64aee2f44b3c00055b",
    )
    util.Dumpf("Document Is Exists >> %v", doc.IsExists())
}

/**
 * TestIsNotModified
 */
func TestIsNotModified() {
    var doc = _document(
        "_id", "0f1eb3ba90772b64aee2f44b3c00055b",
        "_rev", "1-3c92d3e67136c8b206d90ea37a3ee76d",
    )
    util.Dumpf("Document Is Not Modified >> %v", doc.IsNotModified())
}

/**
 * TestFind
 */
func TestFind() {
    var doc = _document(
        "_id", "0f1eb3ba90772b64aee2f44b3c00055b",
    )
    data, err := doc.Find(nil)
    if err != nil {
        panic(err)
    }
    util.Dumpf("Document Find >> %v", data)
    util.Dumpf("Document Find >> _id: %s", data["_id"])
}

/**
 * TestFindRevisions
 */
func TestFindRevisions() {
    var doc = _document(
        "_id", "83b5e0a0b3bd41d9a21cee7ae8000615",
    )
    data, err := doc.FindRevisions()
    if err != nil {
        panic(err)
    }
    util.Dumpf("Document Find Revisions >> %v", data)
    util.Dumpf("Document Find Revisions >> start: %d", data["start"])
    util.Dumpf("Document Find Revisions >> ids: %v", data["ids"])
    util.Dumpf("Document Find Revisions >> ids.0: %v", util.Dig("ids.0", data))
}

/**
 * TestFindRevisionsExtended
 */
func TestFindRevisionsExtended() {
    var doc = _document(
        "_id", "83b5e0a0b3bd41d9a21cee7ae8000615",
    )
    data, err := doc.FindRevisionsExtended()
    if err != nil {
        panic(err)
    }
    util.Dumpf("Document Find Revisions Extended >> %v", data)
    util.Dumpf("Document Find Revisions Extended >> 0: %v", data[0])
    util.Dumpf("Document Find Revisions Extended >> 0.rev: %s", data[0]["rev"])
    util.Dumpf("Document Find Revisions Extended >> 0.status: %s", data[0]["status"])
}

/**
 * TestFindAttachments
 */
func TestFindAttachments() {
    data, err := _document("_id", "attc_test").FindAttachments(false, nil)
    if err != nil {
        panic(err)
    }
    util.Dumpf("Document Find Attachments >> %v", data)
    util.Dumpf("Document Find Attachments >> 0: %v", data[0])
    util.Dumpf("Document Find Attachments >> 0.content_type: %v", data[0]["content_type"])
}

/**
 * TestSave
 */
func TestSave() {
    var doc = _document(
        "_id", "go_test_attc",
        "val1", "is val 1",
        "val2", "is val 2",
        "_attachments", []interface{}{
            couch.NewDocumentAttachment(nil, "./attc.txt", "attc1"),
            couch.NewDocumentAttachment(nil, "./attc.txt", "attc2"),
            map[string]interface{}{"file": "./attc.txt", "fileName": "attc3"},
        },
    )
    data, err := doc.Save(false, false)
    if err != nil {
        panic(err)
    }
    util.Dumpf("Document Save >> %v", data)
    util.Dumpf("Document Save >> ok: %v", data["ok"])
    util.Dumpf("Document Save >> id: %v", data["id"])
    util.Dumpf("Document Save >> rev: %v", data["rev"])
}

/**
 * TestRemove
 */
func TestRemove() {
    var doc = _document(
        "_id", "e90636c398458a9d5969d2e71b04bc2a",
        "_rev", "1-5637fdf6ae62130da1dda54be05d7da7",
    )
    data, err := doc.Remove()
    if err != nil {
        panic(err)
    }
    util.Dumpf("Document Remove >> %v", data)
    util.Dumpf("Document Remove >> ok: %v", data["ok"])
    util.Dumpf("Document Remove >> id: %v", data["id"])
    util.Dumpf("Document Remove >> rev: %v", data["rev"])
}
