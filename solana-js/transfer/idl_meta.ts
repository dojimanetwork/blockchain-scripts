export const IDL: Dojima = {
    "version": "0.1.0",
    "name": "dojima",
    "instructions": [
        {
            "name": "transferNativeTokens",
            "accounts": [
                {
                    "name": "from",
                    "isMut": true,
                    "isSigner": true
                },
                {
                    "name": "to",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "systemProgram",
                    "isMut": false,
                    "isSigner": false
                }
            ],
            "args": [
                {
                    "name": "tokenAmount",
                    "type": "string"
                },
                {
                    "name": "memo",
                    "type": "string"
                }
            ]
        }
    ]
};

export type Dojima = {
    "version": "0.1.0",
    "name": "dojima",
    "instructions": [
        {
            "name": "transferNativeTokens",
            "accounts": [
                {
                    "name": "from",
                    "isMut": true,
                    "isSigner": true
                },
                {
                    "name": "to",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "systemProgram",
                    "isMut": false,
                    "isSigner": false
                }
            ],
            "args": [
                {
                    "name": "tokenAmount",
                    "type": "string"
                },
                {
                    "name": "memo",
                    "type": "string"
                }
            ]
        }
    ]
};