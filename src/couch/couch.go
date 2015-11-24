package couch

import _client "./../couch/client"
// import _stream "./http/stream"
// import _request "./http/request"
// import _response "./http/response"

import u "./util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Couch struct {
    config map[string]interface{}
}

const (
    NAME    = "Couch"
    VERSION = "1.0.0"
)

const (
    DEFAULT_SCHEME string = "http"
    DEFAULT_HOST   string = "localhost"
    DEFAULT_PORT   uint16 = 5984
)

func New(config interface{}) *Couch {
    couch := &Couch{}
    if config, ok := config.(map[string]interface{}); ok {
        couch.SetConfig(config)
    }
    return couch
}

func (this *Couch) SetConfig(config map[string]interface{}) {
    this.config = config
}
func (this *Couch) GetConfig() map[string]interface{} {
    return this.config
}

func NewClient(couch *Couch, scheme, host interface{}, port interface{}) *_client.Client {
    scheme = u.IsEmptySet(scheme, DEFAULT_SCHEME)
    host   = u.IsEmptySet(host, DEFAULT_HOST)
    port   = u.IsEmptySet(port, DEFAULT_PORT)
    client := &_client.Client{
        Scheme: u.String(scheme),
        Host: host.(string),
        Port: port.(uint16),
        Couch: map[string]interface{}{
            "NAME": NAME,
            "VERSION": VERSION,
        },
    }
    var config = couch.GetConfig()
    if value, _ := config["Host"].(string); value != "" {
        client.Host = value
    }
    if value, _ := config["Port"].(uint16); value != 0 {
        client.Port = value
    }
    if value, _ := config["Username"].(string); value != "" {
        client.Username = value
    }
    if value, _ := config["Password"].(string); value != "" {
        client.Password = value
    }
    return client;
}
