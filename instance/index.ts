import { ApiPromise, WsProvider } from '@polkadot/api'

export default async function CreateInstance(): Promise<ApiPromise> {
    const wsProvider = new WsProvider('wss://westend-rpc.polkadot.io')
    const api: ApiPromise = await ApiPromise.create({provider: wsProvider})
    return api
}


