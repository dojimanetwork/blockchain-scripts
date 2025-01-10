const { ethers} = require("hardhat");
const {utils} = require("web3")
const fetch = require("cross-fetch");
const _find = require("lodash/find");


async function FetchInboundAddr(chain){
    var requestOptions = {
        method: 'GET',
        redirect: 'follow'
    };

    const result = await fetch(`http://localhost:1317/hermeschain/inbound_addresses`, requestOptions)
    const data = await result.json()
    const inbound_add = _find(data, {chain})
    return inbound_add["address"]
}

async function main() {

    const [acc]= await ethers.getSigners()
    const to_address = await FetchInboundAddr("DOJ")
    console.log(`${acc.address} :`,await ethers.provider.getBalance(acc.address))
    const nonce = await acc.getNonce()
    const tx = await acc.sendTransaction({
            to: to_address,
            from: acc.address,
            nonce: nonce,
            value: ethers.parseEther('1'),
            data: utils.utf8ToHex("SWAP:DOT.DOT:5GKzQ1XFAxnHRbiHKGYoWbtzDMhdpb8dTEY46LBba2S5yB7L"),
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