package food

import (
	"github.com/adamcolton/assert"
	"testing"
)

func TestNutDataFromStr(t *testing.T) {
	a := assert.A{t}
	nutData := NutDataFromStr("~01001~^~203~^0.85^16^0.074^~1~^~~^~~^~~^^^^^^^~~^11/1976^")
	a.True(nutData.FoodId == 1001, "Expected 1001")
	a.True(nutData.NutrId == 203, "Expected 203")
	a.True(nutData.Val == 0.85, "Expected 0.85")
}
