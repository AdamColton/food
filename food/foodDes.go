package food

import (
	"encoding/binary"
	"fmt"
	"github.com/adamcolton/err"
	"github.com/boltdb/bolt"
	"strconv"
	"strings"
)

type FoodDes struct {
	Id           uint32
	FoodGroup    uint16
	LongDesc     string
	ShortDesc    string
	Name         string
	Manufacturer string
}

// Takes a line from FOOD_DES.txt and converts it to a FoodDes object
func FoodDesFromStr(str string) *FoodDes {
	data := splitLine(str)
	id, e := strconv.Atoi(data[0])
	err.Panic(e)
	foodGroup, e := strconv.Atoi(data[1])
	return &FoodDes{
		Id:           uint32(id),
		FoodGroup:    uint16(foodGroup),
		LongDesc:     data[2],
		ShortDesc:    data[3],
		Name:         data[4],
		Manufacturer: data[5],
	}
}

func (f *FoodDes) Key() []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, f.Id)
	return bs
}

func (f *FoodDes) Get() {
	db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(foodDesBkt)
		b := bkt.Get(f.Key())
		dec(b, f)
		return nil
	})
}

func (f *FoodDes) String() string {
	return strconv.FormatUint(uint64(f.Id), 10) + ": " + f.LongDesc
}

func (f *FoodDes) Detailed() string {
	nutData := NewNutDataGrp(f.Id)
	nutData.Get()
	s := []string{fmt.Sprintf("== %s ==", f.LongDesc)}
	for key, val := range nutData.Data {
		if val > 0 {
			nutrDef := &NutrDef{
				Id: key,
			}
			nutrDef.Get()
			if nutrDef.Display {
				s = append(s, fmt.Sprintf("%s: %.2f %s / 100 g", nutrDef.Name, val, nutrDef.Units))
			}
		}
	}
	return strings.Join(s, "\n")
}
