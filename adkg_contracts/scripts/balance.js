const {ethers} = require("hardhat");
const {utils} = require("web3");

async function main() {

    const [acc]= await ethers.getSigners()

    console.log(`${acc.address} :`,await ethers.provider.getBalance(acc.address))
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});