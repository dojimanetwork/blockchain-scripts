import CreateInstance from '../instance'
import {ApiPromise} from '@polkadot/api'
import MnemonicAccount from '../account/with-seed';

(async () => {
    const inst: ApiPromise = await CreateInstance()
    // Some mnemonic phrase
    const mnemonic = process.env.MNEMONIC as string;
    const to_address = process.env.TO_ADDRESS as string;
    const amt = process.env.AMOUNT as string
    const keypair = await MnemonicAccount({mnemonic})

    const info = await inst.tx.utility.batchAll(
        [inst.tx.system.remark("testing"), inst.tx.balances.transfer(to_address, amt)
        ]).paymentInfo(keypair)

    // log relevant info, partialFee is Balance, estimated for current
    console.log(`
      class=${info.class.toString()},
      weight=${info.weight.toString()},
      partialFee=${info.partialFee.toHuman()},
      PlanckFee=${info.partialFee}
    `);

})()