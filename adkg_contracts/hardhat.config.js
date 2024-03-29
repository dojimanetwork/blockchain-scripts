require("@nomicfoundation/hardhat-toolbox");
// hardhat.config.js
require('@openzeppelin/hardhat-upgrades');


/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: {
    compilers: [{
      version: "0.8.19",
    },
      {
        version: "0.5.0"
      },
      {
        version: "0.8.8"
      }
    ]
  },
  networks: {
    dojimachain: {
      url: 'http://localhost:8545', // The URL of your custom blockchain node
      chainId: 1001, // Replace with the chain ID of your custom blockchain
      gas: "auto",
      gasPrice: 35000000000,
      gasMultiplier: 2,
      // accounts: [
      //   "0cd05860b6508db3d1534c3b22ed008e191214147ec7ef794e1878bd9a14cc70"
      // ]
      // accounts: [
      //   "b490c9ecb2f47ff95de39760a5b86d423b2aa0147ae2b507bebbfd727b55bb49"
      // ]
      accounts: [
        "da047afc0824231a870876cb89321de362e922a23b8e4cf068473347246dd954"
      ]
      // accounts: [
      //   "e548607fed18668ac9ad3b47d56c86fc63620a1139776ed966c2e37104c7e581"
      // ]
      // accounts: [
      //   "09463587b7501a5a18e86291da6403f77a17c0689527622fbf4f2a30db38f596"
      // ]
      // accounts: {
      //   mnemonic: "letter ethics correct bus asset pipe tourist vapor envelope kangaroo warm dawn",
      //   path: "m/44'/60'/0'/0",
      //   initialIndex: 0,
      //   count: 20,
      //   passphrase: "",
      // }
    },
    local_eth: {
      url: "http://localhost:9545",
      chainId: 1337,
      gas: "auto",
      gasPrice: 35000000000,
      gasMultiplier: 2,
      accounts: [
          "09463587b7501a5a18e86291da6403f77a17c0689527622fbf4f2a30db38f596"
      ]
      // accounts: {
      //   mnemonic: "letter ethics correct bus asset pipe tourist vapor envelope kangaroo warm dawn",
      //   path: "m/44'/60'/0'/0",
      //   initialIndex: 0,
      //   count: 20,
      //   passphrase: "",
      // }
    },
    avalanche: {
      url: "http://127.0.0.1:9652/ext/bc/C/rpc",
      chainId: 43112 ,
      gas: 5000000, //units of gas you are willing to pay, aka gas limit
      gasPrice: 225000000000, //gas is typically in units of gwei, but you must enter it as wei here
      accounts:[
        "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"
      ]
    },
    fuji:{
      url: "https://ava-testnet.public.blastapi.io/ext/bc/C/rpc",
      chainId: 43113 ,
      gasPrice: 25000000000,
      gas: 21000,
    }
  }
};
