package test

import (
   "couch"
   "couch/util"
)

var (
   DEBUG = true
)

var (
   Couch  *couch.Couch
   Client *couch.Client
   Server *couch.Server
)

func init() {
   Couch  = couch.New(nil, DEBUG)
   Client = couch.NewClient(Couch)
   Server = couch.NewServer(Client)
}

/**
 * TestAll
 */
func TestAll() {}

/**
 * TestPing
 */
func TestPing() {
   util.Dumpf("Server Ping >> %v ", Server.Ping())
}

/**
 * TestInfo
 */
func TestInfo() {
   data, err := Server.Info()
   if err != nil {
      panic(err)
   }
   util.Dumpf("Server Info >> %+v", data)
   util.Dumpf("Server Info >> couchdb: %s", data["couchdb"])
   // or
   util.Dumpf("Server Info >> couchdb: %s", util.Dig("couchdb", data))
   util.Dumpf("Server Info >> uuid: %s", util.Dig("uuid", data))
   util.Dumpf("Server Info >> version: %s", util.Dig("version", data))
   // or
   util.Dumpf("Server Info >> vendor.name: %s", util.Dig("vendor.name", data))
   util.Dumpf("Server Info >> vendor.version: %s", util.Dig("vendor.version", data))
   util.Dumpf("Server Info >> vendor.name: %s", data["vendor"].(map[string]string)["name"])
   util.Dumpf("Server Info >> vendor.version: %s", data["vendor"].(map[string]string)["version"])
}

/**
 * TestVersion
 */
func TestVersion() {
   data, err := Server.Version()
   if err != nil {
      panic(err)
   }
   util.Dumpf("Server Version >> %s", data)
}

/**
 * TestGetActiveTasks
 */
func TestGetActiveTasks() {
   data, err := Server.GetActiveTasks()
   if err != nil {
      panic(err)
   }
   util.Dumpf("Server Active Tasks >> %+v", data)
   for _, task := range data {
      util.Dumpf("Server Active Tasks >> 0.pid: %s", task["pid"])
      util.Dumpf("Server Active Tasks >> 0.database: %s", task["database"])
   }
}

/**
 * TestGetAllDatabases
 */
func TestGetAllDatabases() {
   data, err := Server.GetAllDatabases()
   if err != nil {
      panic(err)
   }
   util.Dumpf("Server Databases >> %+v", data)
   util.Dumpf("Server Databases >> first: %s", data[0])
}

/**
 * TestGetDatabaseUpdates
 */
func TestGetDatabaseUpdates() {
   data, err := Server.GetDatabaseUpdates(nil)
   if err != nil {
      panic(err)
   }
   util.Dumpf("Server Updates >> %+v", data)
   util.Dumpf("Server Updates >> db_name: %s", data["db_name"])
   util.Dumpf("Server Updates >> type: %s", data["type"])
   util.Dumpf("Server Updates >> ok: %v", data["ok"])
}

/**
 * TestGetLogs
 */
func TestGetLogs() {
   var data = Server.GetLogs(nil)
   util.Dumps(data)
}

/**
 * TestGetStats
 */
func TestGetStats() {
   data, err := Server.GetStats("")
   if err != nil {
      panic(err)
   }
   util.Dumpf("Server Stats >> %+v", data)
   util.Dumpf("Server Stats >> couchdb: %+v", data["couchdb"])
   util.Dumpf("Server Stats >> couchdb.request_time: %+v", data["couchdb"]["request_time"])
   util.Dumpf("Server Stats >> couchdb.request_time.description: %s", data["couchdb"]["request_time"]["description"])
   util.Dumpf("Server Stats >> couchdb.request_time.description: %f", data["couchdb"]["request_time"]["current"])
   util.Dumpf("Server Stats >> httpd_request_methods.GET.max: %v", data["httpd_request_methods"]["GET"]["max"])
}

/**
 * TestGetUuid
 */
func TestGetUuid() {
   data, err := Server.GetUuid()
   if err != nil {
      panic(err)
   }
   util.Dumpf("Server Uuid >> %s", data)
}

/**
 * TestGetUuids
 */
func TestGetUuids() {
   data, err := Server.GetUuids(3)
   if err != nil {
      panic(err)
   }
   util.Dumpf("Server Uuids >> %+v", data)
   for i, _ := range data {
      util.Dumpf("Server Uuids >> %d: %s", i, data[i])
   }
}

/**
 * TestReplicate
 */
func TestReplicate() {
   data, err := Server.Replicate(map[string]interface{}{
      "source": "foo",
      "target": "foo_replicate",
      "create_target": true,
   })
   if err != nil {
      panic(err)
   }
   util.Dumpf("Server Database Replicate >> %+v", data)
   util.Dumpf("Server Database Replicate >> ok: %v", data["ok"])
   util.Dumpf("Server Database Replicate >> history.0: %v", util.Dig("0", data["history"]))
   util.Dumpf("Server Database Replicate >> history.0.start_time: %s", util.Dig("0.start_time", data["history"]))
}

/**
 * TestRestart
 */
func TestRestart() {
   util.Dumpf("Server Restart >> %v ", Server.Restart())
}

/**
 * TestGetConfig
 */
func TestGetConfig() {
   // data1, err := Server.GetConfig()
   // if err != nil {
   //     panic(err)
   // }
   // util.Dumpf("Server Get Config >> %+v", data1)
   // util.Dumpf("Server Get Config >> couchdb: %v", data1["couchdb"])
   // util.Dumpf("Server Get Config >> couchdb.uuid: %s", data1["couchdb"]["uuid"])
   // // or
   // data2, err := Server.GetConfigSection("couchdb")
   // if err != nil {
   //     panic(err)
   // }
   // util.Dumpf("Server Get Config Section >> couchdb: %v", data2)
   // util.Dumpf("Server Get Config Section >> couchdb.uuid: %s", data2["uuid"])
   // or
   data3, err := Server.GetConfigSectionKey("couchdb", "uuid")
   if err != nil {
      panic(err)
   }
   util.Dumpf("Server Get Config Section Key >> couchdb.uuid: %s", data3)
}

/**
 * TestSetConfig
 */
func TestSetConfig() {
   data, err := Server.SetConfig("couch", "foo", "The foo!")
   if err != nil {
      panic(err)
   }
   util.Dumpf("Server Set Config >> %s", data)
}

/**
 * TestRemoveConfig
 */
func TestRemoveConfig() {
   data, err := Server.RemoveConfig("couch", "foo")
   if err != nil {
      panic(err)
   }
   util.Dumpf("Server Remove Config >> %s", data)
}
