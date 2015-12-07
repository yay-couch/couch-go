package couch

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

func (this *Couch) SetConfig(config map[string]interface{}) {
    this.Config = config
}
func (this *Couch) GetConfig() map[string]interface{} {
    return this.Config
}
