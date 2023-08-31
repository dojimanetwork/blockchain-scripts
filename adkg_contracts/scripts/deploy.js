// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// You can also run a script with `npx hardhat run <script>`. If you do that, Hardhat
// will compile your contracts, add the Hardhat Runtime Environment's members to the
// global scope, and execute the script.
const hre = require("hardhat");
const fs = require("fs")
async function main() {

  const [acc] = await hre.ethers.getSigners()
  console.log(await hre.ethers.provider.getBalance(acc.address))
  const contracts = ["Migrations", "NodeList"]
    const Migration = await hre.ethers.getContractFactory(contracts[0])
    const migration = await Migration.deploy()
    await migration.waitForDeployment()

    const NodeList = await hre.ethers.getContractFactory(contracts[1])
    const nodelist = await NodeList.deploy()
    await nodelist.waitForDeployment()
    console.log('%c \n migration address:', 'color:', await migration.getAddress() );
    console.log('%c \n nodelist address:', 'color:', await nodelist.getAddress() );
    const migrationAddr = await migration.getAddress()
  const nodeListAddr = await nodelist.getAddress()

    await fs.rm(`${process.cwd()}/contract.json`, (err) => {
      // ignore the error
      // console.log(err)
    })

    fs.appendFile(`${process.cwd()}/contract.json`, `\n { \n "migration_address": "${migrationAddr}",\n "nodelist_address": "${nodeListAddr}"\n }`, (err) => {
      console.log(err)
    })

}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
