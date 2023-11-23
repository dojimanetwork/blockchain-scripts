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
            data: utils.utf8ToHex("SWAP:DOJ.DOJ:0xa14655a5e856564341b4a659eff54e1932c9afd3"),
            chainId: '1337'
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