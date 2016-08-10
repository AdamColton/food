package food

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"strings"
	"text/tabwriter"
)

type Recipe struct {
	Name        string
	Ingredients map[uint32]float32 // food => amount
}

func (r *Recipe) Key() []byte {
	return []byte(r.Name)
}

func (r *Recipe) Save() {
	write(recipeBkt, r)
}

func (r *Recipe) Get() {
	db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(recipeBkt)
		b := bkt.Get(r.Key())
		dec(b, r)
		return nil
	})
	if r.Ingredients == nil {
		r.Ingredients = map[uint32]float32{}
	}
}

func (r *Recipe) Add(foodID uint32, amount float32) {
	if amount == 0 {
		delete(r.Ingredients, foodID)
	} else {
		r.Ingredients[foodID] = amount
	}
	r.Save()
}

func (r *Recipe) String() string {
	s := []string{fmt.Sprintf("== %s ==", r.Name)}
	for key, val := range r.Ingredients {
		foodDes := &FoodDes{
			Id: key,
		}
		foodDes.Get()
		s = append(s, fmt.Sprintf("%10.2f g %s", val, foodDes.LongDesc))
	}
	return strings.Join(s, "\n")
}

func (r *Recipe) Detailed() string {
	details := map[uint32]map[uint16]float32{} // foodId => nutrientId => amount
	nutrTotals := map[uint16]float32{}
	nutrCache := map[uint16]*NutrDef{}
	foodCache := map[uint32]*FoodDes{}

	for key, ingrdtAmg := range r.Ingredients {
		foodDes := &FoodDes{
			Id: key,
		}
		foodDes.Get()
		foodCache[key] = foodDes
		nutData := NewNutDataGrp(foodDes.Id)
		nutData.Get()
		details[foodDes.Id] = map[uint16]float32{}
		for key, nutrAmt := range nutData.Data {
			if nutrAmt > 0 {
				nutrDef, ok := nutrCache[key]
				if !ok {
					nutrDef = &NutrDef{
						Id: key,
					}
					nutrDef.Get()
					nutrCache[key] = nutrDef
				}
				if nutrDef.Display {
					nt := nutrTotals[nutrDef.Id]
					amt := (ingrdtAmg / 100.0) * nutrAmt
					details[foodDes.Id][nutrDef.Id] = amt
					nt += amt
					nutrTotals[nutrDef.Id] = nt
				}
			}
		}
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 2, 0, '\t', tabwriter.AlignRight)

	nutrList := []uint16{}
	// build header
	fmt.Fprint(w, " \t|")
	for nutrId, _ := range nutrTotals {
		nutrList = append(nutrList, nutrId)
		nutrDef := nutrCache[nutrId]
		fmt.Fprint(w, " ", nutrDef.Name, "\t|")
	}
	fmt.Fprintln(w)

	// build data
	for foodId, nutrs := range details {
		foodDes := foodCache[foodId]
		fmt.Fprint(w, foodDes.LongDesc, "\t|")
		for _, nutrId := range nutrList {
			nutrAmt := nutrs[nutrId]
			fmt.Fprintf(w, " %.2f\t|", nutrAmt)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprint(w, "TOTAL\t|")
	for _, nutrId := range nutrList {
		nutrAmt := nutrTotals[nutrId]
		fmt.Fprintf(w, " %.2f\t|", nutrAmt)
	}
	w.Flush()
	return string(buf.Bytes())
}

func AllRecipes() []string {
	recipes := []string{}

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(recipeBkt))

		b.ForEach(func(k, v []byte) error {
			recipes = append(recipes, string(k))
			return nil
		})
		return nil
	})

	return recipes
}
