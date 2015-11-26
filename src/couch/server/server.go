package server

import (
    // "fmt"
)

import _client   "./../client"
// import _response "./../http/response"


type Server struct {
    Client *_client.Client
}

func Shutup() {}

func New(client *_client.Client) *Server {
    return &Server{
        Client: client,
    }
}


