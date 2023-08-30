package main

import (
	"gopkg.in/h2non/baloo.v3"
)

func main() {
	// mr := jsonrpc.NewMethodRepository()
	//
	// rec := httptest.NewRecorder()
	// r, rerr := http.NewRequestWithContext(context.Background(), "", "http://localhost:8080/rpc/", bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":"test","method":"HealthCheck","params":{}}`)))
	//
	// if rerr != nil {
	//
	// }
	// r.Header.Set("Content-Type", "application/json")
	// mr.ServeHTTP(rec, r)
	//
	// res := jsonrpc.Response{}
	// err := json.NewDecoder(rec.Body).Decode(&res)
	// if err != nil {
	//
	// }

	// test stores the HTTP testing client preconfigured
	var test = baloo.New("http://localhost:8080")

	var payload = `{"jsonrpc": "2.0","method": "HealthCheck","id": 6,"params": {}}`
	test.Post("/rpc").
		JSON([]byte(payload))
}
