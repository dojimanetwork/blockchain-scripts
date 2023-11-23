const hre = require("hardhat");
const contracts = require("../contract.json");

async function main() {
    const signers = await hre.ethers.getSigners()

    const nlAddr = contracts["nodelist_address"]
    const NL = await hre.ethers.getContractFactory("NodeList")
    const nl = await NL.attach(nlAddr)
    const pubkX = await hre.ethers.getBigInt("51432792828285186105341724747602456458180317510569626148064061823289530889378")
    const pubkY = await hre.ethers.getBigInt("2414052357091825218403691575415080113564244696039162350189568725686031397863")
    const result = await nl.listNode(5, "localhost:8083", pubkX,
        pubkY, "a620f5ac96d484cd80149afd993b862a1eb96a36@127.0.0.1:26956", "/ip4/127.0.0.1/tcp/1083/p2p/16Uiu2HAmLJsKV6kn8Yq7S8kKdynNvJg6g1sPPWNPKGpDckXQ9QSD", { from: signers[0], gasLimit: '500000' })
    const receipt = await result.wait()
    console.log(receipt)
    const logs = receipt.logs?.length > 0 ? receipt.logs[0] : null
    console.log(logs)
}


// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
