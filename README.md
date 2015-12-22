Just a wrapper for CouchDB like [Couch PHP](//github.com/qeremy/couch) or [Couch JS](//github.com/qeremy/couch-js) libraries.

Notice: See CouchDB's official documents before using this library.

## In a Nutshell
```go
// create a fresh document
doc := Couch.NewDocument(db)
doc.Set("name", "The Doc!")
doc.Save()

// append an attachment to same document
doc.SetAttachment(Couch.NewDocumentAttachment(doc, "./file.txt"))
doc.Save()
```

## Configuration
Configuration is optional but you can provide all these options as `map`;
```go
map[string]interface{}{
    Scheme        : "http",
    Host          : "localhost",
    Port   uint16 : 5984,
    Username      : "",
    Password      : "",
}
```
## Objects

### Couch Object
```go
// init couch object with default config and without debug
Couch := couch.New(nil, false)

// init couch object with given config and debug
Couch := couch.New(config, true)

// or set later but before streaming
Couch := couch.New(nil, true)
Couch.SetConfig(config)
```

### Client Object
```go
// used in Server and Database objects
Client := couch.NewClient(Couch)
```
