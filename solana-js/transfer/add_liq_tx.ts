import { Transaction, SystemProgram, PublicKey } from "@solana/web3.js";
import SolConnection from "./connect";
import SolKeypair from "./keypair";
import {IDL} from "./idl_meta";
import {Program, web3, Provider, AnchorProvider, Wallet} from "@project-serum/anchor";

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
    const dest = process.env.SOL_DEST as string
    const connection = SolConnection({url})
    const keypair = SolKeypair({mnemonic})
    const memo = process.env.SOL_ADD_LIQ as string
    const opts: web3.ConfirmOptions = {
        preflightCommitment: 'processed'
    }

    const wallet = new SOLNodeWallet(keypair)
    const provider = new AnchorProvider(connection, wallet, opts);
    const programIDPPubKey = new PublicKey('2dkwKCkTQz4xXxyjcvhUYdSb5fb3Bw15ra95o94WkyVo');
    const program = new Program(IDL, programIDPPubKey, provider);
    const txhash = await program.rpc.transferNativeTokens(`${amount}`, memo, {
        accounts: {
            from: keypair.publicKey,
            to: new web3.PublicKey(dest),
            systemProgram: web3.SystemProgram.programId,
        },
        signers: [keypair],
    });
    console.log(`txhash: ${txhash}`);

})();
