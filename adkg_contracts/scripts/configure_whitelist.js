
const hre = require("hardhat");

const NodeList = hre.artifacts.readArtifact('contracts/NodeList.sol:NodeList')
const whitelistedAccounts = require('../whiteList.js')
const contracts = require("../contract.json")

function tx(result, call) {
    // const logs = result.logs.length > 0 ? result.logs[0] : { address: null, event: null }

    console.log()
    console.log(`   Calling ${call}`)
    console.log('   ------------------------')
    console.log(`   > transaction hash: ${result.hash}`)
    // console.log(`   > contract address: ${logs.address}`)
    // console.log(`   > gas used: ${result.receipt.gasUsed}`)
    // console.log(`   > event: ${logs.event}`)
    console.log()
}

async function main() {
    const signers = await hre.ethers.getSigners()

    const nlAddr = contracts["nodelist_address"]
    const NL = await hre.ethers.getContractFactory("NodeList")
    const nl = await NL.attach(nlAddr)
    for (let i = 0; i < whitelistedAccounts.length; i++) {
        const acc = whitelistedAccounts[i]
        tx(await nl.updateWhitelist(5, acc, true, { from: signers[0], gas: '100000' }), `adding ${acc} to whitelist`)
    }
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
