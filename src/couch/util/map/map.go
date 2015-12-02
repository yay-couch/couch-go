package util

type Map struct {
    Data map[string]interface{}
}

func New() *Map {
    return &Map{
        Data: make(map[string]interface{}),
    }
}

func (this *Map) Set(key string, value interface{}) *Map {
    this.Data[key] = value
    return this
}

func (this *Map) Get(key string) interface{} {
    return this.Data[key]
}
