import CreateInstance from '../instance'
import {ApiPromise} from '@polkadot/api'
import '@polkadot/api-augment'

(async () => {
    const inst:ApiPromise = await CreateInstance()

    // Retrieve the chain name
    const chain = await inst.rpc.system.chain();

    // Subscribe to the new headers
    await inst.rpc.chain.subscribeNewHeads((lastHeader) => {
        console.log(`${chain}: last block #${lastHeader.number} has hash ${lastHeader.hash}`);
    });
})()