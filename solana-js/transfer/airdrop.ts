import { Connection, Keypair } from "@solana/web3.js";
import * as bs58 from "bs58";
import SolConnection from "./connect";
import SolKeypair from "./keypair";

(async () => {
    // Some mnemonic phrase
    const mnemonic = process.env.MNEMONIC as string;
    const url = process.env.SOL_URL as string;
    const amount = process.env.SOL_AMT as string;
    const connection = SolConnection({url})
    const keypair = SolKeypair({mnemonic})
    // 1e9 lamports = 10^9 lamports = 1 SOL
    let txhash = await connection.requestAirdrop(keypair.publicKey, Number(amount));
    console.log(`txhash: ${txhash}`);
})();