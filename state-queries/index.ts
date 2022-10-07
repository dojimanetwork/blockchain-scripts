import CreateInstance from '../instance'
import {ApiPromise} from '@polkadot/api'
import '@polkadot/api-augment'

(async () => {
    const inst:ApiPromise = await CreateInstance()
    // The actual address that we will use
    const ADDR = '5Gq3owRKkXLneUckXUc5UxKugXiqq78b71UQC4uHxcXFPdwH';

    // Retrieve the last timestamp
    const now = await inst.query.timestamp.now();

    // Retrieve the account balance & nonce via the system module
    const { nonce, data: balance } = await inst.query.system.account(ADDR);

    console.log(`${now}: balance of ${balance.free} and a nonce of ${nonce}`);


})()