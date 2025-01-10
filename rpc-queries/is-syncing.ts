import CreateInstance from "../instance";

const { ApiPromise } = require('@polkadot/api');

async function checkSyncStatus() {
    // Replace with your node's WebSocket endpoint
    const inst = await CreateInstance()
    // Fetch sync status
    const syncState = await inst.rpc.system.health();

    // Fetch the current block number
    const state = await inst.rpc.system.syncState();



    console.log(`Current Block: ${state.currentBlock}, ${state.startingBlock}`);
    console.log(`Finalized Block: ${state.highestBlock}`);

    if (syncState.isSyncing) {
        console.log('Node is syncing...');
    } else {
        console.log('Node is fully synced.');
    }

    // Close connection
    await inst.disconnect();
}

checkSyncStatus().catch(console.error);
