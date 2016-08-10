package food

import (
	"encoding/binary"
	"github.com/adamcolton/err"
	"github.com/boltdb/bolt"
	"strconv"
)

type NutrDef struct {
	Id      uint16
	Units   string
	Tag     string
	Name    string
	Display bool
}

// Takes a line from NUTR_DEF.txt and converts it to a NutrDef object
func NutrDefFromStr(str string) *NutrDef {
	data := splitLine(str)
	id, e := strconv.Atoi(data[0])
	err.Panic(e)
	return &NutrDef{
		Id:    uint16(id),
		Units: data[1],
		Tag:   data[2],
		Name:  data[3],
	}
}

func (n *NutrDef) Key() []byte {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, n.Id)
	return bs
}

// Gets the search term from the database
func (n *NutrDef) Get() {
	db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(nutrDefBkt)
		b := bkt.Get(n.Key())
		dec(b, n)
		return nil
	})
}

func (n *NutrDef) Save() {
	write(nutrDefBkt, n)
}

func AllNutrients() []*NutrDef {
	nutrDefs := []*NutrDef{}

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(nutrDefBkt))

		b.ForEach(func(k, v []byte) error {
			var n NutrDef
			dec(v, &n)
			nutrDefs = append(nutrDefs, &n)
			return nil
		})
		return nil
	})

	return nutrDefs
}
