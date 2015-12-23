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
bool Server.Ping()
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
bool Server.Restart()
data, err := Server.GetConfig()
data, err := Server.GetConfigSection(section)
data, err := Server.GetConfigSectionKey(section, key)
data, err := Server.SetConfig(section, key, value)
data, err := Server.RemoveConfig(section, key)
```
