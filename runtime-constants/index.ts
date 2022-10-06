import CreateInstance from '../instance'
import {ApiPromise} from '@polkadot/api'

(async () => {
    const inst:ApiPromise = await CreateInstance()
    console.log("constants ----- \n",inst.consts)
})()
