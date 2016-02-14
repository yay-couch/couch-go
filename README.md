Simply port of [Couch](https://github.com/yay-couch/couch) library for Go.

Notice: See CouchDB's official documents before using this library.

## In a Nutshell

```go
// create a fresh document
doc := Couch.NewDocument(db)
doc.Set("name", "The Doc!")
doc.Save()

// append an attachment to the same document above
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
uri       := "/<URI>"
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
ok            := Database.Ping()
data, err     := Database.Info()
ok            := Database.Create()
ok            := Database.Remove()
data, err     := Database.Replicate(target, targetCreate)
data, err     := Database.GetDocument(key)
data, err     := Database.GetDocumentAll(query)
data, err     := Database.CreateDocument(document)
data, err     := Database.CreateDocumentAll([]document)
data, err     := Database.UpdateDocument(document)
data, err     := Database.UpdateDocumentAll([]document)
data, err     := Database.DeleteDocument(document)
data, err     := Database.DeleteDocumentAll([]document)
data, err     := Database.GetChanges(query, docIds)
data, err     := Database.Compact(ddoc)
ok, time, err := Database.EnsureFullCommit()
ok,   err     := Database.ViewCleanup()
data, err     := Database.ViewTemp(map, reduce)
data, err     := Database.GetSecurity()
data, err     := Database.SetSecurity(admins, members)
data, err     := Database.Purge(object)
data, err     := Database.GetMissingRevisions(object)
data, err     := Database.GetMissingRevisionsDiff(object)
limit,err     := Database.GetRevisionLimit()
ok,   err     := Database.SetRevisionLimit(limit)

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

// examples
var doc = couch.NewDocument(Database)
doc.Set("_id", "1ec1098a")
// or like key=>value pairs
doc.Set(
    "_id", "1ec1098a",
    "_rev", "1-5637fd00",
)

var doc = couch.NewDocument(
    Database,
    "_id", "1ec1098a",
    // ...
)
data, err := doc.Find(nil)
if err != nil {
    panic(err)
}
util.Dumpf("Document Find >> %v", data)
util.Dumpf("Document Find >> _id: %s", data["_id"])

var doc = couch.NewDocument(
    Database,
    "_id", "1ec1098a",
)
type Doc struct {
    Id   string `json:"_id"`
    Rev  string `json:"_rev"`
    Name string
    // ...
}
data, err := doc.FindStruct(&Doc{}, nil)
if err != nil {
    panic(err)
}
util.Dumpf("Document Find Func >> doc: %+v", data)
util.Dumpf("Document Find Func >> doc._id: %s", data.(*Doc).Id)
util.Dumpf("Document Find Func >> doc._rev: %s", data.(*Doc).Rev)
util.Dumpf("Document Find Func >> doc.name: %s", data.(*Doc).Name)
```

### DocumentAttachment Object

```go
DocumentAttachment := couch.NewDocumentAttachment(Document, "./attc.txt", "")

// methods
DocumentAttachment.SetDocument(document *Document)
DocumentAttachment.GetDocument() *Document

attcArray := DocumentAttachment.ToArray(?encode)
attcJson  := DocumentAttachment.ToJson(?encode)
ok        := DocumentAttachment.Ping(statusCodes...[200,304])
attc      := DocumentAttachment.Find()
data, err := DocumentAttachment.Save()
data, err := DocumentAttachment.Remove(args...[?batch, ?fullCommit])
DocumentAttachment.ReadFile(?encode)

// examples
attc := couch.NewDocumentAttachment(do, "./attc1.txt", "")
attc.Find()

// find an attachment by digest
attc.Digest   = "U1p5BLvdnOZVRyR6YrXBoQ=="
attc.Find()

// add an attachment to document
attc.File     = "attc2.txt"
attc.FileName = "attc2"
attc.Save()

// remove an attachment from document
attc.FileName = "attc2"
attc.Remove()
```

### DocumentDesign Object

```go
// @todo
```

## Uuid

```go
import "couch/uuid"

// create uuid
uuid := uuid.New(true)      // auto-generate randomly using "crypto/rand"
uuid := uuid.New("<DOCID>") // set given value

uuid := uuid.New(nil)       // set later
uuid.SetValue("...")

// methods
uuid.SetValue(value)
value := uuid.GetValue() interface{}
value := uuid.ToString() string
uuid  := uuid.Generate(limit /* type */) string

// generation types (limits)
RFC
HEX_8
HEX_32
HEX_40
TIMESTAMP
TIMESTAMP_NANO

// examples
dump uuid.Generate(uuid.RFC)       // rfc uuid/v4
dump uuid.Generate(uuid.HEX_8)     // hexed 8 bytes
dump uuid.Generate(uuid.TIMESTAMP) // unix epoch
```

## Query

```go
import "couch/query"

query := query.New(data) // data = map | nil

// methods
query.Set(key, value) *Query
query.Skip(value) *Query
query.Limit(value) *Query

value := query.Get(key)
data  := query.ToData()
query := query.ToString()

// examples
query.Set("conflicts", true).
      Set("stale", "ok").
      Skip(1).
      Limit(2)

dump query.ToString() // conflicts=true&stale=ok&skip=1&limit=2
```

## Request / Response

```go
// after any http stream (server ping, database ping, document save etc)

// ie.
client.DoRequest("GET /", nil, nil, nil)

// get raw stuffs
dump client.GetRequest().ToString()
dump client.GetResponse().ToString()

/*
GET / HTTP/1.0
Host: localhost:5984
Connection: close
Accept: application/json
...

HTTP/1.0 200 OK
Server: CouchDB/1.5.0 (Erlang OTP/R16B03)
Date: Sun, 01 Nov 2015 18:04:42 GMT
Content-Type: application/json
...

{"couchdb":"Welcome","uuid":"5a660f4695a5fa9ab2cd22722bc01e96", ...
*/

// get response header
date := client.GetResponse().GetHeader("Date")

// get response all headers
headers := client.GetResponse().GetHeaderAll()
for key, value := range headers {
    dump key, value
}

// get response body
body := client.GetResponse().GetBody()

// get response body to
type Body struct {
    CouchDB string
    ...
}
body := client.GetResponse().GetBodyData(&Body{})
dump body.(*Body).CouchDB
```

## Error Handling

Couch will not throw any server response error, such as 409 Conflict etc. It only throws library-related errors ie. wrong usages of the library (ie. when _id is required for some action but you did not provide it).

```go
// create issue
doc := Couch.NewDocument(db)
doc.Set("_id", "an_existing_docid")

// no error will be displayed
doc.Save()

// but could be so
if 201 != client.GetResponse().GetStatusCode() {
    dump "nö!"
    // or print response error
    dump client.GetResponse().GetError()
}
```

Note: See the test folders for more.

