package main

import (
	"bytes"
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/vedhavyas/go-subkey"
	"math/big"
	"strconv"
	"strings"
)

const (
	westend = "wss://westend-rpc.polkadot.io"
)

func main() {
	api, err := gsrpc.NewSubstrateAPI(westend)
	if err != nil {
		panic(err)
	}

	metadata, err := api.RPC.State.GetMetadataLatest()
	var metadataV14 types.Metadata
	encoded, err := codec.EncodeToHex(metadata)

	if err != nil {
		panic(err)
	}

	if err := codec.DecodeFromHex(encoded, &metadataV14); err != nil {
		panic(err)
	}

	types.SetSerDeOptions(types.SerDeOptions{NoPalletIndices: true})
	storageKey, err := types.CreateStorageKey(metadata, "System", "Events", nil, nil)

	if err != nil {
		panic(err)
	}

	blockHash, err := api.RPC.Chain.GetBlockHash(5482576)
	if err != nil {
		panic(err)
	}

	raw, err := api.RPC.State.GetStorageRaw(storageKey, blockHash)
	if err != nil {
		panic(err)
	}

	events := types.EventRecords{}
	err = types.EventRecordsRaw(*raw).DecodeEventRecords(metadata, &events)
	if err != nil {
		panic(err)
	}

	// Get the block
	block, err := api.RPC.Chain.GetBlock(blockHash)
	if err != nil {
		panic(err)
	}
	// Loop through successful utility batch events
	for _, event := range events.Utility_BatchCompleted {

		ext := block.Block.Extrinsics[int(event.Phase.AsApplyExtrinsic)]
		_, err := codec.EncodeToHex(ext)
		if err != nil {
			panic(err)
		}

		resInter := FeeDetails{}
		// var resInter interface{}

		err = api.Client.Call(&resInter, "payment_queryFeeDetails", ext, blockHash.Hex())
		if err != nil {
			panic(err)
		}
		resInterDispatchInfo := DispatchInfo{}
		err = api.Client.Call(&resInterDispatchInfo, "payment_queryInfo", ext, blockHash.Hex())
		if err != nil {
			panic(err)
		}
		fmt.Println("BXL:  payment_queryInfo PartialFee: ", resInterDispatchInfo.PartialFee)
		// Get the feeRates to post network fee
		adjustedWeightFee := new(big.Int)
		adjustedWeightFeeStrConv, err := strconv.ParseInt(hexaNumberToInteger(resInter.InclusionFee.AdjustedWeightFee), 16, 64)
		if err != nil {
			panic(err)
		}
		adjustedWeightFee = adjustedWeightFee.SetInt64(adjustedWeightFeeStrConv)
		//fmt.Println("BXL:  adjustedWeightFee: ", adjustedWeightFee)

		baseFee := new(big.Int)
		baseFeeStrConv, err := strconv.ParseInt(hexaNumberToInteger(resInter.InclusionFee.BaseFee), 16, 64)
		if err != nil {
			panic(err)
		}
		baseFee = baseFee.SetInt64(baseFeeStrConv)
		//fmt.Println("BXL:  baseFee: ", baseFee)

		lenFee := new(big.Int)
		lenFeeStrConv, err := strconv.ParseInt(hexaNumberToInteger(resInter.InclusionFee.LenFee), 16, 64)
		if err != nil {
			panic(err)
		}
		lenFee = lenFee.SetInt64(lenFeeStrConv)
		//fmt.Println("BXL:  lenFee: ", lenFee)

		partialFeeCal := new(big.Int)
		partialFeeCal = partialFeeCal.Add(baseFee, lenFee)
		partialFeeCal = partialFeeCal.Add(partialFeeCal, adjustedWeightFee)
		//fmt.Println("BXL:  partialFeeCal: ", partialFeeCal)

		feeRateCal := new(big.Int)

		// Polkadot Fees
		// 1 calculate amount of DOT consumed as fees in that block (PartialFee = lenFee + baseFee + adjustedWeight)
		// 2 calculated block size (lenFee)
		// 3 feeRate is (1)/(2) (baseFee + adjustedWeight)
		feeRateCal = feeRateCal.Add(baseFee, adjustedWeightFee)
		fmt.Println("BXL:  feeRateCal: ", feeRateCal)
		decoder := scale.NewDecoder(bytes.NewReader(ext.Method.Args))
		//accountID := ext.Signature.Signer.AsID
		sender, _ := subkey.SS58Address(ext.Signature.Signer.AsID[:], uint8(42))
		fmt.Println("BXL: sender: ", sender, "")
		ncalls, err := decoder.DecodeUintCompact()
		if err != nil {
			panic(err)
		}
		for call := uint64(0); call < ncalls.Uint64(); call++ {
			callIndex := types.CallIndex{}
			err = decoder.Decode(&callIndex)
			if err != nil {
				fmt.Println("decoder Error", " = ", err)

				// return err
			}
			//

			if metadata.Version == 14 {
				for _, mod := range metadata.AsMetadataV14.Pallets {
					if uint8(mod.Index) == callIndex.SectionIndex {
						// if mod.Name == "Staking" {
						// 	fmt.Println("This should be System  ", mod.Name)
						// 	modSystem := metadata.AsMetadataV12.Modules[0]
						// 	return modSystem.Calls[1]
						// }
						callType := mod.Calls.Type.Int64()

						if typ, ok := metadata.AsMetadataV14.EfficientLookup[callType]; ok {
							if len(typ.Def.Variant.Variants) > 0 {
								for _, _ = range typ.Def.Variant.Variants {

								}
							}
						}
					}
				}
			}
			//for _, callArg := range callFunction.Args {
			//	if callArg.Type == "<T::Lookup as StaticLookup>::Source" {
			//argValue1 := types.AccountID{}
			//_ = decoder.Decode(&argValue1)
			//// https://github.com/paritytech/substrate/blob/master/ss58-registry.json
			//ss58, _ := subkey.SS58Address(argValue1[:], uint8(42))
			//fmt.Println(" dest = ", ss58)
			//
			////} else if callArg.Type == "Compact<T::Balance>" {
			//argValue2 := types.U128{}
			//_ = decoder.Decode(&argValue2)
			//fmt.Println(" = ", argValue2)
			//} else if callArg.Type == "Vec<u8>" {
			//	var argValue = types.Bytes{}
			//	// hex.DecodeString(a.Value.(string))
			//	_ = decoder.Decode(&argValue)
			//	value := string(argValue)
			//	fmt.Println("BXL: FetchTxs: Vec<u8> ", callArg.Name, "memo =", value)
			//
			//} else {
			//	var argValue = types.Bytes{}
			//	_ = decoder.Decode(&argValue)
			//	fmt.Println("BXL: FetchTxs: UNKNOWN argValue", callArg.Name, "=", argValue)
			//
			//}
			//}
		}

	}

}

func findModule(metadata *types.Metadata, index types.CallIndex) types.FunctionMetadataV4 {
	//fmt.Println(metadata.AsMetadataV14.Extrinsic)
	if metadata.Version == 14 {
		for _, mod := range metadata.AsMetadataV14.Pallets {
			if uint8(mod.Index) == index.SectionIndex {
				// if mod.Name == "Staking" {
				// 	fmt.Println("This should be System  ", mod.Name)
				// 	modSystem := metadata.AsMetadataV12.Modules[0]
				// 	return modSystem.Calls[1]
				// }
				callType := mod.Calls.Type.Int64()

				if typ, ok := metadata.AsMetadataV14.EfficientLookup[callType]; ok {
					if len(typ.Def.Variant.Variants) > 0 {
						for _, vars := range typ.Def.Variant.Variants {
							fmt.Println("Find module  ", mod.Name, vars.Name)
						}
					}
				}

				return types.FunctionMetadataV4{}
			}
		}
	}

	panic("Unknown call")
}

// hexaNumberToInteger
func hexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}

type FeeDetails struct {
	InclusionFee InclusionFee
}

type InclusionFee struct {
	AdjustedWeightFee string `json:"adjustedWeightFee"`
	BaseFee           string `json:"baseFee"`
	LenFee            string `json:"lenFee"`
}

// DispatchInfo
type DispatchInfo struct {
	// Weight of this transaction
	Weight interface{} `json:"weight"`
	// Class of this transaction
	Class string `json:"class"`
	// PaysFee indicates whether this transaction pays fees
	PartialFee string `json:"partialFee"`
}
