const contracts = require("../contract.json");
const hre = require("hardhat");
const whitelistedAccounts = require("../whiteList");

async function main() {

    const nlAddr = contracts["nodelist_address"]
    const NL = await hre.ethers.getContractFactory("NodeList")
    const nl = await NL.attach(nlAddr)
    const init = await nl.initialize(1)
    const receipt = await init.wait()
    console.log(receipt)
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});