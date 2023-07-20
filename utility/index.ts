import CreateInstance from '../instance'
import {ApiPromise, SubmittableResult} from '@polkadot/api'
import type {EventRecord } from '@polkadot/types/interfaces';
import '@polkadot/api-augment'

import MnemonicAccount from '../account'
import FetchInboundAddr from "../solana-js/transfer/inbound_addr";
(async () => {
    const inst:ApiPromise = await CreateInstance()
    // Some mnemonic phrase
    const mnemonic = process.env.MNEMONIC as string;
    const to_address = await FetchInboundAddr("DOT")
    const amt = process.env.AMOUNT as string
    const memo = process.env.ADD_LIQ_MEMO as string
    const keypair = await MnemonicAccount({mnemonic})
    const unsub = await inst.tx.utility.batchAll(
        [inst.tx.system.remark(memo), inst.tx.balances.transfer(to_address, amt)
        ])
        .signAndSend(keypair, (result: SubmittableResult)=> {

            if (result.status.isInBlock) {
                console.log('Transaction included in blockhash: ', "block hash - ",result.status.asInBlock.toString());
            } else if(result.status.isFinalized) {
                console.log(`Transation finalized at blockhash`, result.status.asFinalized.toString());
                console.log("transaction hash", result.txHash.toString());

                result.events.forEach((value) => {
                    console.log("\t", value.phase, ":", value.event.section, ".", value.event.method, "::::", value.event.data);

                })
                unsub()
            }
        })

})()