const { ethers} = require("hardhat");
const {utils} = require("web3")


async function main() {

    const [acc]= await ethers.getSigners()

    console.log(`${acc.address} :`,await ethers.provider.getBalance(acc.address))
    const nonce = await acc.getNonce()
    const tx = await acc.sendTransaction({
            to: '0xB82D838b5E3E2Dc7A2610F9C228F69E6edBcDCCD',
            from: acc.address,
            nonce: nonce,
            value: ethers.parseEther('100'),
            data: utils.utf8ToHex("faucet ahmedabad"),
            chainId: '184'
        }
    )
    console.log(tx)
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});