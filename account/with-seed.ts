import { Keyring} from '@polkadot/api'
import type { KeypairType } from '@polkadot/util-crypto/types';

export default async function MnemonicAccount({mnemonic}: {mnemonic: string}) {
    const sign_type = process.env.SIGN_TYPE as KeypairType
    const keyring = new Keyring({type: sign_type })
    const newPair = keyring.addFromUri(mnemonic)
    return newPair
}