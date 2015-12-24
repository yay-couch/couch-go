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

If you need any direct request for any reason, you can use the methods below.

```go
res := Client.DoRequest("GET /<URI>", uriParams<map>, body<any>, headers<map>)

// data type is not specified
data, err := res.GetBodyData(nil)

// data type is MyDoc
type MyDoc struct {
    Id  string
    Rev string
    // ...
}
data, err := res.GetBodyData(&MyDoc{})

// args
uri       := "/<URI>";
uriParams := util.ParamList("param_name", "param_value", ...)
headers   := util.ParamList("X-Foo", "The foo!", ...)
body      := ""

// shortcut methods that handle HEAD, GET, POST, PUT, COPY, DELETE
Client.Head(uri, uriParams, headers)
Client.Get(uri, uriParams, headers)
Client.Copy(uri, uriParams, headers)
Client.Delete(uri, uriParams, headers)

// with body
Client.Put(uri, uriParams, body, headers)
Client.Post(uri, uriParams, body, headers)

// after request operations
Request  := Client.GetRequest()  // *couch.http.Request
Response := Client.GetResponse() // *couch.http.Response
```

### Server Object
```go
Server := couch.NewServer(Client)

// methods
ok        := Server.Ping()
data, err := Server.Info()
data, err := Server.Version()
data, err := Server.GetActiveTasks()
data, err := Server.GetAllDatabases()
data, err := Server.GetDatabaseUpdates(query)
data, err := Server.GetLogs(query)
data, err := Server.GetStats(path)
data, err := Server.GetUuid()
data, err := Server.GetUuids(count)
data, err := Server.Replicate(body)
ok        := Server.Restart()
data, err := Server.GetConfig()
data, err := Server.GetConfigSection(section)
data, err := Server.GetConfigSectionKey(section, key)
data, err := Server.SetConfig(section, key, value)
data, err := Server.RemoveConfig(section, key)
```

### Database Object
```go
Database := couch.NewDatabase(Client, "foo")

// methods
ok        := Database.Ping()
data, err := Database.Info()
ok        := Database.Create()
ok        := Database.Remove()
data, err := Database.Replicate(target, targetCreate)
data, err := Database.GetDocument(key)
data, err := Database.GetDocumentAll(query)
data, err := Database.CreateDocument(document)
data, err := Database.CreateDocumentAll([]document)
data, err := Database.UpdateDocument(document)
data, err := Database.UpdateDocumentAll([]document)
data, err := Database.DeleteDocument(document)
data, err := Database.DeleteDocumentAll([]document)
data, err := Database.GetChanges(query, docIds)
data, err := Database.Compact(ddoc)
data, err := Database.EnsureFullCommit()
ok,   err := Database.ViewCleanup()
data, err := Database.ViewTemp(map, reduce)
data, err := Database.GetSecurity()
data, err := Database.SetSecurity(admins, members)
data, err := Database.Purge(object)
data, err := Database.GetMissingRevisions(object)
data, err := Database.GetMissingRevisionsDiff(object)
limit,err := Database.GetRevisionLimit()
ok,   err := Database.SetRevisionLimit(limit)

// examples
data, err := Database.CreateDocument(map[string]interface{}{
    "name": "CouchDB", "is_nosql": true,
})
if err != nil {
    panic(err)
}
util.Dumpf("Create Document >> %+v", data)
util.Dumpf("Create Document >> doc.ok: %v", data["ok"])
util.Dumpf("Create Document >> doc.id: %s", data["id"])
util.Dumpf("Create Document >> doc.rev: %s", data["rev"])
// or
for key, value := range data {
    util.Dumpf("Create Document >> doc.%s: %v", key, value)
}

data, err := Database.CreateDocumentAll([]interface{}{
    0: map[string]interface{}{"name": "CouchDB", "is_nosql": true},
    1: map[string]interface{}{"name": "MongoDB", "is_nosql": true},
    2: map[string]interface{}{"name": "MySQL", "is_nosql": false},
})
if err != nil {
    panic(err)
}
util.Dumpf("Create Document All >> %+v", data)
util.Dumpf("Create Document All >> doc.0.ok: %v", util.Dig("0.ok", data))
util.Dumpf("Create Document All >> doc.0.id: %s", util.Dig("0.id", data))
util.Dumpf("Create Document All >> doc.0.rev: %s", util.Dig("0.rev", data))
// or
for i, doc := range data {
    for key, value := range doc {
        util.Dumpf("Create Document All >> doc.%d.%s: %v", i, key, value)
    }
}
```

### Document Object
```go
Document := couch.NewDocument(Database, data...)

// methods
Document.SetDatabase(*Database)
Document.GetDatabase() *Database
Document.Set(data...) *Document
Document.SetId(id)
Document.SetRev(rev)
Document.SetDeleted(deleted)
Document.SetAttachment(attachment)
Document.SetData(data)
Document.Get(key)
Document.GetId()
Document.GetRev()
Document.GetDeleted()
Document.GetAttachment(fileName)
Document.GetData()

ok        := Document.Ping()
ok        := Document.IsExists()
ok        := Document.IsNotModified()
data, err := Document.Find(query)
data, err := Document.FindStruct(struct, query)
data, err := Document.FindRevisions()
data, err := Document.FindRevisionsExtended()
data, err := Document.FindAttachments(?attEncInfo, []attsSince)
data, err := Document.Save(args...[?batch, ?fullCommit])
data, err := Document.Remove(args...[?batch, ?fullCommit])
data, err := Document.Copy(dest, args...[?batch, ?fullCommit])
data, err := Document.CopyFrom(dest, args...[?batch, ?fullCommit])
data, err := Document.CopyTo(dest, destRev, args...[?batch, ?fullCommit])
```
