package tss

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

	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32"
	"github.com/decred/dcrd/dcrec/edwards"
	"github.com/dojimanetwork/dojima-tss/common"
	"github.com/dojimanetwork/dojima-tss/conversion"
	"github.com/dojimanetwork/dojima-tss/keygen"
	"github.com/dojimanetwork/dojima-tss/keysign"
	"github.com/dojimanetwork/dojima-tss/tss"
	gsrpc "github.com/dojimanetwork/go-polka-rpc/v5"
	gsrpcTypes "github.com/dojimanetwork/go-polka-rpc/v5/types"
	common2 "github.com/dojimanetwork/hermes/common"
	"github.com/dojimanetwork/hermes/narada/chainclients/polkadot"
	btsskeygen "github.com/dojimanetwork/tss-lib/ecdsa/keygen"
	btss "github.com/dojimanetwork/tss-lib/tss"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog/log"
	cryptoEd "golang.org/x/crypto/ed25519"
)

const (
	partyNum         = 4
	testFileLocation = "../test_data"
	preParamTestFile = "preParam_test.data"
	endpoint         = "ws://localhost:9944"
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
			16666, 16667, 16668, 16669,
		},
		bootstrapPeer: "/ip4/127.0.0.1/tcp/16666/p2p/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp",
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

	if err != nil {
		panic(err)
	}

	address, err := pubkey.GetAddress(common2.DOTCHAIN)
	log.Info().Interface("dot address", address).Msg("POLKA::::")
	ed25519PubKey, err := pubkey.GetEd25519PubK()
	if err != nil {
		panic(err)
	}

	meta, err := s.polkaApi.RPC.State.GetMetadataLatest()

	if err != nil {
		panic(err)
	}

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
	payload, err := polkadot.GetEd25519Payload(extrinsic, sigOpts)

	if err != nil {
		panic(err)
	}

	keysignResult := make(map[int]keysign.Response)
	for i := 0; i < partyNum; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			localPubKeys := append([]string{}, testPubKeys...)
			var keysignReq keysign.Request
			if newJoinParty {
				keysignReq = keysign.NewRequest(poolPubKey, []string{base64.StdEncoding.EncodeToString(payload)}, 10, localPubKeys, "0.14.0", s.algo)
			} else {
				keysignReq = keysign.NewRequest(poolPubKey, []string{base64.StdEncoding.EncodeToString(payload)}, 10, localPubKeys, "0.13.0", s.algo)
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

	checkSignResult(keysignResult, poolPubKey, payload)
}

func checkSignResult(keysignResult map[int]keysign.Response, poolPubkey string, payload []byte) {
	for i := 0; i < len(keysignResult)-1; i++ {
		currentSignatures := keysignResult[i].Signatures
		sigBytes, err := getSignature(currentSignatures[0].R, currentSignatures[0].S)
		if err != nil {

		}
		s := sigBytes[63] & 224
		pk, err := legacybech32.UnmarshalPubKey(legacybech32.AccPK, poolPubkey)
		if err != nil {

		}

		edPubK := cryptoEd.PublicKey(pk.Bytes())

		verify := cryptoEd.Verify(edPubK, payload, sigBytes)

		fmt.Println(verify, s)
	}
}

func getSignature(r, s string) ([]byte, error) {
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

func (s *FourNodeTestSuite) getPolkaMsgToSign(meta *gsrpcTypes.Metadata) (gsrpcTypes.SignatureOptions, gsrpcTypes.Extrinsic) {
	hexaAddr := strings.Join([]string{"0x", "d2c2e63069b7422f37f5c6bb6cf4241d406eb0bb33a8333649a6b77151244c2e"}, "")
	dest, err := gsrpcTypes.NewMultiAddressFromHexAccountID(hexaAddr)
	if err != nil {
		panic(err)
	}

	memo := []byte("memo:OUT:E250EBC0EBF271ED23C41B23D5024C65BAE5563819F7537E63605EEA86485839")
	call1, err := gsrpcTypes.NewCall(meta, "System.remark", memo)

	if err != nil {
		panic(err)
	}

	amount := gsrpcTypes.NewUCompactFromUInt(346506515540)
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

	signOpts := gsrpcTypes.SignatureOptions{
		BlockHash:          genesisHash, // using genesis since we're using immortal era
		Era:                gsrpcTypes.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              gsrpcTypes.NewUCompactFromUInt(uint64(0)),
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
