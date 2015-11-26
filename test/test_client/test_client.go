package test_client

import _couch    "./../../src/couch"
import _client   "./../../src/couch/client"
import _stream   "./../../src/couch/http/stream"
import _request  "./../../src/couch/http/request"
import _response "./../../src/couch/http/response"

import u "./../../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

func Shutup() {
    _client.Shutup()
    _stream.Shutup()
    _request.Shutup()
    _response.Shutup()
}

func _doRequest(uri string) *_response.Response {
    couch  := _couch.New(nil)
    client := _couch.NewClient(couch, nil)
    // or
    // client := _couch.NewClient(couch, map[string]interface{}{
    //     "Host": "127.0.0.1",
    // })
    return client.DoRequest(uri, nil, "", nil)
}

/**
 * TestAll
 */
func TestAll() {
    // status
    TestClientResponseStatus()
    TestClientResponseStatusCode()
    TestClientResponseStatusText()
    _dump("")

    // headers
    TestClientResponseHeaders("")
    TestClientResponseHeaders("0") // status line
    TestClientResponseHeaders("Server")
    _dump("")

    // body
    TestClientResponseBody()
    _dump("")

    // body parsed
    TestClientResponseBodyParse()
}

/**
 * TestClientResponseStatus
 */
func TestClientResponseStatus() {
    var response = _doRequest("GET /")
    _dumpf("Response Status >> %s", response.GetStatus())
}

/**
 * TestClientResponseStatusCode
 */
func TestClientResponseStatusCode() {
    var response = _doRequest("GET /")
    _dumpf("Response Status Code >> %d", response.GetStatusCode())
}

/**
 * TestClientResponseStatusText
 */
func TestClientResponseStatusText() {
    var response = _doRequest("GET /")
    _dumpf("Response Status Text >> %s", response.GetStatusText())
}

/**
 * TestClientResponseHeaders.
 */
func TestClientResponseHeaders(key string) {
    var response = _doRequest("GET /")
    if key == "" {
        _dumpf("Response Headers >> %+v", response.GetHeaderAll())
    } else {
        _dumpf("Response Headers >> %s: %+v", key, response.GetHeader(key))
    }
}

/**
 * TestClientResponseBody
 */
func TestClientResponseBody() {
    var response = _doRequest("GET /")
    _dumpf("Response Body >> len: %d body: %+v",
        len(response.GetBody()), response.GetBody())
}

/**
 * TestClientResponseBodyParse.
 */
func TestClientResponseBodyParse() {
    type Response struct {
        CouchDB string
        Uuid    string
        Version string
        Vendor  map[string]string
    }

    var response = _doRequest("GET /")
    data, err := u.ParseBody(response.GetBody(), &Response{})
    if err != nil {
        _dumps(err)
        return
    }
    _dumpf("Response Body Parsed >> type: %T value: %+v", data, data)
    _dumpf("Response Body Parsed >> couchdb: %s", data.(*Response).CouchDB)
    _dumpf("Response Body Parsed >> uuid: %s", data.(*Response).Uuid)
    _dumpf("Response Body Parsed >> version: %s", data.(*Response).Version)
    _dumpf("Response Body Parsed >> vendor: %s", data.(*Response).Vendor)
    _dumpf("Response Body Parsed >> vendor.name: %s", data.(*Response).Vendor["name"])
    _dumpf("Response Body Parsed >> vendor.version: %s", data.(*Response).Vendor["version"])
}
