import {HermesInit} from '@dojima-wallet/connection'
import {Network} from "@dojima-wallet/types";
import {AssetDOJNative, baseAmount} from '@dojima-wallet/utils'

(async ()=>{
    const mnemonic = process.env.MNEMONIC as string;
    const memo = process.env.HERM_MEMO as string;
    const amt = process.env.HERM_AMT as string;
    const net = process.env.HERM_NET as string
    let network: Network
    if(net == "stagenet") {
        network = Network.Stagenet
    }else {
        network = Network.Testnet
    }
    const hc = new HermesInit(mnemonic, network)
    const address = hc.h4sConnect.getAddress(0)
    console.log("address: = ", address)
    const hash = await hc.h4sConnect.deposit({walletIndex: 0,asset: AssetDOJNative, amount: baseAmount(amt, 8),memo})
    console.log(hash)
})()
