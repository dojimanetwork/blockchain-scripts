package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	gsrpc "github.com/dojimanetwork/go-polka-rpc/v5"
	"github.com/dojimanetwork/go-polka-rpc/v5/rpc/author"
	"github.com/dojimanetwork/go-polka-rpc/v5/signature"
	gsrpcTypes "github.com/dojimanetwork/go-polka-rpc/v5/types"
	"github.com/rs/zerolog/log"
)

const (
	endpoint = "wss://dotws-test.h4s.dojima.network:9944"
	westend  = "wss://westend-rpc.polkadot.io"
)

func main() {
	log := log.Logger.With().Str("module", "polkadot").Logger()
	api, err := gsrpc.NewSubstrateAPI(endpoint)
	opts := types.SerDeOptions{NoPalletIndices: true}
	types.SetSerDeOptions(opts)
	if err != nil {
		log.Info().Err(err).Msg("api not initilaized")
	}
	meta, err := api.RPC.State.GetMetadataLatest()

	if err != nil {
		fmt.Errorf("error %w", err)
	}

	hexaAddr := strings.Join([]string{"0x", "d2c2e63069b7422f37f5c6bb6cf4241d406eb0bb33a8333649a6b77151244c2e"}, "")
	dest, err := gsrpcTypes.NewMultiAddressFromHexAccountID(hexaAddr)
	// log.Info().Str("Hexa address", hexaAddr).Interface("dest", dest).Msg("address details")
	if err != nil {
		panic(err)
	}
	memo := []byte("memo:OUT:E250EBC0EBF271ED23C41B23D5024C65BAE5563819F7537E63605EEA86485839")
	call1, err := gsrpcTypes.NewCall(meta, "System.remark", memo)
	log.Info().Interface("remark call", call1).Msg("remark details")
	if err != nil {
		panic(err)
	}

	// keyringPair := signature.KeyringPair{Address: kp.Public().Hex(), PublicKey: kp.Public().Encode(), URI: mnemonic}

	if err != nil {
		panic(err)
	}

	// bal, ok := new(big.Int).SetString("246506515540", 10)
	amount := gsrpcTypes.NewUCompactFromUInt(346506515540)
	call2, err := gsrpcTypes.NewCall(meta, "Balances.transfer", dest, amount)
	log.Info().Interface("transfer call", call2).Msg("Transfer details")
	if err != nil {
		panic(err)
	}

	batchCall, err := gsrpcTypes.NewCall(meta, "Utility.batch_all", []gsrpcTypes.Call{call1, call2})
	if err != nil {
		panic(err)
	}
	// log.Info().Interface("batch call", batchCall).Msg("Batch details")
	genesisHash, err := api.RPC.Chain.GetBlockHash(0)

	if err != nil {
		fmt.Errorf("error %w", err)
	}

	mnemonic := "hero eagle luxury slight survey catch toy goat model general alarm inner"
	// mnemonic := "letter ethics correct bus asset pipe tourist vapor envelope kangaroo warm dawn"
	kp, err := signature.KeyringPairFromSecret(mnemonic, 42)
	var sub *author.ExtrinsicStatusSubscription

	// for {
	aliceStorageKey, err := gsrpcTypes.CreateStorageKey(meta, "System", "Account", kp.PublicKey)
	log.Info().Interface("storage key call", aliceStorageKey).Msg("Storage details")
	if err != nil {
		panic(err)
	}

	var accountInfo gsrpcTypes.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(aliceStorageKey, &accountInfo)

	if err != nil || !ok {
		panic(err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	// log.Info().Interface("runtime call", rv).Msg("runtime details")
	if err != nil {
		panic(err)
	}

	ext := gsrpcTypes.NewExtrinsic(batchCall)
	// nonce := uint32(accountInfo.Nonce)
	fmt.Println(ext)
	signOpts := gsrpcTypes.SignatureOptions{
		BlockHash:          genesisHash, // using genesis since we're using immortal era
		Era:                gsrpcTypes.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              gsrpcTypes.NewUCompactFromUInt(uint64(0)),
		SpecVersion:        rv.SpecVersion,
		Tip:                gsrpcTypes.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}
	if err := ext.Sign(kp, signOpts); err != nil {
		panic(err)
	}

	if err != nil {
		fmt.Println(err)
	}

	sub, err = api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	// log.Info().Msgf("sub %v", sub)
	if err != nil {
		panic(err)
		// continue
	}

	//	break
	// }

	defer sub.Unsubscribe()

	select {
	case <-time.After(1 * time.Minute):
		panic("Timeout reached")
	case st := <-sub.Chan():
		extStatus, _ := st.MarshalJSON()
		fmt.Println("Done with status -", string(extStatus))
	case err := <-sub.Err():
		panic(err)
	}
}
