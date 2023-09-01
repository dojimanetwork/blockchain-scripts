import {ethers} from "ethers";


(async () => {

const wallet = new ethers.Wallet("0cd05860b6508db3d1534c3b22ed008e191214147ec7ef794e1878bd9a14cc70")
    console.log(wallet.address)
})()