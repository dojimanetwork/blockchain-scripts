const hre = require("hardhat");
const contracts = require("../contract.json");
const whitelistedAccounts = require("../whiteList");

async function main() {
    const signers = await hre.ethers.getSigners()
    const nlAddr = contracts["nodelist_address"]
    const NL = await hre.ethers.getContractFactory("NodeList")
    const nl = await NL.attach(nlAddr)
    const result = await nl.updateEpoch(5, 3, 2, 2, [], 4, 6,{ from: signers[0], gasLimit: '500000' })
    const receipt = await result.wait()
    console.log(receipt)
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});