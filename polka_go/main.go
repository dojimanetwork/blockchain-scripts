package main

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	iScale "github.com/itering/scale.go"
	iTypes "github.com/itering/scale.go/types"
	scaleBytes "github.com/itering/scale.go/types/scaleBytes"
	iUtil "github.com/itering/subscan/util"
	iSS58 "github.com/itering/subscan/util/ss58"
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

	//fmt.Printf("coded metadata hash %v", codedMetadataAtHash)

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
	//fmt.Printf("block %v", v.ToBlock())
	rpcBlock := v.ToBlock()
	blockHeight, err := strconv.ParseInt(hexaNumberToInteger(rpcBlock.Block.Header.Number), 16, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("BXL: readBlockUsingItering: blockHeight: ", blockHeight)
	decodedExtrinsics, _ := decodeExtrinsics(rpcBlock.Block.Extrinsics, currentMetadata, 12)
	//fmt.Printf("decoded extrinsics %v", decodedExtrinsics[2])
	for _, e := range decodedExtrinsics {
		ParseUtilityBatch(&e, blockHeight)

	}
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

func ParseUtilityBatch(e *iScale.ExtrinsicDecoder, blockHeight int64) {

	callDecoder := iTypes.Call{ScaleDecoder: e.ScaleDecoder}
	callDecoder.Process()
	module, _ := e.Value.(map[string]interface{})
	if module["call_module"] == "Utility" && module["call_module_function"] == "batch_all" {
		fromAddressStr := fmt.Sprintf("%v", e.Address)
		fromAddress := iSS58.Encode(fromAddressStr, iUtil.StringToInt("42"))
		fmt.Println(fromAddress)

	}

	//{
	//{
	//{[
	//109 2 132 0 210 194 230 48 105 183 66 47 55 245 198 187 108 244 36 29 64 110 176 187 51 168 51 54 73 166 183 113 81 36 76 46 1 108 66 116 2 254 22 187 211 118 218 240 128 102 137 21 225 180 180 37 100 202 214 251 86 82 169 198 227 249 173 221 53 164 7 210 179 183 206 255 101 18 44 193 93 2 118 5 212 228 71 109 107 2 201 98 19 0 15 163 30 11 152 44 131 37 2 16 0 16 2 8 0 1 28 116 101 115 116 105 110 103 4 0 0 61 188 217 129 149 46 91 205 77 164 232 238 20 215 107 154 72 61 56 152 199 106 36 82 165 128 210 221 110 92 199 254 2 90 98 2
	//] 157}
	//map[
	//account_id:0xd2c2e63069b7422f37f5c6bb6cf4241d406eb0bb33a8333649a6b77151244c2e
	//address_type:Id
	//call_code:1002
	//call_module:Utility
	//call_module_function:batch_all
	//era:2502
	//extrinsic_hash:a9a708657612ce331f0ced07f648d46d810d1eb920aece0ba462226c56f2d1eb
	//extrinsic_length:155
	//nonce:4
	//params:[
	//{
	//calls Vec<Call> Vec<<T as Config>::Call> [
	//map[
	//call_index:0001
	//call_module:System
	//call_name:remark
	//params:[
	//{remark Vec<U8> testing}]]
	//map[
	//call_index:0400
	//call_module:Balances
	//call_name:transfer
	//params:[{
	//dest sp_runtime:multiaddress:MultiAddress
	//map[
	//Id:0x3dbcd981952e5bcd4da4e8ee14d76b9a483d3898c76a2452a580d2dd6e5cc7fe]} {
	//value compact<U128> 10000000}]]]}]
	//signature:0x6c427402fe16bbd376daf080668915e1b4b42564cad6fb5652a9c6e3f9addd35a407d2b3b7ceff65122cc15d027605d4e4476d6b02c96213000fa31e0b982c83
	//tip:0 version_info:84]
	//0800011c74657374696e670400003dbcd981952e5bcd4da4e8ee14d76b9a483d3898c76a2452a580d2dd6e5cc7fe025a6202
	//<nil> 0x14000036640 12 Utility map[]  map[]}
	//155
	//a9a708657612ce331f0ced07f648d46d810d1eb920aece0ba462226c56f2d1eb
	//84 true
	//0xd2c2e63069b7422f37f5c6bb6cf4241d406eb0bb33a8333649a6b77151244c2e
	//0x6c427402fe16bbd376daf080668915e1b4b42564cad6fb5652a9c6e3f9addd35a407d2b3b7ceff65122cc15d027605d4e4476d6b02c96213000fa31e0b982c83
	//4 2502 1002
	//[{calls Vec<Call>
	//Vec<<T as Config>::Call> [map[
	//call_index:0001
	//call_module:System call_name:remark
	//params:[{remark Vec<U8> testing}]] map[
	//call_index:0400 call_module:Balances
	//call_name:transfer params:[{
	//dest sp_runtime:multiaddress:MultiAddress map[
	//Id:0x3dbcd981952e5bcd4da4e8ee14d76b9a483d3898c76a2452a580d2dd6e5cc7fe]}
	//{value compact<U128> 10000000}]]]}] 0x14000036640 []}

	//if e.Module
}
