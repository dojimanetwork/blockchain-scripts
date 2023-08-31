
const hre = require("hardhat");

const NodeList = hre.artifacts.readArtifact('contracts/NodeList.sol:NodeList')
const whitelistedAccounts = require('../whiteList.js')


function tx(result, call) {
    const logs = result.logs.length > 0 ? result.logs[0] : { address: null, event: null }

    console.log()
    console.log(`   Calling ${call}`)
    console.log('   ------------------------')
    console.log(`   > transaction hash: ${result.tx}`)
    console.log(`   > contract address: ${logs.address}`)
    console.log(`   > gas used: ${result.receipt.gasUsed}`)
    console.log(`   > event: ${logs.event}`)
    console.log()
}

async function main() {
    const signers = await hre.ethers.getSigners()
    const NodeListInstance = await NodeList
    tx(await NodeListInstance.updateEpoch(1, 5, 3, 1, [], 0, 2), 'Updated epoch')
    for (let i = 0; i < whitelistedAccounts.length; i++) {
        const acc = whitelistedAccounts[i]
        tx(await NodeListInstance.updateWhitelist(1, acc, true, { from: signers[0], gas: '100000' }), `adding ${acc} to whitelist`)
    }

    for (var i = 0; i < whitelistedAccounts.length; i++) {
        var res = await NodeListInstance.IsWhitelisted(1, whitelistedAccounts[i])
        console.log('should be whitelisted, isWhitelisted: ', res)
    }
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
