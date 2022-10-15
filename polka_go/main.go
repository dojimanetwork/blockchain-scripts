package main

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
)

const (
	westendpoint  = "wss://westend-rpc.polkadot.io"
	localEndpoint = "ws://127.0.0.1:9944"
)

type Client struct {
	dotRpcClient *gsrpc.SubstrateAPI
}

func main() {
	dotRpcClient, err := gsrpc.NewSubstrateAPI(localEndpoint)

	cli := &Client{
		dotRpcClient: dotRpcClient,
	}

	if err != nil {
		_ = fmt.Errorf("err %v", err)
	}

	height, err := cli.GetHeight()

	if err != nil {
		_ = fmt.Errorf("err %v", err)
	}

	fmt.Println("height", height)

}

func (c *Client) GetHeight() (int64, error) {
	BlockInfo, err := c.dotRpcClient.RPC.Chain.GetBlockLatest()
	// TODO: verify conversion of uint32 dot block number to int64
	return int64(BlockInfo.Block.Header.Number), err
}
