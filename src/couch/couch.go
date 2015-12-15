package couch

type Couch struct {
    Config map[string]interface{}
}

const (
    NAME    = "Couch"
    VERSION = "1.0.0"
    DEBUG   = false
)

func Shutup() {}

func New(config interface{}, debug bool) *Couch {
    var couch = &Couch{
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
