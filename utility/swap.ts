import CreateInstance from '../instance'
import {ApiPromise, SubmittableResult} from '@polkadot/api'
import type {EventRecord } from '@polkadot/types/interfaces';
import {stringToU8a, u8aToHex} from '@polkadot/util'
import MnemonicAccount from '../account'
(async () => {
    const inst:ApiPromise = await CreateInstance()
    // Some mnemonic phrase
    const mnemonic = process.env.MNEMONIC as string;
    const to_address = process.env.TO_ADDRESS as string;
    const amt = process.env.AMOUNT as string
    const memo = process.env.SWAP_MEMO as string
    const keypair = await MnemonicAccount({mnemonic})
    const unsub = await inst.tx.utility.batchAll(
        [inst.tx.system.remark(memo), inst.tx.balances.transfer(to_address, amt)
        ])
        .signAndSend(keypair, (result: SubmittableResult)=> {
            console.log("Current status is", result.status);

            if (result.status.isInBlock) {
                console.log('Transaction included in blockhash', result.status.asInBlock);
            } else if(result.status.isFinalized) {
                console.log(`Transation finalized at blockhash`, result.status.asFinalized);
                console.log("transaction hash", result.txHash);

                result.events.forEach((value) => {
                    console.log("\t", value.phase, ":", value.event.section, ".", value.event.method, "::::", value.event.data);

                })
                unsub()
            }
        })

})()