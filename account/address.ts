import CreateInstance from '../instance'
import MnemonicAccount from './with-seed';

(async () => {
    await CreateInstance()
    // Some mnemonic phrase
    const mnemonic = process.env.MNEMONIC as string;
    const keypair = await MnemonicAccount({mnemonic})

    // Log some info
    console.log(keypair.meta,`has address ${keypair.address} with publicKey [${keypair.publicKey}]`);

})()