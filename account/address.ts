import CreateInstance from '../instance'
import MnemonicAccount from './with-seed';
import {getAvaxWallet, getBtcWallet, getEthWallet} from "./addr_utils";
import SolKeypair from "../solana-js/transfer/keypair";

(async () => {
    // await CreateInstance()
    // Some mnemonic phrase
    const mnemonic = process.env.MNEMONIC as string;
    // const keypair = await MnemonicAccount({mnemonic})

    const ethWallet = getEthWallet(mnemonic)
    const avaxWallet = getAvaxWallet(mnemonic)
    const solKp = SolKeypair({ mnemonic})
    const btcWallet = getBtcWallet(mnemonic)


    // Log some info
    // console.log("polkadot address", keypair.meta,`has address ${keypair.address} with publicKey [${keypair.publicKey}]`);
    console.log("ethereum address", ethWallet.address);
    console.log("avalanche address", avaxWallet.address)
    console.log("solana address", solKp.publicKey.toString())
    console.log("bitcoin address", btcWallet.address)

})()