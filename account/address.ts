import CreateInstance from '../instance'
import MnemonicAccount from './with-seed';

(async () => {
    await CreateInstance()
    // Some mnemonic phrase
    const mnemonic = 'entire material egg meadow latin bargain dutch coral blood melt acoustic thought';
    const keypair = await MnemonicAccount({mnemonic})

    // Log some info
    console.log(keypair.meta,`has address ${keypair.address} with publicKey [${keypair.publicKey}]`);

})()