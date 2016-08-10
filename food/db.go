package food

import (
	"github.com/adamcolton/err"
	"github.com/boltdb/bolt"
)

var foodDesBkt = []byte("foodDes")
var nutrDefBkt = []byte("nutrDef")
var nutDataBkt = []byte("nutData")
var searchBkt = []byte("search")
var recipeBkt = []byte("recipe")

var db *bolt.DB

func init() {
	var e error
	db, e = bolt.Open("food.db", 0600, nil)
	err.Panic(e)

	db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists(foodDesBkt)
		err.Panic(e)
		_, e = tx.CreateBucketIfNotExists(nutrDefBkt)
		err.Panic(e)
		_, e = tx.CreateBucketIfNotExists(nutDataBkt)
		err.Panic(e)
		_, e = tx.CreateBucketIfNotExists(searchBkt)
		err.Panic(e)
		_, e = tx.CreateBucketIfNotExists(recipeBkt)
		err.Panic(e)
		return nil
	})
}

type keyed interface {
	Key() []byte
}

// writes a record to a bucket
func write(bktId []byte, obj keyed) {
	db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bktId).Put(obj.Key(), enc(obj))
	})
}
