import CreateInstance from '../instance'
import {ApiPromise, Keyring} from '@polkadot/api'

export default async function MnemonicAccount({mnemonic}: {mnemonic: string}) {
    const keyring = new Keyring({type: 'ed25519'})
    const newPair = keyring.addFromUri(mnemonic)
    return newPair
}