package test_database

import _couch    "./../../src/couch"
import _client   "./../../src/couch/client"
import _database "./../../src/couch/database"

import u "./../../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

var (
    DEBUG  = true
    DBNAME = "foo"
)

var (
    couch    *_couch.Couch
    client   *_client.Client
    database *_database.Database
)

func init() {
    couch    = _couch.New(nil, DEBUG)
    client   = _couch.NewClient(couch, nil)
    database = _couch.NewDatabase(client, DBNAME);
}

/**
 * TestAll
 */
func TestAll() {}

/**
 * TestPing
 */
func TestPing() {
    _dumpf("Database Ping >> %v", database.Ping())
}

/**
 * TestInfo
 */
func TestInfo() {
    data, err := database.Info()
    if err != nil {
        panic(err)
    }
    _dumpf("Database Info >> %+v", data)
    _dumpf("Database Info >> db_name: %s", data["db_name"])
    for key, value := range data {
        _dumpf("Database Info >> %s: %v", key, value)
    }
}

/**
 * TestCreate
 */
func TestCreate() {
    _dumpf("Database Create >> %v", database.Create())
    // error?
    // if err := client.GetResponse().GetError(); err != "" {
    //     _dumpf("Response Status: %s", client.GetResponse().GetStatus())
    //     _dumpf("Response Body  : %s", client.GetResponse().GetBody())
    //     panic(err)
    // }
}

/**
 * TestRemove
 */
func TestRemove() {
    _dumpf("Database Remove >> %v", database.Remove())
    // error?
    // if err := client.GetResponse().GetError(); err != "" {
    //     _dumpf("Response Status: %s", client.GetResponse().GetStatus())
    //     _dumpf("Response Body  : %s", client.GetResponse().GetBody())
    //     panic(err)
    // }
}

/**
 * TestReplicate
 */
func TestReplicate() {
    data, err := database.Replicate("foo_replicate", true)
    if err != nil {
        panic(err)
    }
    _dumpf("Server Database Replicate >> %+v", data)
    _dumpf("Server Database Replicate >> ok: %v", data["ok"])
    _dumpf("Server Database Replicate >> history.0: %v", u.Dig("0", data["history"]))
    _dumpf("Server Database Replicate >> history.0.start_time: %s", u.Dig("0.start_time", data["history"]))
}
