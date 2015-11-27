package test_server

import _couch    "./../../src/couch"
import _client   "./../../src/couch/client"
import _stream   "./../../src/couch/http/stream"
import _request  "./../../src/couch/http/request"
import _response "./../../src/couch/http/response"

import _server   "./../../src/couch/server"

import u "./../../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

func Shutup() {
    _client.Shutup()
    _stream.Shutup()
    _request.Shutup()
    _response.Shutup()
    _server.Shutup()
}

func _newServer() *_server.Server {
    couch  := _couch.New(nil)
    client := _couch.NewClient(couch, nil)
    server := _couch.NewServer(client)
    return server
}

/**
 * TestAll
 */
func TestAll() {
    TestInfo()
}

/**
 * TestInfo
 */
func TestInfo() {
    var info = _newServer().Info()
    _dumpf("Server Info >> %+v", info)
    _dumpf("Server Info >> couchdb: %s", info["couchdb"])
    // or
    _dumpf("Server Info >> couchdb: %s", u.Dig("couchdb", info))
    _dumpf("Server Info >> uuid: %s", u.Dig("uuid", info))
    _dumpf("Server Info >> version: %s", u.Dig("version", info))
    _dumpf("Server Info >> vendor.name: %s", u.Dig("vendor.name", info))
    _dumpf("Server Info >> vendor.version: %s", u.Dig("vendor.version", info))
}

/**
 * TestVersion
 */
func TestVersion() {
    var version = _newServer().Version()
    _dumpf("Server Version >> %s", version)
}

/**
 * TestGetActiveTasks
 */
func TestGetActiveTasks() {
    var tasks = _newServer().GetActiveTasks()
    _dumpf("Server Active Tasks >> %+v", tasks)
    _dumpf("Server Active Tasks >> 0.pid %s", tasks[0]["pid"])
    _dumpf("Server Active Tasks >> 0.database %s", tasks[0]["databases"])
}
