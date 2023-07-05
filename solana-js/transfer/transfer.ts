import { Connection, Keypair, Transaction, SystemProgram, PublicKey, LAMPORTS_PER_SOL } from "@solana/web3.js";
import * as bs58 from "bs58";
import SolConnection from "./connect";
import SolKeypair from "./keypair";

// connection
const connection = new Connection("https://api.devnet.solana.com");

// 5YNmS1R9nNSCDzb5a7mMJ1dwK9uHeAAF4CmPEwKgVWr8
const feePayer = Keypair.fromSecretKey(
    bs58.decode("588FU4PktJWfGfxtzpAAXywSNt74AvtroVzGfKkVN1LwRuvHwKGr851uH8czM5qm4iqLbs1kKoMKtMJG4ATR7Ld2")
);

// G2FAbFQPFa5qKXCetoFZQEvF9BVvCKbvUZvodpVidnoY
const alice = Keypair.fromSecretKey(
    bs58.decode("4NMwxzmYj2uvHuq8xoqhY8RXg63KSVJM1DXkpbmkUY7YQWuoyQgFnnzn6yo3CMnqZasnNPNuAT2TLwQsCaKkUddp")
);

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