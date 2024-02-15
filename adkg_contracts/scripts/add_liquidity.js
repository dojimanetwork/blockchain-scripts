const { ethers} = require("hardhat");
const {utils} = require("web3")


async function main() {

    const [acc]= await ethers.getSigners()

    console.log(`${acc.address} :`,await ethers.provider.getBalance(acc.address))
    const nonce = await acc.getNonce()
    const tx = await acc.sendTransaction({
            to: '0x23a9a914d6c325e033355ff7faf56fd142af10e1',
            from: acc.address,
            nonce: nonce,
            value: ethers.parseEther('10'),
            data: utils.utf8ToHex("ADD:ETH.ETH:tdojima19jqq0xqle6lh4ggr6rrrvw5utz2k52q2h4kgnv"),
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