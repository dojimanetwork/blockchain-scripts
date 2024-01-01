const { ethers} = require("hardhat");
const {utils} = require("web3")


async function main() {

    const [acc]= await ethers.getSigners()

    console.log(`${acc.address} :`,await ethers.provider.getBalance(acc.address))
    const nonce = await acc.getNonce()
    const tx = await acc.sendTransaction({
            to: '0xD9233f6D8a37167d77314DaA87EafB48ad4aBD47',
            from: acc.address,
            nonce: nonce,
            value: ethers.parseEther('0.499475'),
            chainId: '43113'
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