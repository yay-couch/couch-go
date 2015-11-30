package server

/**
 * Links
 * - http://blog.golang.org/json-and-go
 * - http://golang.org/pkg/encoding/json/#example_Unmarshal
 */

import _client   "./../client"
// import _response "./../http/response"

import u "./../util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Server struct {
    Client *_client.Client
}

func Shutup() {}

func New(client *_client.Client) *Server {
    return &Server{
        Client: client,
    }
}

func (this *Server) Ping() bool {
    return (200 == this.Client.Head("/", nil, nil).GetStatusCode())
}

func (this *Server) Info() (map[string]interface{}, error) {
    type Data struct {
        CouchDB string
        Uuid    string
        Version string
        Vendor  map[string]string
    }
    data, err := this.Client.Get("/", nil, nil).GetData(&Data{})
    if err != nil {
        return nil, err
    }
    var _return = make(map[string]interface{});
    _return["couchdb"] = data.(*Data).CouchDB
    _return["uuid"]    = data.(*Data).Uuid
    _return["version"] = data.(*Data).Version
    _return["vendor"]  = map[string]string{
           "name": data.(*Data).Vendor["name"],
        "version": data.(*Data).Vendor["version"],
    }
    return _return, nil
}

func (this *Server) Version() (string, error) {
    data, err := this.Info()
    if err != nil {
        return "", err
    }
    return u.Dig("version", data).(string), nil
}

func (this *Server) GetActiveTasks() ([]map[string]interface{}, error) {
    type Data []map[string]interface{}
    data, err := this.Client.Get("/_active_tasks", nil, nil).GetData(&Data{})
    if err != nil {
        return nil, err
    }
    var _return = make([]map[string]interface{}, len(*data.(*Data)));
    for i, data := range *data.(*Data) {
        _return[i] = make(map[string]interface{})
        for key, value := range data {
            _return[i][key] = value
        }
    }
    return _return, nil
}

func (this *Server) GetAllDatabases() ([]string, error) {
    type Data []string
    data, err := this.Client.Get("/_all_dbs", nil, nil).GetData(&Data{})
    if err != nil {
        return nil, err
    }
    var _return = make([]string, len(*data.(*Data)))
    for i, db := range *data.(*Data) {
        _return[i] = db
    }
    return _return, nil
}

func (this *Server) GetDatabaseUpdates(query interface{}) (map[string]interface{}, error) {
    type Data interface{}
    data, err := this.Client.Get("/_db_updates", query, nil).GetData(*new(Data))
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "db_name": data.(map[string]interface{})["db_name"],
           "type": data.(map[string]interface{})["type"],
             "ok": data.(map[string]interface{})["o"],
    }, nil
}

func (this *Server) GetLogs(query interface{}) string {
    return this.Client.Get("/_log", query, map[string]interface{}{
        "Accept": "text/plain",
    }).GetBody()
}

func (this *Server) GetStats(path string) (map[string]map[string]map[string]interface{}, error) {
    type Data map[string]map[string]map[string]interface{}
    data, err := this.Client.Get("/_stats/"+ path, nil, nil).GetData(&Data{})
    if err != nil {
        return nil, err
    }
    var _return = make(map[string]map[string]map[string]interface{})
    for i, data := range *data.(*Data) {
        _return[i] = make(map[string]map[string]interface{})
        for ii, ddata := range data {
            _return[i][ii] = make(map[string]interface{})
            for key, value := range ddata {
                _return[i][ii][key] = value
            }
        }
    }
    return _return, nil
}
