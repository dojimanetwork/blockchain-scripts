import { Transaction, SystemProgram, PublicKey } from "@solana/web3.js";
import SolConnection from "./connect";
import SolKeypair from "./keypair";

(async () => {
    // Some mnemonic phrase
    const mnemonic = process.env.MNEMONIC as string;
    const url = process.env.SOL_URL as string;
    const amount = process.env.SOL_TRANSFER_AMT as string;
    const dest = process.env.SOL_DEST as string
    const connection = SolConnection({url})
    const keypair = SolKeypair({mnemonic})

    let tx = new Transaction().add(
        SystemProgram.transfer({
            fromPubkey: keypair.publicKey,
            toPubkey: new PublicKey(dest),
            lamports: Number(amount),
        })
    );
    tx.feePayer = keypair.publicKey;

    let txhash = await connection.sendTransaction(tx, [keypair]);
    console.log(`txhash: ${txhash}`);
})();