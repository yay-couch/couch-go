package server

import (
    "./../../src/couch"
    "./../../src/couch/util"
)

var (
    DEBUG  = true
    DBNAME = "foo"
)

var (
    Couch    *couch.Couch
    Client   *couch.Client
    Database *couch.Database
)

func init() {
    Couch    = couch.New(nil, DEBUG)
    Client   = couch.NewClient(Couch)
    Database = couch.NewDatabase(Client, DBNAME);
}

func _document(data map[string]interface{}) *couch.Document {
    return couch.NewDocument(Database, data)
}

/**
 * TestAll
 */
func TestAll() {}


