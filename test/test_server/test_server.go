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
    data, err := _newServer().Info()
    if err != nil {
        panic(err)
    }
    _dumpf("Server Info >> %+v", data)
    _dumpf("Server Info >> couchdb: %s", data["couchdb"])
    // or
    _dumpf("Server Info >> couchdb: %s", u.Dig("couchdb", data))
    _dumpf("Server Info >> uuid: %s", u.Dig("uuid", data))
    _dumpf("Server Info >> version: %s", u.Dig("version", data))
    _dumpf("Server Info >> vendor.name: %s", u.Dig("vendor.name", data))
    _dumpf("Server Info >> vendor.version: %s", u.Dig("vendor.version", data))
}

/**
 * TestVersion
 */
func TestVersion() {
    data, err := _newServer().Version()
    if err != nil {
        panic(err)
    }
    _dumpf("Server Version >> %s", data)
}

/**
 * TestGetActiveTasks
 */
func TestGetActiveTasks() {
    data, err := _newServer().GetActiveTasks()
    if err != nil {
        panic(err)
    }
    _dumpf("Server Active Tasks >> %+v", data)
    _dumpf("Server Active Tasks >> 0.pid %s", data[0]["pid"])
    _dumpf("Server Active Tasks >> 0.database %s", data[0]["databases"])
}
