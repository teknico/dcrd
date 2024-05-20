// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package edwards

import (
	"encoding/hex"
	"math/big"

	newEd "filippo.io/edwards25519/field"
	"github.com/agl/ed25519/edwards25519"
)

var (
	// zero through eight are big.Int numbers useful in
	// elliptical curve math.
	zero  = new(big.Int).SetInt64(0)
	one   = new(big.Int).SetInt64(1)
	two   = new(big.Int).SetInt64(2)
	three = new(big.Int).SetInt64(3)
	four  = new(big.Int).SetInt64(4)
	eight = new(big.Int).SetInt64(8)

	// fieldIntSize is the size of a field element encoded
	// as bytes.
	fieldIntSize = 32
)

// feOne is the field element representation of one. This is
// also the neutral (null) element.
var feOne = edwards25519.FieldElement{
	1, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
}

// fed is the field element representation of D.
var fed = edwards25519.FieldElement{
	-10913610, 13857413, -15372611, 6949391, 114729,
	-8787816, -6275908, -3247719, -18696448, -12055116,
}
var fedHex newEd.Element
fedHex, _ = new(Element).SetBytes(decodeHex(
	"a3785913ca4deb75abd841414d0a700098e879777940c78c73fe6f2bee6c0352"
))

// fed2 is the field element representation of D^2.
// 59f1b226949bd6eb56b183829a14e00030d1f3eef2808e19e7fcdf56dcd90624
var fed2 = edwards25519.FieldElement{
	-21827239, -5839606, -30745221, 13898782, 229458,
	15978800, -12551817, -6495438, 29715968, 9444199,
}

// feI is the field element representation of I.
// b0a00e4a271beec478e42fad0618432fa7d7fb3d99004d2b0bdfc14f8024832b
var feI = edwards25519.FieldElement{
	-32595792, -7943725, 9377950, 3500415, 12389472,
	-272473, -25146209, -2005654, 326686, 11406482,
}

func decodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}
