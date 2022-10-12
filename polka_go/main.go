package main

import (
	"bytes"
	"fmt"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/scale"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
)

func main() {
	//client, err := dotc.Connect("wss://westend-rpc.polkadot.io")
	api, err := gsrpc.NewSubstrateAPI("wss://westend-rpc.polkadot.io")

	if err != nil {
		fmt.Errorf("error in connecting %w", err)
	}

	//rpcClient := dotRpc.NewChain(client)

	//if err != nil {
	//	fmt.Errorf("error in rpc connection %w", err)
	//}

	blockHash, err := api.RPC.Chain.GetBlockHash(12782886)

	if err != nil {
		fmt.Errorf("error in rpc connection %w", err)
	}

	blockInfo, err := api.RPC.Chain.GetBlock(blockHash)

	for _, extrinsic := range blockInfo.Block.Extrinsics {
		if extrinsic.Method.CallIndex.SectionIndex == 16 && extrinsic.Method.CallIndex.MethodIndex == 0 {
			fmt.Printf("extrinsic %v\n\n\n\n", extrinsic)
			//fmt.Printf("version %d\n", extrinsic.Type())
			sender, _ := subkey.SS58Address(extrinsic.Signature.Signer.AsID[:], uint8(42))
			fmt.Println("BXL: DOTChain pk bytes: ", sender)

			_ = scale.NewDecoder(bytes.NewReader(extrinsic.Method.Args))
			fmt.Println("BXL: ext.Method.Args: ", extrinsic.Method.Args)
		}
	}

}
