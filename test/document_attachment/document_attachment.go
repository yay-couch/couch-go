package test

import (
    "./../../src/couch"
    "./../../src/couch/util"
)

var (
    DEBUG   = true
    DBNAME  = "foo2"
    DOCID   = "attc_test"
)

var (
    Couch    *couch.Couch
    Client   *couch.Client
    Database *couch.Database
    Document *couch.Document
)

func init() {
    Couch    = couch.New(nil, DEBUG)
    Client   = couch.NewClient(Couch)
    Database = couch.NewDatabase(Client, DBNAME);
    Document = couch.NewDocument(Database, util.ParamList("_id", DOCID));
}

func _documentAttachment(file, fileName string) *couch.DocumentAttachment {
    return couch.NewDocumentAttachment(Document, file, fileName)
}

/**
 * TestAll
 */
func TestAll() {}

/**
 * TestPing
 */
func TestPing() {
    var docAttc = _documentAttachment("./attc.txt", "")
    util.Dumpf("Document Ping >> %v", docAttc.Ping(200))
}

/**
 * TestFind
 */
func TestFind() {
    var docAttc = _documentAttachment("./attc.txt", "").Find()
    util.Dumpf("Document Find >> %v", docAttc)
    util.Dumpf("Document Find >> content: %s", docAttc["content"])
    util.Dumpf("Document Find >> content_type: %s", docAttc["content_type"])
    util.Dumpf("Document Find >> content_length: %d", docAttc["content_length"])
    util.Dumpf("Document Find >> digest: %s", docAttc["digest"])
}
