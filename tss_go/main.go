package main

import (
	"crypto/sha256"
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

	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32"
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
	common2 "github.com/dojimanetwork/hermes/common"
	"github.com/dojimanetwork/hermes/common/cosmos"
	"github.com/dojimanetwork/hermes/narada/chainclients/polkadot"
	btsskeygen "github.com/dojimanetwork/tss-lib/ecdsa/keygen"
	btss "github.com/dojimanetwork/tss-lib/tss"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ed25519"
)

const (
	partyNum         = 4
	testFileLocation = "../test_data"
	preParamTestFile = "preParam_test.data"
	endpoint         = "ws://localhost:9944"
	// as bytes.
	fieldIntSize = 32
)

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
	ed25519PubKey ed25519.PublicKey
	poolPubkey    string
	sigOpts       gsrpcTypes.SignatureOptions
	extrinsic     gsrpcTypes.Extrinsic
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
	s.poolPubkey = poolPubKey
	pubkey, err := common2.NewPubKey(poolPubKey)

	log.Info().Interface("ed25519 pubkey", pubkey).Msg("POLKA::::")

	if err != nil {
		panic(err)
	}

	address, err := pubkey.GetAddress(common2.DOTCHAIN)
	log.Info().Interface("dot address", address).Msg("POLKA::::")

	ed25519PubKey, err := pubkey.GetEd25519PubK()
	s.ed25519PubKey = ed25519PubKey
	if err != nil {
		panic(err)
	}

	meta, err := s.polkaApi.RPC.State.GetMetadataLatest()

	if err != nil {
		panic(err)
	}

	s.transferToPubAddress(address.String(), meta)

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
	sigOpts, extrinsic := s.getPolkaMsgToSign(meta)
	s.sigOpts = sigOpts
	s.extrinsic = extrinsic
	payload, err := polkadot.GetEd25519Payload(extrinsic, sigOpts)

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
				keysignReq = keysign.NewRequest(poolPubKey, []string{base64.StdEncoding.EncodeToString(hash(payload))}, 10, localPubKeys, "0.14.0", s.algo)
			} else {
				keysignReq = keysign.NewRequest(poolPubKey, []string{base64.StdEncoding.EncodeToString(hash(payload))}, 10, localPubKeys, "0.13.0", s.algo)
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

	s.checkSignResult(keysignResult, poolPubKey)
}

func GetPubKeyBytes(pubkey string) []byte {
	unmarshal, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, pubkey)
	if err != nil {
		panic(fmt.Sprintf("unmarshal pubkey %v", err))
	}

	return unmarshal.Bytes()
}

func (s *FourNodeTestSuite) checkSignResult(keysignResult map[int]keysign.Response, poolPubkey string) {
	var signature gsrpcTypes.Signature
out:
	for i := 0; i < len(keysignResult)-1; i++ {
		currentSignatures := keysignResult[i].Signatures
		for j := 0; j <= len(currentSignatures)-1; j++ {
			sigBytes, sig, err := getSignature(currentSignatures[j].R, currentSignatures[j].S)
			if err != nil {

			}

			pk, err := legacybech32.UnmarshalPubKey(legacybech32.AccPK, poolPubkey)
			if err != nil {

			}

			edPubK, err := edwards.ParsePubKey(pk.Bytes())

			if err != nil {
				fmt.Errorf("inval ed25519 key with error %w", err)
			}

			buf, err := base64.StdEncoding.DecodeString(currentSignatures[j].Msg)
			verify := edwards.Verify(edPubK, buf, sigBytes.R, sigBytes.S)

			if verify {

				signature = gsrpcTypes.NewSignature(sig)
				break out
			}

		}
	}

	multiSig := gsrpcTypes.MultiSignature{IsEd25519: true, AsEd25519: signature}

	multiSigPubkey, err := gsrpcTypes.NewMultiAddressFromAccountID(GetPubKeyBytes(s.poolPubkey))

	if err != nil {
		panic(fmt.Sprintf("failed to create new account id %v", err))
	}

	extSig := gsrpcTypes.ExtrinsicSignatureV4{
		Signer:    multiSigPubkey,
		Signature: multiSig,
		Era:       s.sigOpts.Era,
		Nonce:     s.sigOpts.Nonce,
		Tip:       s.sigOpts.Tip,
	}

	s.extrinsic.Signature = extSig

	s.extrinsic.Version |= gsrpcTypes.ExtrinsicBitSigned

	// enc, err := codec.EncodeToHex(signature)
	//
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to encode to hex %v", err))
	// }
	//
	// sub, err := s.polkaApi.RPC.Author.SubmitBytesAndWatchExtrinsic(enc)

	// var sub *author.ExtrinsicStatusSubscription

	hash, err := s.polkaApi.RPC.Author.SubmitExtrinsic(s.extrinsic)
	fmt.Println(hash)
	// if err != nil {
	// 	panic(err)
	// 	// continue
	// }
	//
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to subscribe %v", err))
	// }
	//
	// defer sub.Unsubscribe()
	//
	// select {
	// case <-time.After(1 * time.Minute):
	// 	panic("Timeout reached")
	// case st := <-sub.Chan():
	// 	extStatus, _ := st.MarshalJSON()
	// 	fmt.Println("Done with status -", string(extStatus))
	// 	return
	// case err := <-sub.Err():
	// 	panic(err)
	// }

}

func getSignature(r, s string) (*edwards.Signature, []byte, error) {
	rBytes, err := base64.StdEncoding.DecodeString(r)
	if err != nil {
		return nil, nil, err
	}
	sBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, nil, err
	}

	R := encodedBytesToBigInt(copyBytes(rBytes))
	S := encodedBytesToBigInt(copyBytes(sBytes))

	sig := append([]byte{0}, append(rBytes[:], sBytes[:]...)...)
	return &edwards.Signature{R: R, S: S}, sig, nil
}

// encodedBytesToBigInt converts a 32 byte little endian representation of
// an integer into a big, big endian integer.
func encodedBytesToBigInt(s *[32]byte) *big.Int {
	// Use a copy so we don't screw up our original
	// memory.

	bi := new(big.Int).SetBytes(s[:])

	return bi
}

// bigIntToEncodedBytes converts a big integer into its corresponding
// 32 byte little endian representation.
func bigIntToEncodedBytes(a *big.Int) *[32]byte {
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

	// // Reverse the byte string --> little endian after
	// // encoding.
	// reverse(s)

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

func (s *FourNodeTestSuite) getPolkaMsgToSign(meta *gsrpcTypes.Metadata) (gsrpcTypes.SignatureOptions, gsrpcTypes.Extrinsic) {
	hexaAddr := strings.Join([]string{"0x", "d2c2e63069b7422f37f5c6bb6cf4241d406eb0bb33a8333649a6b77151244c2e"}, "")
	dest, err := gsrpcTypes.NewMultiAddressFromHexAccountID(hexaAddr)
	if err != nil {
		panic(err)
	}

	// memo := []byte("memo:OUT:TSS_TESTING")
	// call1, err := gsrpcTypes.NewCall(meta, "System.remark", memo)
	//
	// if err != nil {
	// 	panic(err)
	// }

	amount := gsrpcTypes.NewUCompactFromUInt(346506515540)
	call2, err := gsrpcTypes.NewCall(meta, "Balances.transfer", dest, amount)

	if err != nil {
		panic(err)
	}

	// batchCall, err := gsrpcTypes.NewCall(meta, "Utility.batch_all", []gsrpcTypes.Call{call1, call2})
	// if err != nil {
	// 	panic(err)
	// }

	genesisHash, err := s.polkaApi.RPC.Chain.GetBlockHash(0)

	if err != nil {
		fmt.Errorf("error %w", err)
	}

	rv, err := s.polkaApi.RPC.State.GetRuntimeVersionLatest()
	// log.Info().Interface("runtime call", rv).Msg("runtime details")
	if err != nil {
		panic(err)
	}

	ext := gsrpcTypes.NewExtrinsic(call2)

	// create storage key
	storageKey, err := gsrpcTypes.CreateStorageKey(meta, "System", "Account", s.ed25519PubKey)

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

	return signOpts, ext
}

func hash(payload []byte) []byte {
	h := sha256.New()
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
