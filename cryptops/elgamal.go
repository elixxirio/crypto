package cryptops

import "gitlab.com/elixxir/crypto/cyclic"

type ElGamalSignature func(g *cyclic.Group, ecrKeys, cypher, keyInv, privatekey, publicCypherKey *cyclic.Int)

//Implements the modified version of ElGamal within the cryptops interface.
//Modifies ecrkeys and cypher to make its output.
//ecrkeys = ecrkeys*keyInv*(g^privatekey)%p
//cypher  = cypher*(publicCypherKey^privatekey)%p
//More details can be found in the appendix of https://drive.google.com/open?id=1ha8QtUI9Tk_sCIKWN-QE8YHZ7AKofKrV
var ElGamal ElGamalSignature = func(g *cyclic.Group, ecrKeys, cypher, keyInv, privateKey, publicCypherKey *cyclic.Int) {
	tmp := g.NewMaxInt()

	g.ExpG(privateKey, tmp)
	g.Mul(keyInv, tmp, tmp)
	g.Mul(tmp, ecrKeys, ecrKeys)

	g.Exp(publicCypherKey, privateKey, tmp)
	g.Mul(tmp, cypher, cypher)
}

func (ElGamalSignature) GetName() string {
	return "ElGamal"
}

func (ElGamalSignature) GetMinSize() uint32 {
	return 1
}
