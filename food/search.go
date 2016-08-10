package food

import (
	"github.com/boltdb/bolt"
	"regexp"
	"strings"
)

type SearchTerm struct {
	Term    string
	FoodIds []uint32
}

func (s *SearchTerm) Key() []byte {
	return []byte(s.Term)
}

func intersection(searchTerms ...*SearchTerm) *SearchTerm {
	s := &SearchTerm{
		Term:    "",
		FoodIds: []uint32{},
	}
	counts := map[uint32]uint8{}
	for _, st := range searchTerms {
		s.Term += st.Term + " "
		for _, foodId := range st.FoodIds {
			c := counts[foodId]
			counts[foodId] = c + 1
		}
	}
	s.Term = s.Term[:len(s.Term)-1] // remove last space
	terms := uint8(len(searchTerms))
	for foodId, c := range counts {
		if c == terms {
			s.FoodIds = append(s.FoodIds, foodId)
		}
	}
	return s
}

// Gets the search term from the database
func (s *SearchTerm) Get() {
	db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(searchBkt)
		dec(bkt.Get(s.Key()), s)
		return nil
	})
}

func Search(s string) *SearchTerm {
	searchWords := getWords(s)
	terms := []*SearchTerm{}
	for _, word := range searchWords {
		s := &SearchTerm{
			Term: word,
		}
		s.Get()
		terms = append(terms, s)
	}

	return intersection(terms...)
}

func (s *SearchTerm) Foods() []*FoodDes {
	foods := make([]*FoodDes, len(s.FoodIds))
	for i, foodId := range s.FoodIds {
		f := &FoodDes{
			Id: foodId,
		}
		f.Get()
		foods[i] = f
	}
	return foods
}

var wordsRe = regexp.MustCompile("[a-zA-Z]+")

// gets a slice of lower case words
func getWords(str string) []string {
	return wordsRe.FindAllString(strings.ToLower(str), -1)
}

type searchAggregator map[string]*SearchTerm

func (s searchAggregator) Add(foodId uint32, desc string) {
	for _, str := range getWords(desc) {
		searchTerm, ok := s[str]
		if !ok {
			searchTerm = &SearchTerm{
				Term:    str,
				FoodIds: []uint32{},
			}
			s[str] = searchTerm
		}
		searchTerm.FoodIds = append(searchTerm.FoodIds, foodId)
	}
}
