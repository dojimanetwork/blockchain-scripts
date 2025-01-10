const WebSocket = require('ws');
const buffer = require("buffer");

const ws = new WebSocket('ws://20.204.129.56:8546'); // Replace with your WebSocket URL

ws.on('open', function open() {
    console.log('Connected');

    // Send a ping every 30 seconds
    setInterval(() => {
        if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ jsonrpc: "2.0", id: 1, method: "eth_blockNumber", params: [] }));
            console.log('Ping sent');
        }
    }, 30000);
});

ws.on('message', function incoming(data) {


    console.log("Parsed JSON Object:", data);
});

ws.on('close', function close(code, reason) {
    console.log(`Disconnected (code: ${code}, reason: ${reason})`);
});

ws.on('error', function error(err) {
    console.error('WebSocket error:', err);
});