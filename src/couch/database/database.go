package database

import _client "./../client"

import u "./../util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Database struct {
    Client *_client.Client
    Name   string
}

func Shutup() {
    u.Shutup()
}

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
    type Data map[string]interface{}
    data, err := this.Client.Get(this.Name, nil, nil).GetBodyData(&Data{})
    if err != nil {
        return nil, err
    }
    var _return = make(map[string]interface{})
    for key, value := range *data.(*Data) {
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
    var body = map[string]interface{}{
        "source": this.Name,
        "target": target,
        "create_target": targetCreate,
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
