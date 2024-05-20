package edwards

import (
	"bytes"
	"testing"

	newEd "filippo.io/edwards25519/field"
	oldEd "github.com/agl/ed25519/edwards25519"
)

func TestNewDep(t *testing.T) {
	var oldFe oldEd.FieldElement
	var newFe newEd.Element
	var oldBytes [32]byte

	copy(oldFe[:], fed[:])
	oldEd.FeToBytes(&oldBytes, &oldFe)

	newFe.SetBytes(oldBytes[:])
	newBytes := newFe.Bytes()

	if !bytes.Equal(oldBytes[:], newBytes[:]) {
		t.Fatalf("expected %x, got %x", oldBytes, newBytes)
	}
}
