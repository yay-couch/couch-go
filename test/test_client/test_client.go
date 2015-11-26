package test

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
    body, err := u.ParseBody(response.GetBody(), &Response{})
    if err != nil {
        _dumps(err)
        return
    }
    _dumpf("Response Body Parsed >> type: %T value: %+v", body, body)
    _dumpf("Repsonse Body Parsed >> couchdb: %s", body.(*Response).CouchDB)
    _dumpf("Repsonse Body Parsed >> uuid: %s", body.(*Response).Uuid)
    _dumpf("Repsonse Body Parsed >> version: %s", body.(*Response).Version)
    _dumpf("Repsonse Body Parsed >> vendor: %s", body.(*Response).Vendor)
    _dumpf("Repsonse Body Parsed >> vendor.name: %s", body.(*Response).Vendor["name"])
    _dumpf("Repsonse Body Parsed >> vendor.version: %s", body.(*Response).Vendor["version"])
}
