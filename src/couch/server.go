package couch

import (
    "./util"
)

type Server struct {
    Client *Client
}

func NewServer(client *Client) *Server {
    return &Server{
        Client: client,
    }
}

func (this *Server) Ping() bool {
    return (200 == this.Client.Head("/", nil, nil).GetStatusCode())
}

func (this *Server) Info() (map[string]interface{}, error) {
    data, err := this.Client.Get("/", nil, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.Map()
    for key, value := range data.(map[string]interface{}) {
        switch value := value.(type) {
            case map[string]interface{}:
                _return[key] = util.MapString()
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
    data, err := this.Client.Get("/_active_tasks", nil, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.MapList(data)
    for i, data := range data.([]interface{}) {
        _return[i] = util.Map()
        for key, value := range data.(map[string]interface{}) {
            _return[i][key] = value
        }
    }
    return _return, nil
}

func (this *Server) GetAllDatabases() ([]string, error) {
    data, err := this.Client.Get("/_all_dbs", nil, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = util.MapStringSlice(data)
    for i, db := range data.([]interface{}) {
        _return[i] = db.(string)
    }
    return _return, nil
}

func (this *Server) GetDatabaseUpdates(query interface{}) (map[string]interface{}, error) {
    data, err := this.Client.Get("/_db_updates", query, nil).GetBodyData(nil)
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
    return this.Client.Get("/_log", query, nil).GetBody()
}

func (this *Server) GetStats(path string) (map[string]map[string]map[string]interface{}, error) {
    data, err := this.Client.Get("/_stats/"+ path, nil, nil).GetBodyData(nil)
    if err != nil {
        return nil, err
    }
    var _return = make(map[string]map[string]map[string]interface{})
    for i, data := range data.(map[string]interface{}) {
        _return[i] = make(map[string]map[string]interface{})
        for ii, ddata := range data.(map[string]interface{}) {
            _return[i][ii] = make(map[string]interface{})
            for key, value := range ddata.(map[string]interface{}) {
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
    }, nil).GetBodyData(&Data{})
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
    data, err := this.Client.Post("/_replicate", nil, body, nil).GetBodyData(&Data{})
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

func (this *Server) Restart() bool {
    return (202 == this.Client.Post("/_restart", nil, nil, nil).GetStatusCode())
}

func (this *Server) GetConfig() (map[string]map[string]interface{}, error) {
    type Data map[string]map[string]interface{}
    data, err := this.Client.Get("/_config", nil, nil).GetBodyData(&Data{})
    if err != nil {
        return nil, err
    }
    var _return = make(map[string]map[string]interface{})
    for key, value := range *data.(*Data) {
        if _return[key] == nil {
            _return[key] = make(map[string]interface{})
        }
        for kkey, vvalue := range value {
            _return[key][kkey] = vvalue
        }
    }
    return _return, nil
}
func (this *Server) GetConfigSection(section string) (map[string]interface{}, error) {
    data, err := this.GetConfig()
    if err != nil {
        return nil, err
    }
    return data[section], nil
}
func (this *Server) GetConfigSectionKey(section string, key string) (string, error) {
    data, err := this.GetConfig()
    if err != nil {
        return "", err
    }
    return data[section][key].(string), nil
}

func (this *Server) SetConfig(section, key, value string) (string, error) {
    if section == "" || key == "" {
        panic("Both section & key required!")
    }
    var Data string
    data, err := this.Client.Put("/_config/"+ section +"/"+ key, nil, value, nil).GetBodyData(Data)
    if err != nil {
        return "", err
    }
    return data.(string), nil
}

func (this *Server) RemoveConfig(section, key string) (string, error) {
    if section == "" || key == "" {
        panic("Both section & key required!")
    }
    var Data string
    data, err := this.Client.Delete("/_config/"+ section +"/"+ key, nil, nil).GetBodyData(Data)
    if err != nil {
        return "", err
    }
    return data.(string), nil
}
