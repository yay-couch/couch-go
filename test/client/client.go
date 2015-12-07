package client

// import _couch    "./../../src/couch"
// import _client   "./../../src/couch/client"
// import _stream   "./../../src/couch/http/stream"
// import _request  "./../../src/couch/http/request"
// import _response "./../../src/couch/http/response"

import (
    "./../../src/couch"
    "./../../src/couch/util"
    "./../../src/couch/http"
)

var (
    DEBUG = true
)

var (
    Couch  *couch.Couch
    Client *couch.Client
)

func init() {
    Couch  = couch.New(nil, DEBUG)
    Client = couch.NewClient(Couch)
}

func _doRequest(uri string) *http.Response {
    return Client.DoRequest(uri, nil, "", nil)
}

/**
 * TestAll
 */
func TestAll() {
    // status
    TestClientResponseStatus()
    TestClientResponseStatusCode()
    TestClientResponseStatusText()
    util.Dump("")

    // headers
    TestClientResponseHeaders("")
    TestClientResponseHeaders("0") // status line
    TestClientResponseHeaders("Server")
    util.Dump("")

    // body
    TestClientResponseBody()
    util.Dump("")

    // body parsed
    TestClientResponseData()
}

/**
 * TestClientResponseStatus
 */
func TestClientResponseStatus() {
    var response = _doRequest("GET /")
    util.Dumpf("Response Status >> %s", response.GetStatus())
}

/**
 * TestClientResponseStatusCode
 */
func TestClientResponseStatusCode() {
    var response = _doRequest("GET /")
    util.Dumpf("Response Status Code >> %d", response.GetStatusCode())
}

/**
 * TestClientResponseStatusText
 */
func TestClientResponseStatusText() {
    var response = _doRequest("GET /")
    util.Dumpf("Response Status Text >> %s", response.GetStatusText())
}

/**
 * TestClientResponseHeaders.
 */
func TestClientResponseHeaders(key string) {
    var response = _doRequest("GET /")
    if key == "" {
        util.Dumpf("Response Headers >> %+v", response.GetHeaderAll())
    } else {
        util.Dumpf("Response Headers >> %s: %+v", key, response.GetHeader(key))
    }
}

/**
 * TestClientResponseBody
 */
func TestClientResponseBody() {
    var response = _doRequest("GET /")
    util.Dumpf("Response Body >> len: %d body: %+v",
        len(response.GetBody()), response.GetBody())
}

/**
 * TestClientResponseData.
 */
func TestClientResponseData() {
    type Response struct {
        CouchDB string
        Uuid    string
        Version string
        Vendor  map[string]string
    }

    var response = _doRequest("GET /")
    data, err := response.GetBodyData(&Response{})
    if err != nil {
        util.Dumps(err)
        return
    }
    util.Dumpf("Response Body Parsed >> type: %T value: %+v", data, data)
    util.Dumpf("Response Body Parsed >> couchdb: %s", data.(*Response).CouchDB)
    util.Dumpf("Response Body Parsed >> uuid: %s", data.(*Response).Uuid)
    util.Dumpf("Response Body Parsed >> version: %s", data.(*Response).Version)
    util.Dumpf("Response Body Parsed >> vendor: %s", data.(*Response).Vendor)
    util.Dumpf("Response Body Parsed >> vendor.name: %s", data.(*Response).Vendor["name"])
    util.Dumpf("Response Body Parsed >> vendor.version: %s", data.(*Response).Vendor["version"])
}
