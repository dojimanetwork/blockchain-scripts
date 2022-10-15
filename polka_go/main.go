package main

import (
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
	dotRpcClient *gsrpc.SubstrateAPI
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

	err = iws.SendWsRequest(nil, blockResult, irpc.ChainGetBlockHash(1, 1))
	if err != nil {

	}
	fmt.Println(blockResult)
	blockHash := blockResult.Result
	fmt.Println(blockHash)

}

func (c *Client) GetHeight() (int64, error) {
	BlockInfo, err := c.dotRpcClient.RPC.Chain.GetBlockLatest()
	// TODO: verify conversion of uint32 dot block number to int64
	return int64(BlockInfo.Block.Header.Number), err
}
