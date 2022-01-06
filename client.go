package main

import (
	"net/http"
	"sync"
)

type client struct {
	client *http.Client
}

var instance *client = nil
var once sync.Once

func GetNewClient() *client {
	once.Do(func() {
		instance = &client{
			client: &http.Client{},
		}
	})
	return instance
}
