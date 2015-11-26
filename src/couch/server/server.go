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

// func (this *Server) Info() string {
func (this *Server) Info() map[string]interface{} {
    type __ struct {
        CouchDB string
        Uuid    string
        Version string
        Vendor  map[string]string
    }
    data, err := this.Client.Get("/", nil, nil).GetBody(&__{})
    if err != nil {
        return nil
    }
    var result = make(map[string]interface{});
    result["couchdb"] = data.(*__).CouchDB
    result["uuid"]    = data.(*__).Uuid
    result["version"] = data.(*__).Version
    result["vendor"]  = map[string]string{
           "name": data.(*__).Vendor["name"],
        "version": data.(*__).Vendor["version"],
    }
    return result
}
