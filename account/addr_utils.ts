import * as ethers from "ethers"
import bip39 from 'bip39'
// @ts-ignore
import HDKey from 'hdkey'
// @ts-ignore
import CoinKey from 'coinkey'
import {networks} from 'bitcoinjs-lib'

import HDNode from "avalanche/dist/utils/hdnode"
import { Avalanche, Mnemonic, Buffer } from "avalanche"
import { EVMAPI, KeyChain } from "avalanche/dist/apis/evm"
import {SigningKey} from "ethers";



interface Wallet  {
    privKey: string,
    pubKey: string,
    address: string

}

export function getEthWallet(mnemonic: string): Wallet {
    const mnemonicWallet = ethers.Wallet.fromPhrase(mnemonic)
    const privKey = mnemonicWallet.privateKey
    const pubKey = mnemonicWallet.publicKey
    const address = mnemonicWallet.address
    return {privKey, pubKey, address}
}

export function getBtcWallet(_mnemonic: string): Wallet {
// Convert the mnemonic to a seed
    const mnemonic: Mnemonic = Mnemonic.getInstance()
    const seed = mnemonic.mnemonicToSeedSync(_mnemonic);

// Create an HD wallet key from the seed
    // @ts-ignore
    const hdKey = HDKey.fromMasterSeed(Buffer.from(seed, "hex"));

// Define the BIP44 path for Bitcoin (m/44'/0'/0'/0/0)
    const path = "m/44'/0'/0'/0/0";

// Derive a child key from the HD key using the defined path
    const child = hdKey.derive(path);

// Create a CoinKey from the derived child private key and specify the Bitcoin network
    const coinKey = new CoinKey(child.privateKey, networks.bitcoin);

    const address = coinKey.publicAddress
    const privKey = coinKey.privateKey.toString("hex")
    const pubKey = ""
    return {address, privKey, pubKey}
}

export function getAvaxWallet(_mnemonic: string): Wallet {
    const ip: string = "api.avax-test.network"
    const port: number = 443
    const protocol: string = "https"
    const networkID: number = 5
    const avalanche: Avalanche = new Avalanche(ip, port, protocol, networkID)
    const cchain: EVMAPI = avalanche.CChain()
    const mnemonic: Mnemonic = Mnemonic.getInstance()
   const seed: Buffer = mnemonic.mnemonicToSeedSync(_mnemonic)
    const hdnode: HDNode = new HDNode(seed)
    const keyChain: KeyChain = cchain.newKeyChain()
    const child: HDNode = hdnode.derive(`m/44'/60'/0'/0/0`)
    keyChain.importKey(child.privateKey)
    const address = ethers.computeAddress(new SigningKey(child.privateKey))
    const privKey = child.privateKey.toString("hex")
    return {address, privKey, pubKey: ""}
}