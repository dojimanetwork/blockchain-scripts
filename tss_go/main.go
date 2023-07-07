package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/dojimanetwork/dojima-tss/common"
	"github.com/dojimanetwork/dojima-tss/conversion"
	"github.com/dojimanetwork/dojima-tss/keygen"
	"github.com/dojimanetwork/dojima-tss/keysign"
	"github.com/dojimanetwork/dojima-tss/tss"
	gsrpc "github.com/dojimanetwork/go-polka-rpc/v5"
	"github.com/dojimanetwork/go-polka-rpc/v5/rpc/author"
	"github.com/dojimanetwork/go-polka-rpc/v5/signature"
	gsrpcTypes "github.com/dojimanetwork/go-polka-rpc/v5/types"
	"github.com/dojimanetwork/go-polka-rpc/v5/types/codec"
	common2 "github.com/dojimanetwork/hermes/common"
	"github.com/dojimanetwork/hermes/common/cosmos"
	"github.com/dojimanetwork/hermes/narada/chainclients/polkadot"
	btsskeygen "github.com/dojimanetwork/tss-lib/ecdsa/keygen"
	btss "github.com/dojimanetwork/tss-lib/tss"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/btcd/btcec"
)

const (
	partyNum         = 4
	testFileLocation = "../test_data"
	preParamTestFile = "preParam_test.data"
	endpoint         = "ws://localhost:9944"
	// as bytes.
	fieldIntSize = 32
)

// PrivateKeyLength is the fixed Private Key Length
const PrivateKeyLength = 32

// SignatureLength is the fixed Signature Length
const SignatureLength = 64

// MessageLength is the fixed Message Length
const MessageLength = 32

var (
	testPubKeys = []string{
		"dojpub1addwnpepqtdklw8tf3anjz7nn5fly3uvq2e67w2apn560s4smmrt9e3x52nt2q3l8kg",
		"dojpub1addwnpepqtspqyy6gk22u37ztra4hq3hdakc0w0k60sfy849mlml2vrpfr0wvt8c79u",
		"dojpub1addwnpepq2ryyje5zr09lq7gqptjwnxqsy2vcdngvwd6z7yt5yjcnyj8c8cn5ycz9su",
		"dojpub1addwnpepqfjcw5l4ay5t00c32mmlky7qrppepxzdlkcwfs2fd5u73qrwna0vzd44clt",
	}
	testPriKeyArr = []string{
		"MjQ1MDc2MmM4MjU5YjRhZjhhNmFjMmI0ZDBkNzBkOGE1ZTBmNDQ5NGI4NzM4OTYyM2E3MmI0OWMzNmE1ODZhNw==",
		"YmNiMzA2ODU1NWNjMzk3NDE1OWMwMTM3MDU0NTNjN2YwMzYzZmVhZDE5NmU3NzRhOTMwOWIxN2QyZTQ0MzdkNg==",
		"ZThiMDAxOTk2MDc4ODk3YWE0YThlMjdkMWY0NjA1MTAwZDgyNDkyYzdhNmMwZWQ3MDBhMWIyMjNmNGMzYjVhYg==",
		"ZTc2ZjI5OTIwOGVlMDk2N2M3Yzc1MjYyODQ0OGUyMjE3NGJiOGRmNGQyZmVmODg0NzQwNmUzYTk1YmQyODlmNA==",
	}

	testPrivBytes = [][]byte{
		{225, 20, 240, 100, 190, 1, 122, 165, 153, 178, 233, 16, 240, 68, 38, 152, 19, 226, 10, 255, 63, 86, 172, 18, 11, 113, 86, 25, 159, 231, 191, 216},
		{40, 138, 15, 186, 23, 9, 92, 51, 145, 206, 27, 159, 50, 92, 114, 196, 111, 158, 55, 212, 232, 95, 168, 195, 233, 143, 215, 109, 190, 162, 220, 50},
		{26, 144, 7, 123, 28, 35, 134, 239, 233, 173, 1, 80, 176, 117, 194, 17, 72, 28, 226, 229, 94, 183, 20, 79, 38, 61, 160, 126, 243, 49, 222, 84},
		{244, 224, 164, 219, 145, 75, 168, 52, 66, 231, 78, 69, 230, 88, 56, 52, 241, 47, 209, 193, 229, 48, 119, 198, 142, 77, 41, 40, 8, 217, 41, 152},
	}
)

type FourNodeTestSuite struct {
	servers       []*tss.TssServer
	ports         []int
	preParams     []*btsskeygen.LocalPreParams
	bootstrapPeer string
	algo          string
	polkaApi      gsrpc.SubstrateAPI
	ed25519PubKey secp256k1.PublicKey
	poolPubkey    string
	sigOpts       gsrpcTypes.SignatureOptions
	extrinsic     gsrpcTypes.Extrinsic
}

type EcdsaExtrinsic struct {
	Extrinsic gsrpcTypes.Extrinsic
}

func main() {
	common.InitLog("info", true, "four_nodes_test")
	conversion.SetupBech32Prefix()

	api, err := gsrpc.NewSubstrateAPI(endpoint)

	if err != nil {
		panic(err)
	}

	s := FourNodeTestSuite{
		ports: []int{
			17666, 17667, 17668, 17669,
		},
		bootstrapPeer: "/ip4/0.0.0.0/tcp/17666/p2p/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp",
		preParams:     getPreparams(),
		servers:       make([]*tss.TssServer, partyNum),
		polkaApi:      *api,
	}

	conf := common.TssConfig{
		KeyGenTimeout:   90 * time.Second,
		KeySignTimeout:  90 * time.Second,
		PreParamTimeout: 5 * time.Second,
		EnableMonitor:   false,
	}

	var wg sync.WaitGroup
	for i := 0; i < partyNum; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if idx == 0 {
				s.servers[idx] = s.getTssServer(idx, conf, "")
			} else {
				s.servers[idx] = s.getTssServer(idx, conf, s.bootstrapPeer)
			}
		}(i)

		time.Sleep(time.Second)
	}
	wg.Wait()

	for i := 0; i < partyNum; i++ {
		err := s.servers[i].Start()

		if err != nil {
			panic(err)
		}
	}

	s.algo = "eddsa"
	btss.SetCurve(edwards.Edwards())
	s.KeygenAndKeySign(true)
}

func (s *FourNodeTestSuite) KeygenAndKeySign(newJoinParty bool) {
	wg := sync.WaitGroup{}
	lock := &sync.Mutex{}
	keygenResult := make(map[int]keygen.Response)

	for i := 0; i < partyNum; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			var req keygen.Request
			localPubKeys := append([]string{}, testPubKeys...)
			if newJoinParty {
				req = keygen.NewRequest(localPubKeys, 10, "0.14.0", s.algo)
			} else {
				req = keygen.NewRequest(localPubKeys, 10, "0.13.0", s.algo)
			}
			res, err := s.servers[idx].Keygen(req)
			if err != nil {
				panic(err)
			}
			lock.Lock()
			defer lock.Unlock()
			keygenResult[idx] = res
		}(i)
	}

	wg.Wait()

	var poolPubKey string
	for _, item := range keygenResult {
		if len(poolPubKey) == 0 {
			poolPubKey = item.PubKey
		} else {
			if poolPubKey != item.PubKey {
				panic("not same pubkey")
			}
		}
	}

	pubkey, err := common2.NewPubKey(poolPubKey)

	log.Info().Interface("ecdsa pubkey", pubkey).Msg("POLKA::::")

	if err != nil {
		panic(err)
	}

	address, err := pubkey.GetAddress(common2.DOTCHAIN)
	log.Info().Interface("dot address", address).Msg("POLKA::::")

	// ed25519PubKey, err := pubkey.GetSecpk1PubK()

	if err != nil {
		panic(err)
	}

	// metadata
	meta, err := s.polkaApi.RPC.State.GetMetadataLatest()

	if err != nil {
		panic(err)
	}

	// transfer dot to tss pubkey address
	s.transferToPubAddress(address.String(), meta)

	// cosmos pukey
	// unmarshal := GetPubKeyBytes(poolPubKey)
	// get secp publickey
	// secpPub, err := common2.SecpPubkey(unmarshal)

	// get eddsa publickey
	ed25519PubKey, err := pubkey.GetEd25519PubK()
	if err != nil {
		panic(fmt.Sprintf("failed get pubkey from unmarshal:%v", err))
	}

	// get account id bytes
	// pubEncode := pubkey.SecpEncode(secpPub.ToECDSA())
	// accountId := pubkey.SecpAccountId(pubEncode)

	var accountInfo gsrpcTypes.AccountInfo

	for i := 0; i < 5; i++ {
		// create storage key
		storageKey, err := gsrpcTypes.CreateStorageKey(meta, "System", "Account", ed25519PubKey)
		if err != nil {
			fmt.Errorf("failed to get storage key")
			time.Sleep(time.Second * 3)
			continue
		}

		ok, err := s.polkaApi.RPC.State.GetStorageLatest(storageKey, &accountInfo)
		if ok {
			if accountInfo.Data.Free.Cmp(big.NewInt(0)) > 0 {
				break
			}
		}

		time.Sleep(time.Second * 3)
	}

	// balance call
	dest, err := gsrpcTypes.NewMultiAddressFromHexAccountID("0x7a99a7227cd7ddf60976d0c2725b627b35968b51c35e2dc3572c9464e91c1b2b")
	if err != nil {
		panic(err)
	}

	amount := gsrpcTypes.NewUCompactFromUInt(1000000000000)
	call2, err := gsrpcTypes.NewCall(meta, "Balances.transfer", dest, amount)
	if err != nil {
		panic(err)
	}

	// remark call
	memo := []byte("TSS:EDDSA:TEST:E250EBC0EBF271ED23C41B23D5024C65BAE5563819F7537E63605EEA86485839")
	call1, err := gsrpcTypes.NewCall(meta, "System.remark", memo)

	if err != nil {
		panic(err)
	}

	// batch call
	batchCall, err := gsrpcTypes.NewCall(meta, "Utility.batch_all", []gsrpcTypes.Call{call1, call2})
	if err != nil {
		panic(err)
	}

	// extrinsic
	extrinsic := gsrpcTypes.NewExtrinsic(batchCall)

	// genesis hash
	genesisHash, err := s.polkaApi.RPC.Chain.GetBlockHash(0)

	if err != nil {
		fmt.Errorf("error %w", err)
	}

	// runtime version
	rv, err := s.polkaApi.RPC.State.GetRuntimeVersionLatest()
	// log.Info().Interface("runtime call", rv).Msg("runtime details")
	if err != nil {
		panic(err)
	}

	nonce := accountInfo.Nonce

	signatureOpts := gsrpcTypes.SignatureOptions{
		BlockHash:          genesisHash, // using genesis since we're using immortal era
		Era:                gsrpcTypes.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              gsrpcTypes.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                gsrpcTypes.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	if extrinsic.Type() != gsrpcTypes.ExtrinsicVersion4 {
		panic(fmt.Errorf("unsupported extrinsic version: %v (isSigned: %v, type: %v)", extrinsic.Version, extrinsic.IsSigned(), extrinsic.Type()))
	}

	mb, err := codec.Encode(extrinsic.Method)
	if err != nil {
		panic(fmt.Sprintf("failed to encode method:%v", err))
	}

	era := signatureOpts.Era
	if !signatureOpts.Era.IsMortalEra {
		era = gsrpcTypes.ExtrinsicEra{IsImmortalEra: true}
	}

	payload := gsrpcTypes.ExtrinsicPayloadV4{
		ExtrinsicPayloadV3: gsrpcTypes.ExtrinsicPayloadV3{
			Method:      mb,
			Era:         era,
			Nonce:       signatureOpts.Nonce,
			Tip:         signatureOpts.Tip,
			SpecVersion: signatureOpts.SpecVersion,
			GenesisHash: signatureOpts.GenesisHash,
			BlockHash:   signatureOpts.BlockHash,
		},
		TransactionVersion: signatureOpts.TransactionVersion,
	}

	// ecdsa src pubkey
	// srcPubkey, err := codec.HexDecodeString(subkey.EncodeHex(accountId))
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to convert accountid to hex string:%v", err))
	// }
	//
	// signerPubKey, err := gsrpcTypes.NewMultiAddressFromAccountID(srcPubkey)

	eddsaSrcPubkey, err := gsrpcTypes.NewMultiAddressFromAccountID(GetPubKeyBytes(poolPubKey).Bytes())

	if err != nil {
		panic(fmt.Sprintf("failed to get signer pubkey:%v", err))
	}

	b, err := codec.Encode(payload)

	if err != nil {
		panic(fmt.Sprintf("failed to encode payload:%v", err))
	}

	if err != nil {
		panic(err)
	}

	// payload := hash([]byte("helloworld"))
	keysignResult := make(map[int]keysign.Response)
	for i := 0; i < partyNum; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			localPubKeys := append([]string{}, testPubKeys...)
			var keysignReq keysign.Request
			if newJoinParty {
				keysignReq = keysign.NewRequest(poolPubKey, []string{base64.StdEncoding.EncodeToString(b)}, 10, localPubKeys, "0.14.0", s.algo)
			} else {
				keysignReq = keysign.NewRequest(poolPubKey, []string{base64.StdEncoding.EncodeToString(b)}, 10, localPubKeys, "0.13.0", s.algo)
			}
			res, err := s.servers[idx].KeySign(keysignReq)
			if err != nil {
				panic(err)
			}
			lock.Lock()
			defer lock.Unlock()
			keysignResult[idx] = res
		}(i)
	}
	wg.Wait()

	// s.checkSignResult(keysignResult, poolPubKey)

	var eddsasignature gsrpcTypes.Signature
out:
	for i := 0; i < len(keysignResult)-1; i++ {
		currentSignatures := keysignResult[i].Signatures
		for j := 0; j <= len(currentSignatures)-1; j++ {
			sigBytes, sig, err := getEddsaSignature(currentSignatures[j].R, currentSignatures[j].S)
			if err != nil {
				panic(fmt.Sprintf("failed to get signature:%v", err))
			}

			// bRecoveryId, err := base64.StdEncoding.DecodeString(currentSignatures[j].RecoveryID)
			// if err != nil {
			// 	panic(fmt.Sprintf("failed to get recovery id:%v", err))
			// }

			buf, err := base64.StdEncoding.DecodeString(currentSignatures[j].Msg)
			// if len(sigBytes) != SignatureLength {
			// 	panic(errors.New("invalid signature length"))
			// }

			// if len(buf) != MessageLength {
			// 	panic(errors.New("invalid message length: not 32 byte hash"))
			// }

			edPubK, err := edwards.ParsePubKey(GetPubKeyBytes(poolPubKey).Bytes())
			origSig, err := base64.StdEncoding.DecodeString(currentSignatures[i].Signature)

			if err != nil {
				fmt.Errorf("inval ed25519 key with error %w", err)
			}

			val := sig[63] & 224
			cryEdVer := ed25519.Verify(edPubK.Serialize(), buf, origSig)
			verify := edwards.Verify(edPubK, buf, sigBytes.R, sigBytes.S)
			eddsasignature = gsrpcTypes.NewSignature(sig)
			break out
			// verify := gethSecp.VerifySignature(secpPub.SerializeUncompressed(), buf, sigBytes)
			if verify && cryEdVer && val != 0 {
				// add the recovery id at the end
				// result := make([]byte, 65)
				// copy(result, sigBytes)
				// result[64] = bRecoveryId[0]
				// signature = gsrpcTypes.NewEcdsaSignature(result)
				eddsasignature = gsrpcTypes.NewSignature(sig)
				break out
			}

		}
	}

	multiSig := gsrpcTypes.MultiSignature{IsEd25519: true, AsEd25519: eddsasignature}

	if err != nil {
		panic(fmt.Sprintf("failed to create new account id %v", err))
	}

	extSig := gsrpcTypes.ExtrinsicSignatureV4{
		Signer:    eddsaSrcPubkey,
		Signature: multiSig,
		Era:       era,
		Nonce:     signatureOpts.Nonce,
		Tip:       signatureOpts.Tip,
	}

	extrinsic.Signature = extSig

	extrinsic.Version |= gsrpcTypes.ExtrinsicBitSigned

	var sub *author.ExtrinsicStatusSubscription

	// hash, err := s.polkaApi.RPC.Author.SubmitExtrinsic(extrinsic)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Transfer sent with extrinsic hash %#x\n", hash)
	// if err != nil {
	// 	panic(err)
	// 	// continue
	// }
	//
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to subscribe %v", err))
	// }
	//

	sub, err = s.polkaApi.RPC.Author.SubmitAndWatchExtrinsic(extrinsic)

	if err != nil {
		panic(fmt.Sprintf("subcribe err %v", err))
	}
	defer sub.Unsubscribe()

	select {
	case <-time.After(1 * time.Minute):
		panic("Timeout reached")
	case st := <-sub.Chan():
		extStatus, _ := st.MarshalJSON()
		fmt.Println("Done with status -", string(extStatus))
		return
	case err := <-sub.Err():
		panic(err)
	}
}

func GetPubKeyBytes(pubkey string) types.PubKey {
	unmarshal, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, pubkey)
	if err != nil {
		panic(fmt.Sprintf("unmarshal pubkey %v", err))
	}

	return unmarshal
}

func getEcdsaSignature(r, s string) ([]byte, error) {
	rBytes, err := base64.StdEncoding.DecodeString(r)
	if err != nil {
		return nil, err
	}
	sBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	R := new(big.Int).SetBytes(rBytes)
	S := new(big.Int).SetBytes(sBytes)
	N := btcec.S256().N
	halfOrder := new(big.Int).Rsh(N, 1)
	// see: https://github.com/ethereum/go-ethereum/blob/f9401ae011ddf7f8d2d95020b7446c17f8d98dc1/crypto/signature_nocgo.go#L90-L93
	if S.Cmp(halfOrder) == 1 {
		S.Sub(N, S)
	}

	// Serialize signature to R || S.
	// R, S are padded to 32 bytes respectively.
	rBytes = R.Bytes()
	sBytes = S.Bytes()

	sigBytes := make([]byte, 64)
	// 0 pad the byte arrays from the left if they aren't big enough.
	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)
	return sigBytes, nil
}

func getEddsaSignature(r, s string) (*edwards.Signature, []byte, error) {
	rBytes, err := base64.StdEncoding.DecodeString(r)
	if err != nil {
		return nil, nil, err
	}
	sBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, nil, err
	}
	R32 := copyBytes(rBytes)
	S32 := copyBytes(sBytes)
	R := encodedBytesToBigInt(R32)
	S := encodedBytesToBigInt(S32)

	sig := append(R.Bytes(), S.Bytes()...)

	return &edwards.Signature{R: R, S: S}, sig, nil
}

// encodedBytesToBigInt converts a 32 byte little endian representation of
// an integer into a big, big endian integer.
func encodedBytesToBigInt(s *[32]byte) *big.Int {
	// Use a copy so we don't screw up our original
	// memory.
	sCopy := new([32]byte)
	for i := 0; i < 32; i++ {
		sCopy[i] = s[i]
	}
	reverse(sCopy)

	bi := new(big.Int).SetBytes(sCopy[:])

	return bi
}

// bigIntToEncodedBytesNoReverse converts a big integer into its corresponding
// 32 byte big endian representation.
func bigIntToEncodedBytesNoReverse(a *big.Int) *[32]byte {
	s := new([32]byte)
	if a == nil {
		return s
	}
	// Caveat: a can be longer than 32 bytes.
	aB := a.Bytes()

	// If we have a short byte string, expand
	// it so that it's long enough.
	aBLen := len(aB)
	if aBLen < fieldIntSize {
		diff := fieldIntSize - aBLen
		for i := 0; i < diff; i++ {
			aB = append([]byte{0x00}, aB...)
		}
	}

	for i := 0; i < fieldIntSize; i++ {
		s[i] = aB[i]
	}

	return s
}

// bigIntToEncodedBytes converts a big integer into its corresponding
// 32 byte little endian representation.
func bigIntToEncodedBytes(a *big.Int) *[32]byte {
	s := new([32]byte)
	if a == nil {
		return s
	}

	// Caveat: a can be longer than 32 bytes.
	s = copyBytes(a.Bytes())

	// Reverse the byte string --> little endian after
	// encoding.
	reverse(s)

	return s
}

// reverse reverses a byte string.
func reverse(s *[32]byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// copyBytes copies a byte slice to a 32 byte array.
func copyBytes(aB []byte) *[32]byte {
	if aB == nil {
		return nil
	}
	s := new([32]byte)

	// If we have a short byte string, expand
	// it so that it's long enough.
	aBLen := len(aB)
	if aBLen < 32 {
		diff := 32 - aBLen
		for i := 0; i < diff; i++ {
			aB = append([]byte{0x00}, aB...)
		}
	}

	for i := 0; i < 32; i++ {
		s[i] = aB[i]
	}

	return s
}

func (s *FourNodeTestSuite) transferToPubAddress(address string, meta *gsrpcTypes.Metadata) {
	// convert from ss58 to hexadecimal address
	hexAddr := polkadot.DecodeFromSS58(address)
	appendX := strings.Join([]string{"0x", hexAddr}, "")
	dest, err := gsrpcTypes.NewMultiAddressFromHexAccountID(appendX)

	if err != nil {
		panic(err)
	}

	memo := []byte(fmt.Sprintf("TRANSFER:%s", address))
	call1, err := gsrpcTypes.NewCall(meta, "System.remark", memo)

	if err != nil {
		panic(err)
	}

	amount := gsrpcTypes.NewUCompactFromUInt(10000000000000)
	call2, err := gsrpcTypes.NewCall(meta, "Balances.transfer", dest, amount)

	if err != nil {
		panic(err)
	}

	batchCall, err := gsrpcTypes.NewCall(meta, "Utility.batch_all", []gsrpcTypes.Call{call1, call2})
	if err != nil {
		panic(err)
	}

	genesisHash, err := s.polkaApi.RPC.Chain.GetBlockHash(0)

	if err != nil {
		fmt.Errorf("error %w", err)
	}

	rv, err := s.polkaApi.RPC.State.GetRuntimeVersionLatest()
	// log.Info().Interface("runtime call", rv).Msg("runtime details")
	if err != nil {
		panic(err)
	}

	ext := gsrpcTypes.NewExtrinsic(batchCall)

	mnemonic := "hero eagle luxury slight survey catch toy goat model general alarm inner"
	// mnemonic := "letter ethics correct bus asset pipe tourist vapor envelope kangaroo warm dawn"
	kp, err := signature.KeyringPairFromSecret(mnemonic, 42)

	// create storage key
	storageKey, err := gsrpcTypes.CreateStorageKey(meta, "System", "Account", kp.PublicKey)

	if err != nil {
		panic(err)
	}

	// fetch account info for nonce value
	var accountInfo gsrpcTypes.AccountInfo
	ok, err := s.polkaApi.RPC.State.GetStorageLatest(storageKey, &accountInfo)
	if !ok {
		panic(err)
	}
	nonce := accountInfo.Nonce

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

	var sub *author.ExtrinsicStatusSubscription

	sub, err = s.polkaApi.RPC.Author.SubmitAndWatchExtrinsic(ext)
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
		return
	case err := <-sub.Err():
		panic(err)
	}

}

func hash(payload []byte) []byte {
	h := sha256.New()
	h.Write(payload)
	return h.Sum(nil)
}

func hash512(payload []byte) []byte {
	h := sha512.New()
	h.Write(payload)
	return h.Sum(nil)
}

func (s *FourNodeTestSuite) getTssServer(index int, conf common.TssConfig, bootstrap string) *tss.TssServer {
	priKey, err := conversion.GetPriKey(testPriKeyArr[index])
	if err != nil {
		panic(err)
	}

	baseHome := path.Join(os.TempDir(), "4nodes_test", strconv.Itoa(index))
	if _, err := os.Stat(baseHome); os.IsNotExist(err) {
		err := os.MkdirAll(baseHome, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	var peerIDs []maddr.Multiaddr
	if len(bootstrap) > 0 {
		multiAddr, err := maddr.NewMultiaddr(bootstrap)
		if err != nil {
			panic(err)
		}
		peerIDs = []maddr.Multiaddr{multiAddr}
	} else {
		peerIDs = nil
	}

	instance, err := tss.NewTss(peerIDs, s.ports[index], priKey, "Asgard", baseHome, conf, s.preParams[index], "", "true")
	return instance
}

func getPreparams() []*btsskeygen.LocalPreParams {
	var preParamArray []*btsskeygen.LocalPreParams
	buf, err := ioutil.ReadFile(path.Join(testFileLocation, preParamTestFile))
	if err != nil {
		panic(err)
	}
	preParamsStr := strings.Split(string(buf), "\n")
	for _, item := range preParamsStr {
		var preParam btsskeygen.LocalPreParams
		val, err := hex.DecodeString(item)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(val, &preParam)
		preParamArray = append(preParamArray, &preParam)
	}
	return preParamArray
}
