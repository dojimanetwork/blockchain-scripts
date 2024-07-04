import { ethers } from 'ethers'
// Replace these values with your own
const keystoreFilePath = '/Users/luffybhaagi/.dojimachain/data/keystore/UTC--2024-01-21T13-58-01.911338000Z--fad68c705bf42414b65832a100ffd3917fab2637';
const keystorePassword = 'password';
const fromAddress = '0xfAd68c705Bf42414b65832a100fFd3917Fab2637';
const toAddress = '0x71B9CDc5ba6A8747429A82C3f9E95967434BaBfA';
const valueInEther = 1; // Adjust as needed

(async () => {
    try {
        // Read keystore file
        const keystoreJson = require(keystoreFilePath);
        const wallet = await ethers.Wallet.fromEncryptedJson(keystoreJson, keystorePassword);

        // Create transaction
        const transaction = {
            to: toAddress,
            value: ethers.parseEther(valueInEther.toString()),
        };

        // Sign the transaction
        const signedTransaction = await wallet.signTransaction(transaction);

        console.log('Signed Transaction:', signedTransaction);
    } catch (error: any) {
        console.error('Error:', error.message);
    }
})()
