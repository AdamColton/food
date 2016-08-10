package food

import (
	"encoding/binary"
	"errors"
	"github.com/adamcolton/err"
	"github.com/boltdb/bolt"
	"strconv"
)

type NutData struct {
	FoodId uint32
	NutrId uint16
	Val    float32
}

// Takes a line from NUT_DATA.txt and converts it to a NutData object
func NutDataFromStr(str string) *NutData {
	data := splitLine(str)
	err.Debug(data)
	foodId, e := strconv.Atoi(data[0])
	err.Panic(e)
	nutrId, e := strconv.Atoi(data[1])
	err.Panic(e)
	val, e := strconv.ParseFloat(data[2], 32)
	err.Panic(e)
	return &NutData{
		FoodId: uint32(foodId),
		NutrId: uint16(nutrId),
		Val:    float32(val),
	}
}

type NutDataGrp struct {
	FoodId uint32
	Data   map[uint16]float32 // nutrId => amount per 100g
}

func NewNutDataGrp(foodId uint32) *NutDataGrp {
	return &NutDataGrp{
		FoodId: foodId,
		Data:   map[uint16]float32{},
	}
}

func (n *NutDataGrp) Key() []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, n.FoodId)
	return bs
}

func (n *NutDataGrp) Get() {
	db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(nutDataBkt)
		b := bkt.Get(n.Key())
		dec(b, n)
		return nil
	})
}

var FoodIdMismatch = errors.New("Food IDs do not match")

func (n *NutDataGrp) AddData(data *NutData) error {
	if n.FoodId != data.FoodId {
		return FoodIdMismatch
	}
	n.Data[data.NutrId] = data.Val
	return nil
}

type nutDataAggregator map[uint32]*NutDataGrp

func (n nutDataAggregator) AddData(data *NutData) {
	grp, ok := n[data.FoodId]
	if !ok {
		grp = NewNutDataGrp(data.FoodId)
		n[data.FoodId] = grp
	}

	grp.AddData(data)
}
