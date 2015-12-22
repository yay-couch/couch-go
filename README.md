Just a wrapper for CouchDB like [Couch PHP](//github.com/qeremy/couch) or [Couch JS](//github.com/qeremy/couch-js) libraries.

Notice: See CouchDB's official documents before using this library.

## In a Nutshell
```go
// create a fresh document and save it
doc := Couch.NewDocument(database)
doc.Set("name", "The Doc!")
doc.Save()
// append an attachment and save it again
doc.SetAttachment(Couch.NewDocumentAttachment(doc, "./attc.txt"))
doc.Save()
```
