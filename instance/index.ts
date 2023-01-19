import { ApiPromise, WsProvider } from '@polkadot/api'
import dotenv from 'dotenv'

dotenv.config()

export default async function CreateInstance(): Promise<ApiPromise> {
    const host= process.env.HOST as string
    const wsProvider = new WsProvider(host)
    const api: ApiPromise = await ApiPromise.create({provider: wsProvider})
    return api
}


