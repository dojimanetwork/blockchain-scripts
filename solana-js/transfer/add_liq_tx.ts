import { PublicKey } from "@solana/web3.js";
import SolConnection from "./connect";
import SolKeypair from "./keypair";
import {IDL} from "./idl_meta";
import {Program, web3, AnchorProvider, Wallet} from "@project-serum/anchor";
import FetchInboundAddr from "./inbound_addr";

export class SOLNodeWallet implements Wallet {
    constructor(readonly payer: web3.Keypair) {
        this.payer = payer
    }

    async signTransaction(tx: web3.Transaction): Promise<web3.Transaction> {
        tx.partialSign(this.payer);
        return tx;
    }

    async signAllTransactions(txs: web3.Transaction[]): Promise<web3.Transaction[]> {
        return txs.map((t) => {
            t.partialSign(this.payer);
            return t;
        });
    }

    get publicKey(): web3.PublicKey {
        return this.payer.publicKey;
    }
}

(async () => {
    // Some mnemonic phrase
    const mnemonic = process.env.MNEMONIC as string;
    const url = process.env.SOL_URL as string;
    const amount = process.env.SOL_AMT as string;
    const connection = SolConnection({url})
    const keypair = SolKeypair({mnemonic})
    const memo = process.env.SOL_ADD_LIQ as string
    const opts: web3.ConfirmOptions = {
        preflightCommitment: 'processed'
    }

    const inbound_add = await FetchInboundAddr("SOL")

    console.log(inbound_add)
    const wallet = new SOLNodeWallet(keypair)
    const provider = new AnchorProvider(connection, wallet, opts);
    const programIDPPubKey = new PublicKey('2dkwKCkTQz4xXxyjcvhUYdSb5fb3Bw15ra95o94WkyVo');
    const program = new Program(IDL, programIDPPubKey, provider);
    const txhash = await program.rpc.transferNativeTokens(`${amount}`, memo, {
        accounts: {
            from: keypair.publicKey,
            to: new web3.PublicKey(inbound_add),
            systemProgram: web3.SystemProgram.programId,
        },
        signers: [keypair],
    });
    console.log(`txhash: ${txhash}`);

})();
