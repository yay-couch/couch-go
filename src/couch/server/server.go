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
    type Data map[string]interface{}
    data, err := this.Client.Get("/", nil, nil).GetData(&Data{})
    if err != nil {
        return nil, err
    }
    var _return = make(map[string]interface{})
    for key, value := range *data.(*Data) {
        switch value := value.(type) {
            case map[string]interface{}:
                _return[key] = make(map[string]string)
                for kkey, vvalue := range value {
                    _return[key].(map[string]string)[kkey] = vvalue.(string)
                }
            default:
                _return[key] = value
        }
    }
    return _return, nil
}

func (this *Server) Version() (string, error) {
    data, err := this.Info()
    if err != nil {
        return "", err
    }
    return data["version"].(string), nil
}

func (this *Server) GetActiveTasks() ([]map[string]interface{}, error) {
    type Data []map[string]interface{}
    data, err := this.Client.Get("/_active_tasks", nil, nil).GetData(&Data{})
    if err != nil {
        return nil, err
    }
    var _return = make([]map[string]interface{}, len(*data.(*Data)))
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

func (this *Server) GetUuid() (string, error) {
    data, err := this.GetUuids(1)
    if err != nil {
        return "", err
    }
    return data[0], nil
}

func (this *Server) GetUuids(count int) ([]string, error) {
    type Data map[string][]string
    data, err := this.Client.Get("/_uuids", map[string]interface{}{
        "count": count,
    }, nil).GetData(&Data{})
    if err != nil {
        return nil, err
    }
    var _return = make([]string, count)
    for _, uuid := range *data.(*Data) {
        for i := 0; i < len(uuid); i++ {
            _return[i] = uuid[i]
        }
    }
    return _return, nil
}

func (this *Server) Replicate(body interface{}) (map[string]interface{}, error) {
    if body == nil {
        body = make(map[string]interface{})
    }
    if body.(map[string]interface{})["source"] == nil ||
       body.(map[string]interface{})["target"] == nil {
        panic("Both source & target required!")
    }
    type Data map[string]interface{}
    data, err := this.Client.Post("/_replicate", nil, body, nil).GetData(&Data{})
    if err != nil {
        return nil, err
    }
    var _return = make(map[string]interface{})
    for key, value := range *data.(*Data) {
        if key == "history" {
            _return[key] = make(map[int]map[string]interface{})
            for i, history := range value.([]interface{}) {
                _return[key] = make([]map[string]interface{}, len(value.([]interface{})))
                for kkey, vvalue := range history.(map[string]interface{}) {
                    if _return[key].([]map[string]interface{})[i] == nil {
                        _return[key].([]map[string]interface{})[i] = make(map[string]interface{})
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
