import { LAMPORTS_PER_SOL } from "@solana/web3.js";
import SolConnection from "./connect";
import SolKeypair from "./keypair";

(async () => {
    // Some mnemonic phrase
    const mnemonic = process.env.MNEMONIC as string;
    const url = process.env.SOL_URL as string;
    const connection = SolConnection({url})
    const keypair = SolKeypair({mnemonic})
    let balance = await connection.getBalance(keypair.publicKey);
    console.log(`${balance / LAMPORTS_PER_SOL} SOL`);
})();