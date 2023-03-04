package keypair

import "github.com/hyperboloide/lk"

type KeyPair struct {
	Public *lk.PublicKey
	Private *lk.PrivateKey 
}

func New(pub *lk.PublicKey, private *lk.PrivateKey) *KeyPair {
	return &KeyPair{
		Public: pub,
		Private: private,
	}
}