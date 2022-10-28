package main

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/author"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	gsrpcTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"time"
)

const (
	endpoint = "ws://localhost:9944"
	westend  = "wss://westend-rpc.polkadot.io"
)

func main() {

	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		fmt.Errorf("error %w", err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()

	if err != nil {
		fmt.Errorf("error %w", err)
	}

	dest, err := gsrpcTypes.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")

	if err != nil {
		panic(err)
	}

	//memo := gsrpcTypes.NewData([]byte("memo:ADD:DOT.DOT:dojima1nh4y3gqxsn7ymm9t45zwsz3h8p9tm7pev8my62"))
	////memoBytes, err := codec.Encode(memo)
	//if err != nil {
	//	panic(err)
	//}

	call1, err := gsrpcTypes.NewCall(meta, "System.remark", []byte("memo:ADD:DOT.DOT:dojima1nh4y3gqxsn7ymm9t45zwsz3h8p9tm7pev8my62"))
	if err != nil {
		panic(err)
	}

	//keyringPair := signature.KeyringPair{Address: kp.Public().Hex(), PublicKey: kp.Public().Encode(), URI: mnemonic}

	if err != nil {
		panic(err)
	}

	call2, err := gsrpcTypes.NewCall(meta, "Balances.transfer", dest, gsrpcTypes.NewUCompactFromUInt(1000000))
	if err != nil {
		panic(err)
	}

	batchCall, err := gsrpcTypes.NewCall(meta, "Utility.batch_all", []gsrpcTypes.Call{call1, call2})
	if err != nil {
		panic(err)
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		fmt.Errorf("error %w", err)
	}

	//mnemonic := "entire material egg meadow latin bargain dutch coral blood melt acoustic thought"
	mnemonic := "letter ethics correct bus asset pipe tourist vapor envelope kangaroo warm dawn"
	kp, err := signature.KeyringPairFromSecret(mnemonic, 42)
	var sub *author.ExtrinsicStatusSubscription

	//for {
	aliceStorageKey, err := gsrpcTypes.CreateStorageKey(meta, "System", "Account", kp.PublicKey)

	if err != nil {
		panic(err)
	}

	var accountInfo gsrpcTypes.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(aliceStorageKey, &accountInfo)

	if err != nil || !ok {
		panic(err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	ext := gsrpcTypes.NewExtrinsic(batchCall)
	nonce := uint32(accountInfo.Nonce)

	signOpts := gsrpcTypes.SignatureOptions{
		BlockHash:          genesisHash, // using genesis since we're using immortal era
		Era:                gsrpcTypes.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              gsrpcTypes.NewUCompactFromUInt(uint64(nonce)),
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

	if err != nil {
		fmt.Printf("extrinsic submit failde %v", err)
		//continue
	}

	//	break
	//}

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
