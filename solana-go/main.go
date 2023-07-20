package main

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/dojimanetwork/solana-go/v2"
	"github.com/dojimanetwork/solana-go/v2/rpc"
	"github.com/dojimanetwork/solana-go/v2/rpc/ws"
	bin "github.com/gagliardetto/binary"
)

const (
	rpcLocal = "http://localhost:8899"
	wsLocal  = "ws://localhost:8900"
)

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
	rpcC := rpc.New(rpcLocal)
	wsC, err := ws.Connect(context.Background(), wsLocal)

	if err != nil {
		panic(fmt.Errorf("failed to create ws client %w", err))
	}

	path := "/Users/luffybhaagi/dojima/hermes/build/scripts/testnet/solana_testnet/82iP5jLLyiuTHbQRrSwUgZ6sKycT2mjbNkncgpm7Duvg.json"
	kp, err := solana.PrivateKeyFromSolanaKeygenFile(path)

	if err != nil {
		panic(fmt.Errorf("failed to get keypair %w", err))
	}

	out, err := rpcC.GetBalance(
		context.TODO(),
		kp.PublicKey(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		panic(err)
	}
	spew.Dump(out)

	recentBlockHash, err := rpcC.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)

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

	hash, err := rpcC.SendEncodedTransactionWithOpts(context.TODO(), base64.StdEncoding.EncodeToString(output), opts)

	if err != nil {
		panic(fmt.Errorf("failed to send tx %w", err))
	}

	sub, err := wsC.SignatureSubscribe(
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
