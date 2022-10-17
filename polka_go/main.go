package main

import (
	"encoding/json"
	"errors"
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	irpc "github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
	iws "github.com/itering/substrate-api-rpc/websocket"
	"net/http"
)

const (
	westendpoint  = "wss://westend-rpc.polkadot.io"
	localEndpoint = "ws://localhost:9944"
)

type Client struct {
	latestBlock  *irpc.JsonRpcResult
	dotRpcClient *gsrpc.SubstrateAPI
}

//system_syncState
type DotSyncState struct {
	StartingBlock int `json:"startingBlock"`
	CurrentBlock  int `json:"currentBlock"`
	HighestBlock  int `json:"highestBlock"`
}

func main() {
	dotRpcClient, err := gsrpc.NewSubstrateAPI(localEndpoint)
	//fmt.Println(dotRpcClient)
	websocket.SetEndpoint(localEndpoint)
	iapi, err := iws.Init()
	iapi.Conn.Dial(localEndpoint, http.Header{})

	cli := &Client{
		dotRpcClient: dotRpcClient,
	}

	fmt.Println(iapi, cli.dotRpcClient)
	if err != nil {
		_ = fmt.Errorf("err %v", err)
	}

	height, err := cli.GetHeight()

	if err != nil {
		_ = fmt.Errorf("err %v", err)
	}

	fmt.Println("height", height)

	blockResult := &irpc.JsonRpcResult{}

	query, err := json.Marshal(irpc.Param{Id: 2, Params: []string{}, JsonRpc: "2.0", Method: "system_syncState"})
	if err != nil {
		_ = fmt.Errorf("query error %w", err)
	}
	err = iws.SendWsRequest(nil, blockResult, query)
	if err != nil {

	}
	cli.latestBlock = blockResult

	fmt.Println(cli.ToSyncState().CurrentBlock)

}

func (c *Client) GetHeight() (int64, error) {
	BlockInfo, err := c.dotRpcClient.RPC.Chain.GetBlockLatest()
	// TODO: verify conversion of uint32 dot block number to int64
	return int64(BlockInfo.Block.Header.Number), err
}

func (c *Client) ToSyncState() *DotSyncState {

	if c.checkErr() != nil {
		return nil
	}

	if (c.latestBlock).Result == nil {
		return nil
	}

	result := (c.latestBlock).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil
	}
	v := &DotSyncState{}
	marshal, err := json.Marshal(result)
	fmt.Println(result)
	if err != nil {
		return nil
	}
	_ = json.Unmarshal(marshal, v)

	return v
}

func (p *Client) checkErr() error {
	if p.latestBlock.Error != nil {
		return errors.New(p.latestBlock.Error.Message)
	}
	return nil
}
