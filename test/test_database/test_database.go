package test_database

import _couch    "./../../src/couch"
import _database "./../../src/couch/database"

import u "./../../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

var (
    DEBUG  = true
    DBNAME = "foo"
)

func Shutup() {}

func _newDatabase() *_database.Database {
    couch    := _couch.New(nil, DEBUG)
    client   := _couch.NewClient(couch, nil)
    database := _couch.NewDatabase(client, DBNAME);
    return database
}

/**
 * TestAll
 */
func TestAll() {
}

/**
 * TestPing
 */
func TestPing() {
    _dumpf("Database Ping >> %v ", _newDatabase().Ping())
}

/**
 * TestInfo
 */
func TestInfo() {
    data, err := _newDatabase().Info()
    if err != nil {
        panic(err)
    }
    _dumpf("Database Info >> %+v", data)
    // _dumpf("Database Info >> couchdb: %s", data["couchdb"])
    // // or
    // _dumpf("Database Info >> couchdb: %s", u.Dig("couchdb", data))
    // _dumpf("Database Info >> uuid: %s", u.Dig("uuid", data))
    // _dumpf("Database Info >> version: %s", u.Dig("version", data))
    // // or
    // _dumpf("Database Info >> vendor.name: %s", u.Dig("vendor.name", data))
    // _dumpf("Database Info >> vendor.version: %s", u.Dig("vendor.version", data))
    // _dumpf("Database Info >> vendor.name: %s", data["vendor"].(map[string]string)["name"])
    // _dumpf("Database Info >> vendor.version: %s", data["vendor"].(map[string]string)["version"])
}
