package food

import (
	"errors"
	"os"
)

//https://www.ars.usda.gov/SP2UserFiles/Place/80400525/Data/SR/SR28/sr28_doc.pdf
const food_des_url = "https://www.ars.usda.gov/SP2UserFiles/Place/12354500/Data/SR/SR28/asc/FOOD_DES.txt"
const nut_data_url = "https://www.ars.usda.gov/SP2UserFiles/Place/80400525/Data/SR/SR28/asc/NUT_DATA.txt"
const nutr_def_url = "https://www.ars.usda.gov/SP2UserFiles/Place/12354500/Data/SR/SR28/asc/NUTR_DEF.txt"

// Streams FoodDes objects from FOOD_DES.txt
func ReadFoodDes(f *os.File) <-chan *FoodDes {
	ch := make(chan *FoodDes)
	go func(ch chan<- *FoodDes) {
		for line := range readFileByLine(f) {
			ch <- FoodDesFromStr(line)
		}
		close(ch)
	}(ch)
	return ch
}

// Streams NutrDef objects from NUTR_DEF.txt
func ReadNutrDef(f *os.File) <-chan *NutrDef {
	ch := make(chan *NutrDef)
	go func(ch chan<- *NutrDef) {
		for line := range readFileByLine(f) {
			ch <- NutrDefFromStr(line)
		}
		close(ch)
	}(ch)
	return ch
}

// Streams NutData objects from NUT_DATA.txt
func ReadNutData(f *os.File) <-chan *NutData {
	ch := make(chan *NutData)
	go func(ch chan<- *NutData) {
		for line := range readFileByLine(f) {
			ch <- NutDataFromStr(line)
		}
		close(ch)
	}(ch)
	return ch
}

// Reads data from .txt files and stores it in the bolt db
func PopulateDB() error {

	foodDesFile, e := os.Open("FOOD_DES.txt")
	if e != nil {
		return errors.New("Could not open FOOD_DES.txt, download from " + food_des_url)
	}
	defer foodDesFile.Close()

	nutrDefFile, e := os.Open("NUTR_DEF.txt")
	if e != nil {
		return errors.New("Could not open NUTR_DEF.txt, download from " + nutr_def_url)
	}
	defer nutrDefFile.Close()

	nutDataFile, e := os.Open("NUT_DATA.txt")
	if e != nil {
		return errors.New("Could not open NUT_DATA.txt, download from " + nut_data_url)
	}
	defer nutDataFile.Close()

	for nutrDef := range ReadNutrDef(nutrDefFile) {
		write(nutrDefBkt, nutrDef)
	}

	search := searchAggregator{}
	for foodDes := range ReadFoodDes(foodDesFile) {
		write(foodDesBkt, foodDes)
		search.Add(foodDes.Id, foodDes.LongDesc)
	}
	for _, searchTerm := range search {
		write(searchBkt, searchTerm)
	}

	agg := nutDataAggregator{}
	for nutData := range ReadNutData(nutDataFile) {
		agg.AddData(nutData)
	}
	for _, nutDataGrp := range agg {
		write(nutDataBkt, nutDataGrp)
	}

	return nil
}
