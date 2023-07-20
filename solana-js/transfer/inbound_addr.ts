import fetch from "cross-fetch";
import find from "lodash/find";


export default async function FetchInboundAddr(chain: string){
    var requestOptions: any = {
        method: 'GET',
        redirect: 'follow'
    };

    const result = await fetch("http://localhost:1317/hermeschain/inbound_addresses", requestOptions)
    const data = await result.json()
    const inbound_add = find(data, {chain})
    return inbound_add["address"]
}
