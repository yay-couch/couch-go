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
    return client.DoRequest(uri, nil, "", nil)
}

/**
 * Test client.
 */
func TestClient() {
    couch  := _couch.New(nil)
    client := _couch.NewClient(couch, nil)
    // or
    // client := _couch.NewClient(couch, map[string]interface{}{
    //     "Host": "127.0.0.1",
    // })
    var response = client.DoRequest("GET /", nil, "", nil)
    _dumpf("response >> %+v", response)
}

// func TestClientResponse() {
//     var response = _doRequest("GET /")
//     var responseBody = response.GetBody()
//     _dumpf("response: len(%d) body(%+v)", len(responseBody), responseBody)
// }
// func TestClientResponseStatusCode() {
//     var response = _doRequest("GET /")
//     var responseBody = response.GetBody()
//     _dumpf("response: len(%d) body(%+v)", len(responseBody), responseBody)
// }
// func TestClientResponseStatusText() {
//     var response = _doRequest("GET /")
//     var responseBody = response.GetBody()
//     _dumpf("response: len(%d) body(%+v)", len(responseBody), responseBody)
// }

// func TestClientResponseBodyParse() {
//     var response = _doRequest("GET /")
//     type Response struct {
//         CouchDB string
//         Uuid    string
//         Version string
//         Vendor  map[string]string
//     }
//     body, err := u.ParseBody(response.GetBody(), &Response{})
//     if err != nil {
//         _dumps(err)
//         return
//     }
//     _dumps(body)
//     _dumps(body.(*Response).CouchDB)
//     _dumps(body.(*Response).Uuid)
//     _dumps(body.(*Response).Version)
//     _dumps(body.(*Response).Vendor)
//     _dumps(body.(*Response).Vendor["name"])
//     _dumps(body.(*Response).Vendor["version"])
// }
