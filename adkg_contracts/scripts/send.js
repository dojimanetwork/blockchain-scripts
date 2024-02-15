const { ethers} = require("hardhat");
const {utils} = require("web3")


async function main() {

    const [acc]= await ethers.getSigners()

    console.log(`${acc.address} :`,await ethers.provider.getBalance(acc.address))
    const nonce = await acc.getNonce()
    const tx = await acc.sendTransaction({
            to: '0x934E5123fb2D0507b7C4B8A402E2879610268B4f',
            from: acc.address,
            nonce: nonce,
            value: ethers.parseEther('100'),
            data: utils.utf8ToHex("Bandit network faucet"),
            chainId: '1001'
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