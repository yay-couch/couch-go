package couch

import _client "./../couch/client"
import _server "./../couch/server"
// import _stream "./http/stream"
// import _request "./http/request"
// import _response "./http/response"

import u "./util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Couch struct {
    Config map[string]interface{}
}

const (
    NAME    = "Couch"
    VERSION = "1.0.0"
    DEBUG   = false
)

var (
    Username = ""
    Password = ""
)

var (
    DefaultScheme        = "http"
    DefaultHost          = "localhost"
    DefaultPort   uint16 = 5984
)

func Shutup() {}

func New(config interface{}) *Couch {
    couch := &Couch{}
    if config, ok := config.(map[string]interface{}); ok {
        couch.SetConfig(config)
    }
    return couch
}

func NewClient(couch *Couch, config interface{}) *_client.Client {
    var Config = make(map[string]interface{})
    Config["Couch.NAME"]    = NAME
    Config["Couch.VERSION"] = VERSION
    Config["Couch.DEBUG"]   = DEBUG // set default
    if config != nil {
        for key, value := range config.(map[string]interface{}) {
            Config[key] = value
        }
    }
    Config["Scheme"]   = u.IsEmptySet(Config["Scheme"],   DefaultScheme)
    Config["Host"]     = u.IsEmptySet(Config["Host"],     DefaultHost)
    Config["Port"]     = u.IsEmptySet(Config["Port"],     DefaultPort)
    Config["Username"] = u.IsEmptySet(Config["Username"], Username)
    Config["Password"] = u.IsEmptySet(Config["Password"], Password)

    if debug := u.Dig("debug", couch.Config); debug != nil {
        Config["Couch.DEBUG"] = debug
    }

    couch.SetConfig(Config)

    return _client.New(Config,
        Config["Username"].(string), Config["Username"].(string))
    // or
    // return _client.New("https://localhost:1234", "", "", Config???)
}

func NewServer(client *_client.Client) *_server.Server {
    return _server.New(client)
}

func (this *Couch) SetConfig(config map[string]interface{}) {
    this.Config = config
}
func (this *Couch) GetConfig() map[string]interface{} {
    return this.Config
}
