import CreateInstance from '../instance'
import {ApiPromise} from '@polkadot/api'
import '@polkadot/api-augment'

(async () => {
    const inst:ApiPromise = await CreateInstance()

    // Retrieve the chain name
    const chain = await inst.rpc.system.chain();

    let count = 0
    // Subscribe to the new headers
    const unsubHeads = await inst.rpc.chain.subscribeNewHeads((lastHeader) => {
        console.log(`${chain}: last block #${lastHeader.number} has hash ${lastHeader.hash}`);
        if ( ++count === 10) {
            unsubHeads()
        }
    });
})()