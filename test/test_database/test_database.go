package test_database

import _couch    "./../../src/couch"
import _database "./../../src/couch/database"

import u "./../../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

var (
    DEBUG  = true
    DBNAME = "foo_tmp"
)

func _newDatabase() *_database.Database {
    couch    := _couch.New(nil, DEBUG)
    client   := _couch.NewClient(couch, nil)
    database := _couch.NewDatabase(client, DBNAME);
    return database
}

/**
 * TestAll
 */
func TestAll() {}

/**
 * TestPing
 */
func TestPing() {
    _dumpf("Database Ping >> %v", _newDatabase().Ping())
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
    _dumpf("Database Info >> db_name: %s", data["db_name"])
    for key, value := range data {
        _dumpf("Database Info >> %s: %v", key, value)
    }
}

/**
 * TestCreate
 */
func TestCreate() {
    _dumpf("Database Create >> %v", _newDatabase().Create())
}
