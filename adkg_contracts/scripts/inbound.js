const fetch = require("cross-fetch");
const _find = require("lodash/find");

async function FetchInboundAddr(chain){
    var requestOptions = {
        method: 'GET',
        redirect: 'follow'
    };

    const result = await fetch(`http://localhost:1317/hermeschain/inbound_addresses`, requestOptions)
    const data = await result.json()
    const inbound_add = _find(data, {chain})
    return inbound_add["address"]
}