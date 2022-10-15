import CreateInstance from '../instance'
import {stringToU8a, u8aToHex} from '@polkadot/util'
import {ApiPromise} from '@polkadot/api'
import MnemonicAccount from './with-seed';

(async () => {
    // Some mnemonic phrase
    await CreateInstance()
    const mnemonic = 'entire material egg meadow latin bargain dutch coral blood melt acoustic thought';
    const keypair = await MnemonicAccount({mnemonic})

    const message = stringToU8a("testing the signature")
    const signature = keypair.sign(message)
    const isValid = keypair.verify(message, signature, keypair.publicKey)

    // Log info
    console.log(`The signature ${u8aToHex(signature)}, is ${isValid ? '' : 'in'}valid`);

})()