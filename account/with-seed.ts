import { Keyring} from '@polkadot/api'

export default async function MnemonicAccount({mnemonic}: {mnemonic: string}) {
    const keyring = new Keyring({type: 'sr25519'})
    const newPair = keyring.addFromUri(mnemonic)
    return newPair
}