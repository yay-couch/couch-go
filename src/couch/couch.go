package couch

import (
    "./util"
)

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

func New(config interface{}, debug bool) *Couch {
    couch := &Couch{
        Config: map[string]interface{}{
            "debug": debug,
        },
    }
    if config, ok := config.(map[string]interface{}); ok {
        couch.SetConfig(config)
    }
    return couch
}

func NewClient(couch *Couch, config interface{}) *Client {
    var Config = make(map[string]interface{})
    Config["Couch.NAME"]    = NAME
    Config["Couch.VERSION"] = VERSION
    Config["Couch.DEBUG"]   = DEBUG // set default
    if config != nil {
        for key, value := range config.(map[string]interface{}) {
            Config[key] = value
        }
    }
    Config["Scheme"]   = util.IsEmptySet(Config["Scheme"],   DefaultScheme)
    Config["Host"]     = util.IsEmptySet(Config["Host"],     DefaultHost)
    Config["Port"]     = util.IsEmptySet(Config["Port"],     DefaultPort)
    Config["Username"] = util.IsEmptySet(Config["Username"], Username)
    Config["Password"] = util.IsEmptySet(Config["Password"], Password)

    if debug := util.Dig("debug", couch.Config); debug != nil {
        Config["Couch.DEBUG"] = debug
    }

    couch.SetConfig(Config)

    return _client.New(Config,
        Config["Username"].(string), Config["Username"].(string))
    // or
    // return _client.New("https://localhost:1234", "", "", Config???)
}

func (this *Couch) SetConfig(config map[string]interface{}) {
    this.Config = config
}
func (this *Couch) GetConfig() map[string]interface{} {
    return this.Config
}
