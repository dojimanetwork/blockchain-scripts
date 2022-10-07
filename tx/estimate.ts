import CreateInstance from '../instance'
import {ApiPromise} from '@polkadot/api'
import MnemonicAccount from '../account/with-seed';

(async () => {
    const inst: ApiPromise = await CreateInstance()
    // Some mnemonic phrase
    const mnemonic = 'letter ethics correct bus asset pipe tourist vapor envelope kangaroo warm dawn';
    const keypair = await MnemonicAccount({mnemonic})

    const info = await inst.tx.utility.batchAll(
        [inst.tx.system.remark("testing"), inst.tx.balances.transfer('5DTestUPts3kjeXSTMyerHihn1uwMfLj8vU8sqF7qYrFabHE', 10000000)
        ]).paymentInfo(keypair)

    // log relevant info, partialFee is Balance, estimated for current
    console.log(`
      class=${info.class.toString()},
      weight=${info.weight.toString()},
      partialFee=${info.partialFee.toHuman()},
      PlanckFee=${info.partialFee}
    `);

})()