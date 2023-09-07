import {ethers} from "ethers";


(async () => {
    const keys = ["482c4194805040e64e9c902d36373efaea7f4f5f83964081a44f2bd7ceb4c894", "0cd05860b6508db3d1534c3b22ed008e191214147ec7ef794e1878bd9a14cc70", "b490c9ecb2f47ff95de39760a5b86d423b2aa0147ae2b507bebbfd727b55bb49", "da047afc0824231a870876cb89321de362e922a23b8e4cf068473347246dd954", "e548607fed18668ac9ad3b47d56c86fc63620a1139776ed966c2e37104c7e581", "09463587b7501a5a18e86291da6403f77a17c0689527622fbf4f2a30db38f596"]
    keys.forEach(key => {
        const wallet = new ethers.Wallet(key)
        console.log(`${key} key: `, wallet.address)
    })
})()