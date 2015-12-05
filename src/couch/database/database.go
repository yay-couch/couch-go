package database

import _client "./../client"

import u "./../util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Database struct {
    Client *_client.Client
    Name   string
}
type DatabaseDocument struct {
    Id     string
    Key    string
    Value  map[string]interface{}
    Doc    map[string]interface{}
}
type DatabaseDocumentList struct {
    Offset     uint
    TotalRows  uint `json:"total_rows"`
    UpdateSeq  uint `json:"update_seq"`
    Rows       []DatabaseDocument
}

func Shutup() {}

func New(client *_client.Client, name string) *Database {
    return &Database{
        Client: client,
          Name: name,
    }
}

func (this *Database) Ping() bool {
    return (200 == this.Client.Head(this.Name, nil, nil).GetStatusCode())
}

func (this *Database) Info() (map[string]interface{}, error) {
    data, err := this.Client.Get(this.Name, nil, nil).
        GetBodyData(map[string]interface{}{})
    if err != nil {
        return nil, err
    }
    var _return = u.Map()
    for key, value := range data.(map[string]interface{}) {
        _return[key] = value
    }
    return _return, nil
}

func (this *Database) Create() bool {
    return (201 == this.Client.Put(this.Name, nil, nil, nil).
        GetStatusCode())
}

func (this *Database) Remove() bool {
    return (200 == this.Client.Delete(this.Name, nil, nil).
        GetStatusCode())
}

func (this *Database) Replicate(target string, targetCreate bool) (map[string]interface{}, error) {
    var body = u.ParamList(
        "source", this.Name,
        "target", target,
        "create_target", targetCreate,
    )
    data, err := this.Client.Post("/_replicate", nil, body, nil).
        GetBodyData(map[string]interface{}{})
    if err != nil {
        return nil, err
    }
    var _return = u.Map()
    for key, value := range data.(map[string]interface{}) {
        if key == "history" {
            _return[key] = u.MapList(value)
            for i, history := range value.([]interface{}) {
                for kkey, vvalue := range history.(map[string]interface{}) {
                    if _return[key].([]map[string]interface{})[i] == nil {
                        _return[key].([]map[string]interface{})[i] = u.Map()
                    }
                    _return[key].([]map[string]interface{})[i][kkey] = vvalue
                }
            }
            continue
        }
        _return[key] = value
    }
    return _return, nil
}

func (this *Database) GetDocument(key string) (map[string]interface{}, error) {
    var query = u.ParamList(
        "include_docs", true,
        "key"         , u.Quote(key),
    )
    data, err := this.Client.Get(this.Name +"/_all_docs", query, nil).
        GetBodyData(&DatabaseDocumentList{})
    if err != nil {
        return nil, err
    }
    var _return = u.Map()
    for _, doc := range data.(*DatabaseDocumentList).Rows {
        _return["id"]    = doc.Id
        _return["key"]   = doc.Key
        _return["value"] = map[string]string{"rev": doc.Value["rev"].(string)}
        _return["doc"]   = map[string]interface{}{}
        for key, value := range doc.Doc {
            _return["doc"].(map[string]interface{})[key] = value
        }
    }
    return _return, nil
}

func (this *Database) GetDocumentAll(query map[string]interface{}, keys []string) (map[string]interface{}, error) {
    query = u.Param(query)
    if query["include_docs"] == nil {
        query["include_docs"] = true
    }
    // reusable lambda
    var _return = func(data interface{}, err error) (map[string]interface{}, error) {
        if err != nil {
            return nil, err
        }
        var _return = u.Map()
        var _returnRows = data.(*DatabaseDocumentList).Rows
        _return["offset"]     = data.(*DatabaseDocumentList).Offset
        _return["total_rows"] = data.(*DatabaseDocumentList).TotalRows
        _return["rows"]       = u.MapList(len(_returnRows))
        for i, row := range _returnRows {
            _return["rows"].([]map[string]interface{})[i] = map[string]interface{}{
                   "id": row.Id,
                  "key": row.Key,
                "value": map[string]string{"rev": row.Value["rev"].(string)},
                  "doc": row.Doc,
            }
        }
        return _return, nil
    }
    if keys == nil {
        return _return(
            this.Client.Get(this.Name +"/_all_docs", query, nil).
                GetBodyData(&DatabaseDocumentList{}))
    } else {
        return _return(
            this.Client.Post(this.Name +"/_all_docs", query, u.ParamList("keys", keys), nil).
                GetBodyData(&DatabaseDocumentList{}))
    }
}

func (this *Database) CreateDocument(document interface{}) (map[string]interface{}, error) {
    data, err := this.CreateDocumentAll([]interface{}{document})
    if err != nil {
        return nil, err
    }
    if data := data[0]; data != nil {
        return data, nil
    }
    return nil, nil
}

func (this *Database) CreateDocumentAll(documents []interface{}) ([]map[string]interface{}, error) {
    var docs = u.MapList(documents)
    for i, doc := range documents {
        if docs[i] == nil {
            docs[i] = u.Map()
        }
        for key, value := range doc.(map[string]interface{}) {
            // this is create method, no update allowed
            if key == "_id" || key == "_id" || key == "_deleted" {
                continue
            }
            docs[i][key] = value
        }
    }
    data, err := this.Client.Post(this.Name +"/_bulk_docs", nil, u.ParamList("docs", docs), nil).
        GetBodyData([]interface{}{})
    if err != nil {
        return nil, err
    }
    var _return = u.MapList(data)
    for i, doc := range data.([]interface{}) {
        if _return[i] == nil {
            _return[i] = u.Map()
        }
        for key, value := range doc.(map[string]interface{}) {
            _return[i][key] = value
        }
    }
    return _return, nil
}

func (this *Database) UpdateDocument(document interface{}) (map[string]interface{}, error) {
    data, err := this.UpdateDocumentAll([]interface{}{document})
    if err != nil {
        return nil, err
    }
    if data := data[0]; data != nil {
        return data, nil
    }
    return nil, nil
}

func (this *Database) UpdateDocumentAll(documents []interface{}) ([]map[string]interface{}, error) {
    var docs = u.MapList(documents)
    for i, doc := range documents {
        if docs[i] == nil {
            docs[i] = u.Map()
        }
        for key, value := range doc.(map[string]interface{}) {
            docs[i][key] = value
        }
        // these are required params
        if docs[i]["_id"] == nil || docs[i]["_rev"] == nil {
            panic("Both _id & _rev fields are required!")
        }
    }
    data, err := this.Client.Post(this.Name +"/_bulk_docs", nil, u.ParamList("docs", docs), nil).
        GetBodyData([]interface{}{})
    if err != nil {
        return nil, err
    }
    var _return = u.MapList(data)
    for i, doc := range data.([]interface{}) {
        if _return[i] == nil {
            _return[i] = u.Map()
        }
        for key, value := range doc.(map[string]interface{}) {
            _return[i][key] = value
        }
    }
    return _return, nil
}

func (this *Database) DeleteDocument(document interface{}) (map[string]interface{}, error) {
    data, err := this.DeleteDocumentAll([]interface{}{document})
    if err != nil {
        return nil, err
    }
    if data := data[0]; data != nil {
        return data, nil
    }
    return nil, nil
}

func (this *Database) DeleteDocumentAll(documents []interface{}) ([]map[string]interface{}, error) {
    for i, _ := range documents {
        // just add "_deleted" param into document
        documents[i].(map[string]interface{})["_deleted"] = true
    }
    return this.UpdateDocumentAll(documents)
}

func (this *Database) GetChanges(query map[string]interface{}, docIds []string) (map[string]interface{}, error) {
    query = u.Param(query)
    if docIds != nil {
        query["filter"] = "_doc_ids"
    }
    data, err := this.Client.Post(this.Name +"/_changes", query, map[string]interface{}{
        "doc_ids": docIds,
    }, nil).GetBodyData(map[string]interface{}{})
    if err != nil {
        return nil, err
    }
    var _return = u.Map()
    _return["last_seq"] = u.Dig("last_seq", data)
    _return["results"]  = u.MapList(0)
    if results := data.(map[string]interface{})["results"].([]interface{}); results != nil {
        _return["results"] = u.MapList(len(results))
        for i, result := range results {
            _return["results"].([]map[string]interface{})[i] = map[string]interface{}{
                     "id": u.Dig("id", result),
                    "seq": u.Dig("seq", result),
                "deleted": u.Dig("deleted", result),
                "changes": u.Dig("changes", result),
            }
        }
    }
    return _return, nil
}

func (this *Database) Compact(ddoc string) (map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_compact/"+ ddoc, nil, nil, nil).
        GetBodyData(map[string]interface{}{})
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "ok": u.DigBool("ok", data),
    }, nil
}

func (this *Database) EnsureFullCommit() (map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_ensure_full_commit", nil, nil, nil).
        GetBodyData(map[string]interface{}{})
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "ok": u.DigBool("ok", data),
        "instance_start_time": u.DigString("instance_start_time", data),
    }, nil
}

func (this *Database) ViewCleanup() (map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_view_cleanup", nil, nil, nil).
        GetBodyData(map[string]interface{}{})
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "ok": u.DigBool("ok", data),
    }, nil
}

func (this *Database) ViewTemp(map_, reduce string) (map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_temp_view", nil, map[string]interface{}{
        "map": map_,
        "reduce": u.IsEmptySet(reduce, nil), // prevent "missing function" error
    }, nil).GetBodyData(&DatabaseDocumentList{})
    if err != nil {
        return nil, err
    }
    var _return = u.Map()
    _return["offset"]     = data.(*DatabaseDocumentList).Offset
    _return["total_rows"] = data.(*DatabaseDocumentList).TotalRows
    _return["rows"]       = u.MapList(len(data.(*DatabaseDocumentList).Rows))
    for i, row := range data.(*DatabaseDocumentList).Rows {
        _return["rows"].([]map[string]interface{})[i] = map[string]interface{}{
               "id": row.Id,
              "key": row.Key,
            "value": row.Value,
        }
    }
    return _return, nil
}
