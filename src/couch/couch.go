// Copyright 2015 Kerem Güneş
//    <http://qeremy.com>
//
// Apache License, Version 2.0
//    <http://www.apache.org/licenses/LICENSE-2.0>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package couch

// @package couch
// @object  couch.Couch
// @author  Kerem Güneş <qeremy[at]gmail[dot]com>
type Couch struct {
    Config map[string]interface{}
}

const (
    NAME    = "Couch"
    VERSION = "1.0.0"
    // used to dump whole stream
    DEBUG   = false
)

// Suppress "... imported and not used"
func Shutup() {}

// Object constructor.
//
// @param  config map[string]interface{}
// @return *couch.Couch
func New(config interface{}, debug bool) *Couch {
    var this = &Couch{
        Config: map[string]interface{}{
            "debug": debug,
        },
    }

    // apply config options if provided
    if config, ok := config.(map[string]interface{}); ok {
        this.SetConfig(config)
    }

    return this
}

// Set config
//
// @param  config map[string]interface{}
// @return void
func (this *Couch) SetConfig(config map[string]interface{}) {
    this.Config = config
}

// Get config
//
// @return void
func (this *Couch) GetConfig() map[string]interface{} {
    return this.Config
}
