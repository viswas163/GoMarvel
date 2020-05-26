package v1

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/logrusorgru/aurora"
	"github.com/viswas163/MarvelousShipt/db"
	"github.com/viswas163/MarvelousShipt/models"
)

func TestGetCommon(t *testing.T) {
	comics := []int{123, 456, 789, 234, 678, 982}
	var tests = []struct {
		name    string
		char1   models.Character
		comics1 []int
		char2   models.Character
		comics2 []int
		want    interface{}
	}{
		{
			"One Common",
			models.Character{Name: "Jean Gray"},
			[]int{123, 456, 789},
			models.Character{Name: "Emma Frost"},
			[]int{234, 456, 678},
			[]models.Comic{
				{ID: 456, Title: "abc"},
			},
		},
		{
			"Two Common",
			models.Character{Name: "The Professor"},
			[]int{123, 456, 789},
			models.Character{Name: "Cyclops"},
			[]int{123, 982, 789},
			[]models.Comic{
				{ID: 789, Title: "abc"},
				{ID: 123, Title: "abc"},
			},
		},
		{
			"No Common",
			models.Character{Name: "Cyclops"},
			[]int{123, 456, 789},
			models.Character{Name: "Emma Frost"},
			[]int{234, 456, 678},
			[]models.Comic{},
		},
	}
	p, _ := os.Getwd()
	p = filepath.Dir(filepath.Dir(p))
	_, err := db.Open(p + "/db/marvel.db")
	if err != nil {
		t.Errorf("Error init db")
	}
	defer db.GetInstance().Close()
	for _, c := range comics {
		models.ComicsByID.LoadOrStore(c, models.Comic{ID: c, Title: "abc"})
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			models.ComicsByCharacter.LoadOrStore(tt.char1.Name, tt.comics1)
			models.ComicsByCharacter.LoadOrStore(tt.char2.Name, tt.comics2)
			ans := GetCommon(tt.char1, tt.char2)
			a := reflect.ValueOf(ans)
			w := reflect.ValueOf(tt.want)
			if a.Len() > 0 && w.Len() > 0 && a.Index(0) != w.Index(0) {
				t.Errorf("got %d, want %d", ans, tt.want)
			} else {
				t.Log(tt.name, aurora.Green("Passed!"))
			}
		})
	}
}
