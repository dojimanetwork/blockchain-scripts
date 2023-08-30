require("@nomicfoundation/hardhat-toolbox");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: {
    compilers: [{
      version: "0.8.19",
    },
      {
        version: "0.5.0"
      }
    ]
  },
  networks: {
    dojimachain: {
      url: 'https://rpc-test.d11k.dojima.network:8545', // The URL of your custom blockchain node
      chainId: 1001, // Replace with the chain ID of your custom blockchain
      gas: "auto",
      gasPrice: 35000000000,
      gasMultiplier: 2,
      accounts: {
        mnemonic: "letter ethics correct bus asset pipe tourist vapor envelope kangaroo warm dawn",
        path: "m/44'/60'/0'/0",
        initialIndex: 0,
        count: 20,
        passphrase: "",
      }
    },
    local_eth: {
      url: "http://localhost:9545",
      chainId: 1337,
      gas: "auto",
      gasPrice: 35000000000,
      gasMultiplier: 2,
      accounts: {
        mnemonic: "letter ethics correct bus asset pipe tourist vapor envelope kangaroo warm dawn",
        path: "m/44'/60'/0'/0",
        initialIndex: 0,
        count: 20,
        passphrase: "",
      }
    }
  }
};
