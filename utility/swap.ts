import CreateInstance from '../instance'
import {ApiPromise, SubmittableResult} from '@polkadot/api'
import type {EventRecord } from '@polkadot/types/interfaces';
import {stringToU8a, u8aToHex} from '@polkadot/util'
import MnemonicAccount from '../account'
(async () => {
    const inst:ApiPromise = await CreateInstance()
    // Some mnemonic phrase
    const mnemonic = 'letter ethics correct bus asset pipe tourist vapor envelope kangaroo warm dawn';
    const keypair = await MnemonicAccount({mnemonic})
    const unsub = await inst.tx.utility.batchAll(
        [inst.tx.system.remark("memo:SWAP:AR.AR:2eqRgoKFc3QhXhZWFoCoBprtep-9oiHojyaz_Ey_6os"), inst.tx.balances.transfer('5H16DLfWFLdpm5C4f9Qr6UkADsT1PtD9jELWF9WKuiC7St1T', 100000000000)
        ])
        .signAndSend(keypair, (result: SubmittableResult)=> {
            console.log("Current status is", result.status);

            if (result.status.isInBlock) {
                console.log('Transaction included in blockhash', result.status.asInBlock);
            } else if(result.status.isFinalized) {
                console.log(`Transation finalized at blockhash`, result.status.asFinalized);
                console.log("transaction hash", result.txHash);

                result.events.forEach((value: EventRecord) => {
                    console.log("\t", value.phase, ":", value.event.section, ".", value.event.method, "::::", value.event.data);

                })
                unsub()
            }
        })

})()