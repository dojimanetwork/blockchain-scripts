package main

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	iScale "github.com/itering/scale.go"
	iTypes "github.com/itering/scale.go/types"
	scaleBytes "github.com/itering/scale.go/types/scaleBytes"
	iUtil "github.com/itering/subscan/util"
	iMetadata "github.com/itering/substrate-api-rpc/metadata"
	irpc "github.com/itering/substrate-api-rpc/rpc"
	iutils "github.com/itering/substrate-api-rpc/util"
	iws "github.com/itering/substrate-api-rpc/websocket"
	"strconv"
	"strings"
)

const (
	endpoint = "wss://westend-rpc.polkadot.io"
)

func ConnectDot() {
	iws.SetEndpoint(endpoint)
}

func main() {
	ConnectDot()
	api, err := gsrpc.NewSubstrateAPI(endpoint)
	//iapi, err := iws.Init()
	//isconn := iapi.Conn.IsConnected()
	//iapi.Conn.Dial(endpoint, http.Header{})
	//iapi.Conn.RemoteAddr()
	//fmt.Printf("%v%t", iapi.Conn, isconn)
	if err != nil {

	}
	//
	blockHash, err := api.RPC.Chain.GetBlockHash(12782886)
	fmt.Println("block hash", blockHash.Hex())

	if err != nil {

	}

	codedMetadataAtHash, err := irpc.GetMetadataByHash(nil, blockHash.Hex())
	if err != nil {

	}

	fmt.Printf("coded metadata hash %v", codedMetadataAtHash)

	metaDataInBytes := iutils.HexToBytes(codedMetadataAtHash)
	m := iScale.MetadataDecoder{}
	m.Init(metaDataInBytes)
	m.Process()

	iMetadata.Latest(&iMetadata.RuntimeRaw{
		Spec: 12,
		Raw:  strings.TrimPrefix(codedMetadataAtHash, "0x"),
	})
	currentMetadata := iMetadata.RuntimeMetadata[12]
	//fmt.Printf("metadata in bytes %v", currentMetadata)

	v := &irpc.JsonRpcResult{}

	err = iws.SendWsRequest(nil, v, irpc.ChainGetBlock(12782886, blockHash.Hex()))
	if err != nil {
		fmt.Println("Could not read the block", err)
	}
	fmt.Printf("block %v", v.ToBlock())
	rpcBlock := v.ToBlock()
	blockHeight, err := strconv.ParseInt(hexaNumberToInteger(rpcBlock.Block.Header.Number), 16, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("BXL: readBlockUsingItering: blockHeight: ", blockHeight)
	_, _ = decodeExtrinsics(rpcBlock.Block.Extrinsics, currentMetadata, 12)

}

// hexaNumberToInteger
func hexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}

func decodeExtrinsics(list []string, metadata *iMetadata.Instant, spec int) (r []iScale.ExtrinsicDecoder, err error) {
	defer func() {
		if fatal := recover(); fatal != nil {
			err = fmt.Errorf("Recovering from panic in DecodeExtrinsic: %v", fatal)
		}
	}()

	m := iTypes.MetadataStruct(*metadata)
	for _, extrinsicRaw := range list {
		e := iScale.ExtrinsicDecoder{}
		option := iTypes.ScaleDecoderOption{Metadata: &m, Spec: spec}
		e.Init(scaleBytes.ScaleBytes{Data: iUtil.HexToBytes(extrinsicRaw)}, &option)
		e.Process()

		r = append(r, e)
	}
	return r, nil
}
