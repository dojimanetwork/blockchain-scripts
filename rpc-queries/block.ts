import CreateInstance from '../instance'
import {ApiPromise} from '@polkadot/api'
import '@polkadot/api-augment'

(async () => {
    const inst: ApiPromise = await CreateInstance()

    // returns Hash
    const blockHash = await inst.rpc.chain.getBlockHash(12782886);
    // returns SignedBlock
    const signedBlock = await inst.rpc.chain.getBlock(blockHash);

    // the hash for the block, always via header (Hash -> toHex()) - will be
    // the same as blockHash above (also available on any header retrieved,
    // subscription or once-off)
    console.log("hash", signedBlock.block.header.hash.toHex());

    // the hash for each extrinsic in the block
    signedBlock.block.extrinsics.forEach((ex, index) => {
        console.log(index, ex.hash.toHex());
    });
})()