import { Connection } from "@solana/web3.js";

export default function SolConnection({url} : {url :string}) {
    // connection
    const connection = new Connection(url);
    return connection
}
