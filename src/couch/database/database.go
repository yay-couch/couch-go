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
