const { ethers} = require("hardhat");
const {utils} = require("web3")


async function main() {

    const [acc]= await ethers.getSigners()

    console.log(`${acc.address} :`,await ethers.provider.getBalance(acc.address))
    const tx = await ethers.provider.getTransaction("0xdc4d0e694f1fd93fdd8318f2cefe442f88449032989636dc17214a770265fc29")

    console.log(tx)
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});