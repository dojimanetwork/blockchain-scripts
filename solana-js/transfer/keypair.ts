import { Keypair } from "@solana/web3.js";
import { derivePath } from "ed25519-hd-key";
import * as bip39 from "bip39";

export default  function SolKeypair({mnemonic}: {mnemonic: string}) {
    const seed = bip39.mnemonicToSeedSync(mnemonic, ""); // (mnemonic, password)
    const path = `m/44'/501'/0'/0'`;
    const keypair = Keypair.fromSeed(derivePath(path, seed.toString("hex")).key);

    return keypair
}