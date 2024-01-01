import { Keypair } from "@solana/web3.js";
import { derivePath } from "ed25519-hd-key";
import * as bip39 from "bip39";
import fs from "fs";
import SolKeypair from "./keypair";

(() => {
    const mnemonic = process.env.MNEMONIC as string;
    const keypair = SolKeypair({mnemonic})
    fs.rm(`${process.cwd()}/keypair.json`, (err) => {
        fs.appendFile(`${process.cwd()}/keypair.json`, `[${keypair.secretKey}]`, (err) => {
            if ( err != null) {
                console.log(err)
            }
        })
    })

})()