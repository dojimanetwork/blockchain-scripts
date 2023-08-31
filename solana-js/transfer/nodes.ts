import fetch from "cross-fetch";


export default async function FetchFirstNodeAddr(){
    var requestOptions: any = {
        method: 'GET',
        redirect: 'follow'
    };

    const result = await fetch("http://localhost:1317/hermeschain/nodes", requestOptions)
    const data = await result.json()
    return data[0]["node_address"]
}
