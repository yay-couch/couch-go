// Copyright 2015 Kerem Güneş
//   <k-gun@mail.com>
//
// Apache License, Version 2.0
//   <http://www.apache.org/licenses/LICENSE-2.0>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// @package couch
// @uses    couch.util
// @author  Kerem Güneş <k-gun@mail.com>
package couch

import (
   "couch/util"
)

// @object couch.Server
type Server struct {
   Client *Client
}

// Constructor.
//
// @param  client *couch.Client
// @return *couch.Server
func NewServer(client *Client) (*Server) {
   return &Server{
      Client: client,
   }
}

// Ping.
//
// @return bool
func (this *Server) Ping() (bool) {
   return (200 == this.Client.Head("/", nil, nil).GetStatusCode())
}

// Info.
//
// @return map[string]interface{}, error
func (this *Server) Info() (map[string]interface{}, error) {
   data, err := this.Client.Get("/", nil, nil).GetBodyData(nil)
   if err != nil {
      return nil, err
   }

   ret := util.Map()
   for key, value := range data.(map[string]interface{}) {
      switch value := value.(type) {
         case map[string]interface{}:
            ret[key] = util.MapString()
            for kkey, vvalue := range value {
               ret[key].(map[string]string)[kkey] = vvalue.(string)
            }
         default:
            ret[key] = value
      }
   }

   return ret, nil
}

// Version.
//
// @return string, error
func (this *Server) Version() (string, error) {
   data, err := this.Info()
   if err != nil {
      return "", err
   }

   return data["version"].(string), nil
}

// Get active tasks.
//
// @return []map[string]interface{}, error
func (this *Server) GetActiveTasks() ([]map[string]interface{}, error) {
   data, err := this.Client.Get("/_active_tasks", nil, nil).GetBodyData(nil)
   if err != nil {
      return nil, err
   }

   ret := util.MapList(data)
   for i, data := range data.([]interface{}) {
      ret[i] = util.Map()
      for key, value := range data.(map[string]interface{}) {
         ret[i][key] = value
      }
   }

   return ret, nil
}

// Get all databases.
//
// @return []string, error
func (this *Server) GetAllDatabases() ([]string, error) {
   data, err := this.Client.Get("/_all_dbs", nil, nil).GetBodyData(nil)
   if err != nil {
      return nil, err
   }

   ret := util.MapSliceString(data)
   for i, db := range data.([]interface{}) {
      ret[i] = db.(string)
   }

   return ret, nil
}

// Get database updates.
//
// @param  query map[string]interface{}
// @return map[string]interface{}, error
func (this *Server) GetDatabaseUpdates(query interface{}) (map[string]interface{}, error) {
   data, err := this.Client.Get("/_db_updates", query, nil).GetBodyData(nil)
   if err != nil {
      return nil, err
   }

   return map[string]interface{}{
          "ok": util.DigBool("ok", data),
         "type": util.DigString("type", data),
      "db_name": util.DigString("db_name", data),
   }, nil
}

// Get logs.
//
// @param  query map[string]interface{}
// @return string
func (this *Server) GetLogs(query interface{}) (string) {
   return this.Client.Get("/_log", query, nil).GetBody()
}

// Get stats.
//
// @param  path string
// @return map[string]map[string]map[string]interface{}, error
func (this *Server) GetStats(path string) (map[string]map[string]map[string]interface{}, error) {
   data, err := this.Client.Get("/_stats/"+ path, nil, nil).GetBodyData(nil)
   if err != nil {
      return nil, err
   }

   ret := make(map[string]map[string]map[string]interface{})
   for i, data := range data.(map[string]interface{}) {
      ret[i] = make(map[string]map[string]interface{})
      for ii, ddata := range data.(map[string]interface{}) {
         ret[i][ii] = make(map[string]interface{})
         for key, value := range ddata.(map[string]interface{}) {
            ret[i][ii][key] = value
         }
      }
   }

   return ret, nil
}

// Get UUID.
//
// @return string, error
func (this *Server) GetUuid() (string, error) {
   data, err := this.GetUuids(1)
   if err != nil {
      return "", err
   }

   return data[0], nil
}

// Get UUIDs.
//
// @return []string, error
func (this *Server) GetUuids(count int) ([]string, error) {
   query := util.ParamList(
      "count", count,
   )
   data, err := this.Client.Get("/_uuids", query, nil).GetBodyData(nil)
   if err != nil {
      return nil, err
   }

   ret := util.MapSliceString(count)
   for i, uuid := range data.(map[string]interface{})["uuids"].([]interface{}) {
      ret[i] = uuid.(string)
   }

   return ret, nil
}

// Get UUID.
//
// @param  body map[string]interface{}
// @return map[string]interface{}, error
// @panics
func (this *Server) Replicate(body map[string]interface{}) (map[string]interface{}, error) {
   body = util.Param(body)
   if body["source"] == nil || body["target"] == nil {
      panic("Both source & target required!")
   }

   data, err := this.Client.Post("/_replicate", nil, body, nil).GetBodyData(nil)
   if err != nil {
      return nil, err
   }

   ret := util.Map()
   for key, value := range data.(map[string]interface{}) {
      // grap, set & pass history field
      if key == "history" {
         ret["history"] = util.MapList(value)
         for i, history := range value.([]interface{}) {
            ret["history"].([]map[string]interface{})[i] = util.Map()
            for kkey, vvalue := range history.(map[string]interface{}) {
               ret["history"].([]map[string]interface{})[i][kkey] = vvalue
            }
         }
         continue
      }
      ret[key] = value
   }

   return ret, nil
}

// Restart.
//
// @return bool
func (this *Server) Restart() (bool) {
   return (202 == this.Client.Post("/_restart", nil, nil, nil).GetStatusCode())
}

// Get config.
//
// @return map[string]map[string]interface{}, error
func (this *Server) GetConfig() (map[string]map[string]interface{}, error) {
   data, err := this.Client.Get("/_config", nil, nil).GetBodyData(nil)
   if err != nil {
      return nil, err
   }

   ret := util.MapMapString()
   for key, value := range data.(map[string]interface{}) {
      ret[key] = util.Map()
      for kkey, vvalue := range value.(map[string]interface{}) {
         ret[key][kkey] = vvalue
      }
   }

   return ret, nil
}

// Get config section.
//
// @param  section string
// @return map[string]map[string]interface{}, error
func (this *Server) GetConfigSection(section string) (map[string]interface{}, error) {
   data, err := this.GetConfig()
   if err != nil {
      return nil, err
   }

   return data[section], nil
}

// Get config section value bey key.
//
// @param  section string
// @param  key string
// @return string, error
func (this *Server) GetConfigSectionKey(section string, key string) (string, error) {
   data, err := this.GetConfig()
   if err != nil {
      return "", err
   }

   return data[section][key].(string), nil
}

// Set config.
//
// @param  section string
// @param  key    string
// @param  value   string
// @return string, error
// @panics
func (this *Server) SetConfig(section, key, value string) (string, error) {
   if section == "" || key == "" {
      panic("Both section & key required!")
   }

   data, err := this.Client.Put("/_config/"+ section +"/"+ key, nil, value, nil).
      GetBodyData(nil)
   if err != nil {
      return "", err
   }

   return data.(string), nil
}

// Remove config.
//
// @param  section string
// @param  key    string
// @return string, error
// @panics
func (this *Server) RemoveConfig(section, key string) (string, error) {
   if section == "" || key == "" {
      panic("Both section & key required!")
   }

   data, err := this.Client.Delete("/_config/"+ section +"/"+ key, nil, nil).
      GetBodyData(nil)
   if err != nil {
      return "", err
   }

   return data.(string), nil
}
