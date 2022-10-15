import { ApiPromise, WsProvider } from '@polkadot/api'

export default async function CreateInstance(): Promise<ApiPromise> {
    const wsProvider = new WsProvider('ws://127.0.0.1:9944')
    const api: ApiPromise = await ApiPromise.create({provider: wsProvider})
    return api
}


