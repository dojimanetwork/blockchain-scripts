import fetch from "cross-fetch";


export default async function FetchFirstNodeAddr(){
    var requestOptions: any = {
        method: 'GET',
        redirect: 'follow'
    };
    const endpoint =process.env.HERMES_ENDPOINT as string
    const result = await fetch(`${endpoint}/hermeschain/nodes`, requestOptions)
    const data = await result.json()
    return data[0]["node_address"]
}
