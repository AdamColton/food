package food

import (
	"bytes"
	"github.com/adamcolton/assert"
	"testing"
)

func TestFoodDesFromStr(t *testing.T) {
	a := assert.A{t}
	foodDes := FoodDesFromStr("~01001~^~0100~^~Butter, salted~^~BUTTER,WITH SALT~^~~^~~^~Y~^~~^0^~~^6.38^4.27^8.79^3.87")
	a.True(foodDes.Id == 1001, "Expected 1001")
	a.True(foodDes.FoodGroup == 100, "Expected 100")
	a.True(foodDes.LongDesc == "Butter, salted", "Expected 'Butter, salted'")
	a.True(foodDes.ShortDesc == "BUTTER,WITH SALT", "Expected 'BUTTER,WITH SALT'")
	a.True(foodDes.Name == "", "Expected ''")
	a.True(foodDes.Manufacturer == "", "Expected ''")

	expect := []byte{233, 3, 0, 0}
	if !bytes.Equal(expect, foodDes.Key()) {
		t.Error("Did not get expected key")
	}
}
