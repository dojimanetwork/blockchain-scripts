const hre = require("hardhat");
const contracts = require("../contract.json");
const whitelistedAccounts = require("../whiteList");

async function main() {

    const nlAddr = contracts["nodelist_address"]
    const NL = await hre.ethers.getContractFactory("NodeList")
    const nl = await NL.attach(nlAddr)

    for (let i = 0; i < whitelistedAccounts.length; i++) {
        const res = await nl.IsWhitelisted(1, whitelistedAccounts[i]);
        console.log('should be whitelisted, isWhitelisted: ', res)
    }
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});