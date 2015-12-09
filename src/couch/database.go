package couch

import (
    "./util"
)

type Database struct {
    Client *Client
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

func NewDatabase(client *Client, name string) *Database {
    return &Database{
        Client: client,
          Name: name,
    }
}

func (this *Database) Ping() bool {
    return (200 == this.Client.Head(this.Name, nil, nil).GetStatusCode())
}

func (this *Database) Info() (map[string]interface{}, error) {
    data, err := this.Client.Get(this.Name, nil, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    for key, value := range data.(map[string]interface{}) {
        _return[key] = value
    }
    return _return, nil
}

func (this *Database) Create() bool {
    return (201 == this.Client.Put(this.Name, nil, nil, nil).GetStatusCode())
}

func (this *Database) Remove() bool {
    return (200 == this.Client.Delete(this.Name, nil, nil).GetStatusCode())
}

func (this *Database) Replicate(target string, targetCreate bool) (map[string]interface{}, error) {
    var body = util.ParamList(
        "source", this.Name,
        "target", target,
        "create_target", targetCreate,
    )
    data, err := this.Client.Post("/_replicate", nil, body, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    for key, value := range data.(map[string]interface{}) {
        if key == "history" {
            _return[key] = util.MapList(value)
            for i, history := range value.([]interface{}) {
                _return[key].([]map[string]interface{})[i] = util.Map()
                for kkey, vvalue := range history.(map[string]interface{}) {
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
    var query = util.ParamList(
        "include_docs", true,
        "key"         , util.Quote(key),
    )
    data, err := this.Client.Get(this.Name +"/_all_docs", query, nil).
        GetBodyData(&DatabaseDocumentList{})
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
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

func (this *Database) GetDocumentAll(query map[string]interface{}, keys []string) (
        map[string]interface{}, error) {
    query = util.Param(query)
    if query["include_docs"] == nil {
        query["include_docs"] = true
    }
    // reusable lambda
    var _return = func(data interface{}, err error) (map[string]interface{}, error) {
        if err != nil {
            return nil, err
        }
        var _return = util.Map()
        var _returnRows = data.(*DatabaseDocumentList).Rows
        _return["offset"]     = data.(*DatabaseDocumentList).Offset
        _return["total_rows"] = data.(*DatabaseDocumentList).TotalRows
        _return["rows"]       = util.MapList(len(_returnRows))
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
            this.Client.Post(this.Name +"/_all_docs", query, util.ParamList("keys", keys), nil).
                GetBodyData(&DatabaseDocumentList{}))
    }
}

func (this *Database) CreateDocument(document interface{}) (
        map[string]interface{}, error) {
    data, err := this.CreateDocumentAll([]interface{}{document})
    if err != nil {
        return nil, err
    }
    if data := data[0]; data != nil {
        return data, nil
    }
    return nil, nil
}

func (this *Database) CreateDocumentAll(documents []interface{}) (
        []map[string]interface{}, error) {
    var docs = util.MapList(documents)
    for i, doc := range documents {
        if docs[i] == nil {
            docs[i] = util.Map()
        }
        for key, value := range doc.(map[string]interface{}) {
            // this is create method, no update allowed
            if key == "_id" || key == "_id" || key == "_deleted" {
                continue
            }
            docs[i][key] = value
        }
    }
    var body = util.ParamList("docs", docs)
    data, err := this.Client.Post(this.Name +"/_bulk_docs", nil, body, nil).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.MapList(data)
    for i, doc := range data.([]interface{}) {
        if _return[i] == nil {
            _return[i] = util.Map()
        }
        for key, value := range doc.(map[string]interface{}) {
            _return[i][key] = value
        }
    }
    return _return, nil
}

func (this *Database) UpdateDocument(document interface{}) (
        map[string]interface{}, error) {
    data, err := this.UpdateDocumentAll([]interface{}{document})
    if err != nil {
        return nil, err
    }
    if data := data[0]; data != nil {
        return data, nil
    }
    return nil, nil
}

func (this *Database) UpdateDocumentAll(documents []interface{}) (
        []map[string]interface{}, error) {
    var docs = util.MapList(documents)
    for i, doc := range documents {
        if docs[i] == nil {
            docs[i] = util.Map()
        }
        for key, value := range doc.(map[string]interface{}) {
            docs[i][key] = value
        }
        // these are required params
        if docs[i]["_id"] == nil || docs[i]["_rev"] == nil {
            panic("Both _id & _rev fields are required!")
        }
    }
    var body = util.ParamList("docs", docs)
    data, err := this.Client.Post(this.Name +"/_bulk_docs", nil, body, nil).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.MapList(data)
    for i, doc := range data.([]interface{}) {
        if _return[i] == nil {
            _return[i] = util.Map()
        }
        for key, value := range doc.(map[string]interface{}) {
            _return[i][key] = value
        }
    }
    return _return, nil
}

func (this *Database) DeleteDocument(document interface{}) (
        map[string]interface{}, error) {
    data, err := this.DeleteDocumentAll([]interface{}{document})
    if err != nil {
        return nil, err
    }
    if data := data[0]; data != nil {
        return data, nil
    }
    return nil, nil
}

func (this *Database) DeleteDocumentAll(documents []interface{}) (
        []map[string]interface{}, error) {
    for i, _ := range documents {
        // just add "_deleted" param into document
        documents[i].(map[string]interface{})["_deleted"] = true
    }
    return this.UpdateDocumentAll(documents)
}

func (this *Database) GetChanges(query map[string]interface{}, docIds []string) (
        map[string]interface{}, error) {
    query = util.Param(query)
    if docIds != nil {
        query["filter"] = "_doc_ids"
    }
    var body = util.ParamList("doc_ids", docIds)
    data, err := this.Client.Post(this.Name +"/_changes", query, body, nil).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    _return["last_seq"] = util.Dig("last_seq", data)
    _return["results"]  = util.MapList(0) // set empty as default
    if results := data.(map[string]interface{})["results"].([]interface{});
       results != nil {
        _return["results"] = util.MapList(results) // @overwrite
        for i, result := range results {
            _return["results"].([]map[string]interface{})[i] = map[string]interface{}{
                     "id": util.Dig("id", result),
                    "seq": util.Dig("seq", result),
                "deleted": util.Dig("deleted", result),
                "changes": util.Dig("changes", result),
            }
        }
    }
    return _return, nil
}

func (this *Database) Compact(ddoc string) (map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_compact/"+ ddoc, nil, nil, nil).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "ok": util.DigBool("ok", data),
    }, nil
}

func (this *Database) EnsureFullCommit() (map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_ensure_full_commit", nil, nil, nil).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "ok": util.DigBool("ok", data),
        "instance_start_time": util.DigString("instance_start_time", data),
    }, nil
}

func (this *Database) ViewCleanup() (map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_view_cleanup", nil, nil, nil).
        GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "ok": util.DigBool("ok", data),
    }, nil
}

func (this *Database) ViewTemp(map_ string, reduce interface{}) (
        map[string]interface{}, error) {
    var body = util.ParamList(
        "map", map_,
        // prevent "missing function" error
        "reduce", util.IsEmptySet(reduce, nil),
    )
    data, err := this.Client.Post(this.Name +"/_temp_view", nil, body, nil).
        GetBodyData(&DatabaseDocumentList{})
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    var _returnRows = data.(*DatabaseDocumentList).Rows
    _return["offset"]     = data.(*DatabaseDocumentList).Offset
    _return["total_rows"] = data.(*DatabaseDocumentList).TotalRows
    _return["rows"]       = util.MapList(len(_returnRows))
    for i, row := range _returnRows {
        _return["rows"].([]map[string]interface{})[i] = map[string]interface{}{
               "id": row.Id,
              "key": row.Key,
            "value": row.Value,
        }
    }
    return _return, nil
}

func (this *Database) GetSecurity() (map[string]interface{}, error) {
    data, err := this.Client.Get(this.Name +"/_security", nil, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    return data.(map[string]interface{}), nil
}

func (this *Database) SetSecurity(admins, members map[string]interface{}) (
        map[string]interface{}, error) {
    if admins["names"].([]string) == nil || admins["roles"].([]string)  == nil ||
       members["names"].([]string) == nil || members["roles"].([]string) == nil {
        panic("Specify admins and/or members with names=>roles fields!")
    }
    var body = util.ParamList("admins", admins, "members", members)
    data, err := this.Client.Put(this.Name +"/_security", nil, body, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "ok": util.DigBool("ok", data),
    }, nil
}

func (this *Database) Purge(object map[string]interface{}) (map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_purge", nil, object, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    _return["purge_seq"] = util.DigInt("purge_seq", data)
    _return["purged"]  = util.Map()
    for id, revs := range data.(map[string]interface{})["purged"].
        (map[string]interface{}) {
        _return["purged"].(map[string]interface{})[id] = revs
    }
    return _return, nil
}

func (this *Database) GetMissingRevisions(object map[string]interface{}) (
        map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_missing_revs", nil, object, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    _return["missing_revs"] = util.Map()
    for id, revs := range data.(map[string]interface{})["missing_revs"].
        (map[string]interface{}) {
        _return["missing_revs"].(map[string]interface{})[id] = revs
    }
    return _return, nil
}

func (this *Database) GetMissingRevisionsDiff(object map[string]interface{}) (
        map[string]interface{}, error) {
    data, err := this.Client.Post(this.Name +"/_revs_diff", nil, object, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    for id, _ := range data.(map[string]interface{}) {
        _return[id] = map[string]interface{}{
            "missing": util.Dig(id +".missing", data),
        }
    }
    return _return, nil
}

func (this *Database) GetRevisionLimit() (int, error) {
    data, err := this.Client.Get(this.Name +"/_revs_limit", nil, nil).
        GetBodyData(nil)
    if err != nil {
        return -1, err
    }
    return int(data.(float64)), nil
}

func (this *Database) SetRevisionLimit(limit int) (map[string]interface{}, error) {
    data, err := this.Client.Put(this.Name +"/_revs_limit", nil, limit, nil).
        GetBodyData(limit)
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "ok": util.DigBool("ok", data),
    }, nil
}
