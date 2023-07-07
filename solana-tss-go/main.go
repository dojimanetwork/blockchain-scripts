package main

import (
	"context"
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
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/dojimanetwork/dojima-tss/common"
	"github.com/dojimanetwork/dojima-tss/conversion"
	"github.com/dojimanetwork/dojima-tss/keygen"
	"github.com/dojimanetwork/dojima-tss/keysign"
	"github.com/dojimanetwork/dojima-tss/tss"
	"github.com/dojimanetwork/solana-go/v2"
	"github.com/dojimanetwork/solana-go/v2/rpc"
	"github.com/dojimanetwork/solana-go/v2/rpc/ws"
	btsskeygen "github.com/dojimanetwork/tss-lib/ecdsa/keygen"
	btss "github.com/dojimanetwork/tss-lib/tss"
	bin "github.com/gagliardetto/binary"
	maddr "github.com/multiformats/go-multiaddr"
)

const (
	partyNum         = 4
	testFileLocation = "../test_data"
	preParamTestFile = "preParam_test.data"
	endpoint         = "ws://localhost:9944"
	// as bytes.
	fieldIntSize = 32
	rpcLocal     = "http://localhost:8899"
	wsLocal      = "ws://localhost:8900"
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
	ed25519PubKey secp256k1.PublicKey
	poolPubkey    string
	rpc           rpc.Client
	ws            ws.Client
}

type TransactionInstructions struct {
	accounts  []*solana.AccountMeta
	data      []byte
	programID solana.PublicKey
}

func (in *TransactionInstructions) ProgramID() solana.PublicKey {
	return in.programID
}

func (in *TransactionInstructions) Accounts() []*solana.AccountMeta {
	return in.accounts
}

func (in *TransactionInstructions) Data() ([]byte, error) {
	return in.data, nil
}

func main() {
	common.InitLog("info", true, "four_nodes_test")
	conversion.SetupBech32Prefix()

	rpcC := rpc.New(rpcLocal)
	wsC, err := ws.Connect(context.Background(), wsLocal)

	if err != nil {
		panic(fmt.Errorf("failed to create ws client %w", err))
	}

	s := &FourNodeTestSuite{
		ports: []int{
			17666, 17667, 17668, 17669,
		},
		bootstrapPeer: "/ip4/0.0.0.0/tcp/17666/p2p/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp",
		preParams:     getPreparams(),
		servers:       make([]*tss.TssServer, partyNum),
		rpc:           *rpcC,
		ws:            *wsC,
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

	ed25519PubKey := GetPubKeyBytes(poolPubKey).Bytes()

	var publicKey solana.PublicKey

	copy(publicKey[:], ed25519PubKey)

	s.transferToPubAddress(publicKey.String())
	// get account id bytes
	// pubEncode := pubkey.SecpEncode(secpPub.ToECDSA())
	// accountId := pubkey.SecpAccountId(pubEncode)
	recentBlockHash, err := s.rpc.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)

	if err != nil {
		panic(fmt.Errorf("failed to get recent block hash %w", err))
	}

	amt := "100000000"

	ins_data := append([]byte{0x3a, 0x12, 0x3d, 0x16, 0xd0, 0xff, 0x68, 0xe7}, []byte{byte(len(amt)), 0, 0, 0}...)
	ins_data = append(ins_data, []byte(amt)...)
	ins_data = append(ins_data, []byte{byte(len("Testing")), 0, 0, 0}...)
	ins_data = append(ins_data, []byte("Testing")...)

	programID := "2dkwKCkTQz4xXxyjcvhUYdSb5fb3Bw15ra95o94WkyVo"
	dest := "FV1fdjFezEiKndqJD5DHLEnPpwpigJg8cY7HouqucfSv"

	instruction := []solana.Instruction{
		&TransactionInstructions{
			accounts: []*solana.AccountMeta{
				{PublicKey: publicKey, IsSigner: true, IsWritable: true},
				{PublicKey: solana.MustPublicKeyFromBase58(dest), IsSigner: false, IsWritable: true},
				{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
			},
			data:      ins_data,
			programID: solana.MustPublicKeyFromBase58(programID),
		},
	}

	solTx, err := solana.NewTransaction(
		instruction,
		recentBlockHash.Value.Blockhash,
		solana.TransactionPayer(publicKey),
	)

	if err != nil {
		panic(fmt.Errorf("failed to create new tx %w", err))
	}

	// convert message to bytes
	payload, err := solTx.Message.MarshalBinary()

	// payload := hash([]byte("helloworld"))
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

	var firstSig solana.Signature
out:
	for i := 0; i < len(keysignResult)-1; i++ {
		currentSignatures := keysignResult[i].Signatures
		for j := 0; j <= len(currentSignatures)-1; j++ {
			sigBytes, sig, err := getEddsaSignature(currentSignatures[j].R, currentSignatures[j].S)
			if err != nil {
				panic(fmt.Sprintf("failed to get signature:%v", err))
			}

			buf, err := base64.StdEncoding.DecodeString(currentSignatures[j].Msg)

			edPubK, err := edwards.ParsePubKey(GetPubKeyBytes(poolPubKey).Bytes())
			origSig, err := base64.StdEncoding.DecodeString(currentSignatures[i].Signature)

			if err != nil {
				fmt.Errorf("inval ed25519 key with error %w", err)
			}

			val := sig[63] & 224
			cryEdVer := ed25519.Verify(edPubK.Serialize(), buf, origSig)
			verify := edwards.Verify(edPubK, buf, sigBytes.R, sigBytes.S)
			copy(firstSig[:], origSig)
			break out
			if verify && cryEdVer && val != 0 {
				copy(firstSig[:], origSig)
				break out
			}

		}
	}

	var signatureCount []byte
	bin.EncodeCompactU16Length(&signatureCount, 1)
	output := make([]byte, 0, len(signatureCount)+len(signatureCount)*64+len(payload))
	output = append(output, signatureCount[:]...)
	output = append(output, firstSig[:]...)
	output = append(output, payload[:]...)

	opts := rpc.TransactionOpts{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
	}

	hash, err := s.rpc.SendEncodedTransactionWithOpts(context.TODO(), base64.StdEncoding.EncodeToString(output), opts)

	if err != nil {
		panic(fmt.Errorf("failed to send tx %w", err))
	}

	sub, err := s.ws.SignatureSubscribe(
		hash,
		rpc.CommitmentFinalized,
	)

	if err != nil {
		panic(err)
	}

	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			panic(err)
		}
		if got.Value.Err != nil {
			panic(fmt.Errorf("transaction confirmation failed: %v", got.Value.Err))
		} else {
			fmt.Println(hash.String())
		}
	}

}

func GetPubKeyBytes(pubkey string) ctypes.PubKey {
	unmarshal, err := legacybech32.UnmarshalPubKey(legacybech32.AccPK, pubkey)
	if err != nil {
		panic(fmt.Sprintf("unmarshal pubkey %v", err))
	}

	return unmarshal
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

func (s *FourNodeTestSuite) transferToPubAddress(address string) {
	path := "/Users/luffybhaagi/dojima/hermes/build/scripts/testnet/solana_testnet/82iP5jLLyiuTHbQRrSwUgZ6sKycT2mjbNkncgpm7Duvg.json"
	kp, err := solana.PrivateKeyFromSolanaKeygenFile(path)

	if err != nil {
		panic(fmt.Errorf("failed to get keypair %w", err))
	}

	recentBlockHash, err := s.rpc.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)

	if err != nil {
		panic(fmt.Errorf("failed to get recent block hash %w", err))
	}

	amt := "1000000000"

	ins_data := append([]byte{0x3a, 0x12, 0x3d, 0x16, 0xd0, 0xff, 0x68, 0xe7}, []byte{byte(len(amt)), 0, 0, 0}...)
	ins_data = append(ins_data, []byte(amt)...)
	ins_data = append(ins_data, []byte{byte(len("Testing")), 0, 0, 0}...)
	ins_data = append(ins_data, []byte("Testing")...)

	programID := "2dkwKCkTQz4xXxyjcvhUYdSb5fb3Bw15ra95o94WkyVo"
	dest := address
	instruction := []solana.Instruction{
		&TransactionInstructions{
			accounts: []*solana.AccountMeta{
				{PublicKey: kp.PublicKey(), IsSigner: true, IsWritable: true},
				{PublicKey: solana.MustPublicKeyFromBase58(dest), IsSigner: false, IsWritable: true},
				{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
			},
			data:      ins_data,
			programID: solana.MustPublicKeyFromBase58(programID),
		},
	}

	solTx, err := solana.NewTransaction(
		instruction,
		recentBlockHash.Value.Blockhash,
		solana.TransactionPayer(kp.PublicKey()),
	)

	if err != nil {
		panic(fmt.Errorf("failed to create new tx %w", err))
	}

	// convert message to bytes
	payload, err := solTx.Message.MarshalBinary()

	if err != nil {
		panic(fmt.Errorf("failed to convert payload %w", err))
	}

	signature, err := solTx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			return &kp
		},
	)

	if err != nil {
		panic(fmt.Errorf("failed to sign tx %w", err))
	}

	firstSig := signature[0]

	var signatureCount []byte
	bin.EncodeCompactU16Length(&signatureCount, 1)
	output := make([]byte, 0, len(signatureCount)+len(signatureCount)*64+len(payload))
	output = append(output, signatureCount[:]...)
	output = append(output, firstSig[:]...)
	output = append(output, payload[:]...)

	opts := rpc.TransactionOpts{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
	}

	hash, err := s.rpc.SendEncodedTransactionWithOpts(context.TODO(), base64.StdEncoding.EncodeToString(output), opts)

	if err != nil {
		panic(fmt.Errorf("failed to send tx %w", err))
	}

	sub, err := s.ws.SignatureSubscribe(
		hash,
		rpc.CommitmentFinalized,
	)

	if err != nil {
		panic(err)
	}

	defer sub.Unsubscribe()

	for {
		fmt.Println("waiting.....")
		got, err := sub.Recv()
		if err != nil {
			panic(err)
		}
		if got.Value.Err != nil {
			panic(fmt.Errorf("transaction confirmation failed: %v", got.Value.Err))
		} else {
			fmt.Println(hash.String())
			return
		}
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
