const hre = require("hardhat");
const contracts = require("../contract.json");
const {BigNumberish} = require("ethers")
async function main() {
    const signers = await hre.ethers.getSigners()

    const nlAddr = contracts["nodelist_address"]
    const NL = await hre.ethers.getContractFactory("NodeList")
    const nl = await NL.attach(nlAddr)
    const pubkX = await hre.ethers.getBigInt("47551540815061147812751641264987957592487251838713885348536012124335096645708")
    const pubkY = await hre.ethers.getBigInt("42009559273756857136503594680560873482640804574797420793485870047361720570328")
    const result = await nl.listNode(1, "localhost", pubkX,
        pubkY, "", "", { from: signers[0], gasLimit: '500000' })
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
