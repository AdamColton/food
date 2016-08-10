package food

import (
	"bytes"
	"github.com/adamcolton/assert"
	"testing"
)

func TestNutrDefFromStr(t *testing.T) {
	a := assert.A{t}
	nutr := NutrDefFromStr("~203~^~g~^~PROCNT~^~Protein~^~2~^~600~")
	a.True(nutr.Id == 203, "Expected 203")
	a.True(nutr.Units == "g", "Expected 'g'")
	a.True(nutr.Tag == "PROCNT", "Expected 'PROCNT'")
	a.True(nutr.Name == "Protein", "Expected 'Protein'")

	expect := []byte{203, 0}
	if !bytes.Equal(expect, nutr.Key()) {
		t.Error("Did not get expected key")
	}
}
