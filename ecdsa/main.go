package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os/exec"
	"strings"

	gsrpc "github.com/dojimanetwork/go-polka-rpc/v5"
	"github.com/dojimanetwork/go-polka-rpc/v5/config"
	"github.com/dojimanetwork/go-polka-rpc/v5/signature"
	"github.com/dojimanetwork/go-polka-rpc/v5/types"
	"github.com/dojimanetwork/go-polka-rpc/v5/types/codec"
	"github.com/dojimanetwork/hermes/narada/chainclients/polkadot"
	"golang.org/x/crypto/blake2b"
)

type EcdsaKeyringPair struct {
	// URI is the derivation path for the private key in subkey
	URI string
	// Address is an SS58 address
	Address string
	// PublicKey
	PublicKey string
	// Account Id
	AccountID []byte
}

func main() {
	api := NewSubstrateAPI()
	SetSerDeOptions()
	metadata := GetMetadataLatest(api)
	genesisHash := GetGenesisHash(api)
	runtimeVersion := GetRuntimeVersionLatest(api)

	mnemonic := "bicycle youth kidney ugly actual web thank rate good garage program lend"
	aliceEd25519KeyringPair, accountID, err := polkadot.Ecdsa_KPFromSeedPhrase(mnemonic, 42)
	// kp, err := signature.KeyringPairFromSecret(mnemonic, 42)
	// hexAccountId := subkey.EncodeHex(accountID)
	if err != nil {
		panic(fmt.Sprintf("failed to get kp %v", err))
	}

	storageKey := CreateStorageKey(metadata, accountID)
	accountInfo := GetStorageLatest(api, storageKey)
	nonce := accountInfo.Nonce

	bobSr25519 := NewAddressFromHexAccountID("0x7a99a7227cd7ddf60976d0c2725b627b35968b51c35e2dc3572c9464e91c1b2b")
	options := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        runtimeVersion.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: runtimeVersion.TransactionVersion,
	}
	call := CreateBalanceCall(metadata, bobSr25519, 1)
	memo := []byte("memo:OUT:E250EBC0EBF271ED23C41B23D5024C65BAE5563819F7537E63605EEA86485839")
	call1, err := types.NewCall(metadata, "System.remark", memo)

	if err != nil {
		panic(err)
	}

	batchCall, err := types.NewCall(metadata, "Utility.batch_all", []types.Call{call1, call})
	if err != nil {
		panic(err)
	}

	ext := types.NewExtrinsic(batchCall)

	edExt := polkadot.EcdsaExtrinsic{
		Extrinsic: ext,
	}
	// sign using Ed25519
	err = edExt.EcdsaSign(aliceEd25519KeyringPair, options, accountID)

	if err != nil {
		panic(err)
	}

	hash, err := api.RPC.Author.SubmitExtrinsic(edExt.Extrinsic)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Printf("Transfer sent with extrinsic hash %#x\n", hash)

}

func SignUsingEd25519(e types.Extrinsic, signer signature.KeyringPair, o types.SignatureOptions) (types.Extrinsic, error) {
	if e.Type() != types.ExtrinsicVersion4 {
		return e, fmt.Errorf("unsupported extrinsic version: %v (isSigned: %v, type: %v)", e.Version, e.IsSigned(), e.Type())
	}

	mb, err := codec.Encode(e.Method)
	if err != nil {
		return e, err
	}

	era := o.Era
	if !o.Era.IsMortalEra {
		era = types.ExtrinsicEra{IsImmortalEra: true}
	}

	payload := types.ExtrinsicPayloadV4{
		ExtrinsicPayloadV3: types.ExtrinsicPayloadV3{
			Method:      mb,
			Era:         era,
			Nonce:       o.Nonce,
			Tip:         o.Tip,
			SpecVersion: o.SpecVersion,
			GenesisHash: o.GenesisHash,
			BlockHash:   o.BlockHash,
		},
		TransactionVersion: o.TransactionVersion,
	}

	signerPubKey, err := types.NewMultiAddressFromAccountID(signer.PublicKey)
	data, err := codec.Encode(payload)
	if err != nil {
		return e, err
	}

	sig, err := SignEd25519(data, signer, "Ed25519")

	if err != nil {
		return e, err
	}

	multiSig := types.MultiSignature{IsEd25519: true, AsEd25519: sig}

	// multiSig := types.MultiSignature{IsEd25519: true, AsEd25519: sig}
	// You would use this if you are using Ecdsa since it needs to return bytes

	extSig := types.ExtrinsicSignatureV4{
		Signer:    signerPubKey,
		Signature: multiSig,
		Era:       era,
		Nonce:     o.Nonce,
		Tip:       o.Tip,
	}

	e.Signature = extSig

	// mark the extrinsic as signed
	e.Version |= types.ExtrinsicBitSigned

	return e, nil

}

func SignEd25519(data []byte, signer signature.KeyringPair, scheme string) (types.Signature, error) {
	// if data is longer than 256 bytes, hash it first
	if len(data) > 256 {
		h := blake2b.Sum256(data)
		data = h[:]
	}

	// data to stdin
	dataHex := hex.EncodeToString(data)

	// use "subkey" command for signature
	cmd := exec.Command("/Users/luffybhaagi/dojima/substrate/target/release/subkey", "sign", "--hex", "--message", dataHex, "--scheme", scheme, "--suri", "0x8a81cd73276e0e9fd43472950e3df0eaf38411abbe2838a995db435588fa8737")
	// cmd := exec.Command("subkey", "sign", "--hex", "--suri", privateKeyURI)

	cmd.Stdin = strings.NewReader(dataHex)
	log.Printf("%v sign  --hex --message %v --scheme %v --suri %v ", "/Users/luffybhaagi/dojima/substrate/target/release/subkey", dataHex, scheme, "0x8a81cd73276e0e9fd43472950e3df0eaf38411abbe2838a995db435588fa8737")

	// execute the command, get the output
	out, err := cmd.Output()
	if err != nil {
		return types.Signature{}, fmt.Errorf("failed to sign with subkey: %v", err.Error())
	}

	// remove line feed
	if len(out) > 0 && out[len(out)-1] == 10 {
		out = out[:len(out)-1]
	}

	outStr := fmt.Sprintf("%x", out)
	output, err := hex.DecodeString(outStr)

	if err != nil {
		return types.Signature{}, err
	}

	// hxpubkey := hex.EncodeToString(signer.PublicKey)
	log.Printf("/Users/luffybhaagi/dojima/substrate/target/release/subkey verify --hex --message %v --scheme %v %v %v ", dataHex, scheme, string(out), signer.Address)

	// Return a new Signature
	return types.NewSignature(output), err

}

func NewSubstrateAPI() *gsrpc.SubstrateAPI {
	// Instantiate the API
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}
	return api
}

func SetSerDeOptions() types.SerDeOptions {
	opts := types.SerDeOptions{NoPalletIndices: true}
	types.SetSerDeOptions(opts)
	return opts
}

func GetMetadataLatest(api *gsrpc.SubstrateAPI) *types.Metadata {
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}
	return meta
}

func GetGenesisHash(api *gsrpc.SubstrateAPI) types.Hash {
	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}
	return genesisHash
}

func GetRuntimeVersionLatest(api *gsrpc.SubstrateAPI) *types.RuntimeVersion {
	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}
	return rv
}

func GetAliceEd25519KeyringPair() signature.KeyringPair {
	// Secret phrase:       bicycle youth kidney ugly actual web thank rate good garage program lend
	// Network ID:        substrate
	// Secret seed:       0x614d27ba8970890a65733da33769f17f118d58b615f1d354a8ccab322d76292e
	// Public key (hex):  0x033e59b502f20d752dd4e1b462975dbdf6f0c19a672d6cee21e21cb9aa4707b028
	// Account ID:        0xef076239c21b8b6a5d94a39686f593be68a5c4b3e33e58e8597ac1b876af7afd
	// Public key (SS58): KWA84cn89vShDfpAZjcgktDT8xRqSTBaoM5W5uqrdW3VhWmnT
	// SS58 Address:      5HU7VfVzwN3APaEACDW6J7BArdDqY9MwhAfRQ6t1zXfiQ36G

	publicKey, err := codec.HexDecodeString("0x68da43131f395a1e916535011ccc3d34df3fe4347e3f093720343a57002ba0c5")
	if err != nil {
		panic(err)
	}

	keypair := signature.KeyringPair{
		URI:       "//bhagath",
		Address:   "5ESBg5NXbVauTLPsnXuuDLPT2V2NiLoxXXhLnYU5D2xMnDtm",
		PublicKey: publicKey,
	}
	return keypair
}

func CreateStorageKey(meta *types.Metadata, pubkey []byte) types.StorageKey {
	key, err := types.CreateStorageKey(meta, "System", "Account", pubkey)
	if err != nil {
		panic(err)
	}
	return key
}

func GetStorageLatest(api *gsrpc.SubstrateAPI, key types.StorageKey) types.AccountInfo {
	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(fmt.Sprintf("failed to get account info %v", err))
	}
	return accountInfo
}

func NewAddressFromHexAccountID(hexAccountId string) types.MultiAddress {
	addr, err := types.NewMultiAddressFromHexAccountID(hexAccountId)
	if err != nil {
		panic(err)
	}
	return addr
}

func CreateBalanceCall(meta *types.Metadata, toAddr types.MultiAddress, amount uint64) types.Call {
	call, err := types.NewCall(meta, "Balances.transfer", toAddr, types.NewUCompactFromUInt(10000000000000))
	if err != nil {
		panic(err)
	}
	return call
}
