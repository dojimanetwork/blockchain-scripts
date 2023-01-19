import CreateInstance from '../instance'
import {ApiPromise} from '@polkadot/api'
import type { AnyNumber, Codec } from '@polkadot/types-codec/types'
import '@polkadot/api-augment'

(async () => {
    const inst: ApiPromise = await CreateInstance()
    const block_no = process.env.BLOCK_NO as AnyNumber
    // returns Hash
    const blockHash = await inst.rpc.chain.getBlockHash(block_no);
    // returns SignedBlock
    const signedBlock = await inst.rpc.chain.getBlock(blockHash);
    const apiAt = await inst.at(signedBlock.block.header.hash);
    const allRecords = await apiAt.query.system.events();

    // the hash for the block, always via header (Hash -> toHex()) - will be
    // the same as blockHash above (also available on any header retrieved,
    // subscription or once-off)
    console.log("hash", signedBlock.block.header.hash.toHex());

    // the hash for each extrinsic in the block
    signedBlock.block.extrinsics.forEach((ex, index) => {
        // the extrinsics are decoded by the API, human-like view
        console.log(index, ex.toHuman());

        const { isSigned, meta, method: { args, method, section } } = ex;

        // explicit display of name, args & documentation
        console.log(`${section}.${method}(${args.map((a) => a.toString()).join(', ')})`);
        console.log(meta.docs.map((d) => d.toString()).join('\n'));

        // filter the specific events based on the phase and then the
        // index of our extrinsic in the block
        const events = allRecords
            .filter(({ phase }) =>
                phase.isApplyExtrinsic &&
                phase.asApplyExtrinsic.eq(index)
            )
            .map(({ event }) => `${event.section}.${event.method}`);

        console.log(`${section}.${method}:: ${events.join(', ') || 'no events'}`);

        // signer/nonce info
        if (isSigned) {
            console.log(`signer=${ex.signer.toString()}, nonce=${ex.nonce.toString()}, hash=${ex.hash.toHex()}`);
        }
    });
})()