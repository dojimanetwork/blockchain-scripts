import {generateMnemonic} from "bip39";


(async () => {
    const mnemonic = generateMnemonic(256)
    console.log(mnemonic)
})()