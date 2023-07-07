import {Connection, LAMPORTS_PER_SOL, PublicKey} from "@solana/web3.js";
import SolConnection from "./connect";

(async () => {

// connection
    const url = process.env.SOL_URL as string;
    const connection = SolConnection({url})
    const dest = process.env.SOL_DEST as string
    const tokenAccount1Pubkey = new PublicKey(dest);
    let balance = await connection.getBalance(tokenAccount1Pubkey);
    console.log(`sol: ${balance/ LAMPORTS_PER_SOL}`);
})();