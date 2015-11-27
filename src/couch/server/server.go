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

func (this *Server) Info() map[string]interface{} {
    type Data struct {
        CouchDB string
        Uuid    string
        Version string
        Vendor  map[string]string
    }
    data, err := this.Client.Get("/", nil, nil).GetBody(&Data{})
    if err != nil {
        return nil
    }
    var info = make(map[string]interface{});
    info["couchdb"] = data.(*Data).CouchDB
    info["uuid"]    = data.(*Data).Uuid
    info["version"] = data.(*Data).Version
    info["vendor"]  = map[string]string{
           "name": data.(*Data).Vendor["name"],
        "version": data.(*Data).Vendor["version"],
    }
    return info
}
