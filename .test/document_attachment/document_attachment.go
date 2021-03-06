package test

import (
   "couch"
   "couch/util"
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
   Document = couch.NewDocument(Database, "_id", DOCID);
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
   util.Dumpf("Document Attachment Ping >> %v", docAttc.Ping(200))
}

/**
 * TestFind
 */
func TestFind() {
   var docAttc = _documentAttachment("./attc.txt", "").Find()
   util.Dumpf("Document Attachment Find >> %v", docAttc)
   util.Dumpf("Document Attachment Find >> content: %s", docAttc["content"])
   util.Dumpf("Document Attachment Find >> content_type: %s", docAttc["content_type"])
   util.Dumpf("Document Attachment Find >> content_length: %d", docAttc["content_length"])
   util.Dumpf("Document Attachment Find >> digest: %s", docAttc["digest"])
}

/**
 * TestSave
 */
func TestSave() {
   var doc = couch.NewDocument(Database,
      "_id", "attc_test",
      "_rev", "15-521ec44a6b39f3c5b40b20f94408dd57",
   )
   var docAttc = couch.NewDocumentAttachment(doc, "./attc.txt", "attc2.txt")
   data, err := docAttc.Save()
   if err != nil {
      panic(err)
   }
   util.Dumpf("Document Attachment Save >> %v", data)
   util.Dumpf("Document Attachment Save >> ok: %v", data["ok"])
   util.Dumpf("Document Attachment Save >> id: %s", data["id"])
   util.Dumpf("Document Attachment Save >> rev: %s", data["rev"])
}

/**
 * TestRemove
 */
func TestRemove() {
   var doc = couch.NewDocument(Database,
      "_id", "attc_test",
      "_rev", "16-744735fac2c0bc49ff17f6a8fb8145bf",
   )
   var docAttc = couch.NewDocumentAttachment(doc, "./attc.txt", "attc2.txt")
   data, err := docAttc.Remove(false, false)
   if err != nil {
      panic(err)
   }
   util.Dumpf("Document Attachment Remove >> %v", data)
   util.Dumpf("Document Attachment Remove >> ok: %v", data["ok"])
   util.Dumpf("Document Attachment Remove >> id: %s", data["id"])
   util.Dumpf("Document Attachment Remove >> rev: %s", data["rev"])
}

/**
 * TestReadFile
 */
func TestReadFile() {
   _documentAttachment("./attc.txt", "").ReadFile(true)
}

/**
 * TestToArray
 */
func TestToArray() {
   var array = _documentAttachment("./attc.txt", "").ToArray(true)
   util.Dumpf("Document Attachment To Array >> %v", array)
}

/**
 * TestToJson
 */
func TestToJson() {
   var json = _documentAttachment("./attc.txt", "").ToJson(true)
   util.Dumpf("Document Attachment To JSON >> %v", json)
}
