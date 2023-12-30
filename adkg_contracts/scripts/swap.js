const { ethers} = require("hardhat");
const {utils} = require("web3")


async function main() {

    const [acc]= await ethers.getSigners()

    console.log(`${acc.address} :`,await ethers.provider.getBalance(acc.address))
    const nonce = await acc.getNonce()
    const tx = await acc.sendTransaction({
            to: '0xd526d5f47f863eff32b99bc4f9e77ddb4bd2929b',
            from: acc.address,
            nonce: nonce,
            value: ethers.parseEther('1'),
            data: utils.utf8ToHex("SWAP:ETH.ETH:0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"),
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