package test

import (
    "./../../src/couch"
    "./../../src/couch/util"
)

var (
    DEBUG   = true
    DBNAME  = "foo2"
    DOCNAME = "attc_test"
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
    Document = couch.NewDocument(Database, util.ParamList("id", DOCNAME));
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
    var docAttc = _documentAttachment("attc.txt", "")
    util.Dumpf("Document Ping >> %#v", docAttc)
    // util.Dumpf("Document Ping >> %v", docAttc.Ping(304))
}
