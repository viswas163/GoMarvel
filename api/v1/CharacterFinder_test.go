package v1

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/viswas163/MarvelousShipt/db"
	"github.com/viswas163/MarvelousShipt/models"
)

func TestGetCommon(t *testing.T) {
	var tests = []struct {
		name    string
		char1   models.Character
		comics1 []int
		char2   models.Character
		comics2 []int
		want    []models.Comic
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
				{ID: 123, Title: "abc"},
				{ID: 789, Title: "abc"},
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
	comics := []int{123, 456, 789, 234, 678, 982}
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

			if len(ans) > 0 && len(tt.want) > 0 && !cmp.Equal(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}
