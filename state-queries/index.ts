import CreateInstance from '../instance'
import {ApiPromise} from '@polkadot/api'
import '@polkadot/api-augment'

(async () => {
    const inst:ApiPromise = await CreateInstance()
    const address = process.env.FROM_ADDRESS as string
    // The actual address that we will use
    const ADDR = address;

    // Retrieve the last timestamp
    const now = await inst.query.timestamp.now();

    // Retrieve the account balance & nonce via the system module
    const { nonce, data: balance } = await inst.query.system.account(ADDR);

    console.log(`${now}: balance of ${balance.free} and a nonce of ${nonce}`);


})()