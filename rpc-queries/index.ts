import CreateInstance from '../instance'
import {ApiPromise} from '@polkadot/api'
import '@polkadot/api-augment'

(async () => {
    const inst:ApiPromise = await CreateInstance()

    // without subscription
    // Retrieve the chain name
    const chain = await inst.rpc.system.chain();

    // Retrieve the latest header
    const lastHeader = await inst.rpc.chain.getHeader();

    // Log the information
    console.log(`${chain}: last block #${lastHeader.number} has hash ${lastHeader.hash}`);

})()