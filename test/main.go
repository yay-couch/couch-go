package main

import (
    // test_uuid "./uuid"
    // test_client "./client"
    // test_server "./server"
    // test_database "./database"
    test_document "./document"
)

func main() {
    /* client */
    // test_client.TestClientResponseStatus()
    // ...

    /* server */
    // test_server.TestPing()
    // ...

    /* database */
    // test_database.TestPing()
    // ...

    /* document */
    // test_document.TestPing()
    // test_document.TestIsExists()
    // test_document.TestIsNotModified()
    // test_document.TestFind()
    // test_document.TestFindRevisions()
    test_document.TestFindRevisionsExtended()

    /* uuid */
    // test_uuid.TestAll()
}
