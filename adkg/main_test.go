package main

import (
	"testing"

	"gopkg.in/h2non/baloo.v3"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New("http://localhost:8080")

var payload = `{"jsonrpc":"2.0","method":"HealthCheck","id":10,"params":{"provider":"passwordless","user_id":"yvbgr.blockchain@gmail.com","app_id":"4fe5879540c6db0ba158213a971a09cfbb7d2db4"}}`

func TestShareAssign(t *testing.T) {
	test.Post("/rpc").
		JSON([]byte(payload)).
		Expect(t).
		Status(200).
		Type("json").
		Done()
}
