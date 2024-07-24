package pets

import (
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"math/rand"
)

func singleSpecies(session *gocql.Session, ctx *gin.Context) (Species, error) {
	var allSpecies []Species
	scanner := session.Query(`SELECT id, name, features FROM lexipets.species`).WithContext(ctx).Iter().Scanner()
	for scanner.Next() {
		var (
			id       string
			name     string
			features []map[string]string
			flist    []Feature
		)
		err := scanner.Scan(&id, &name, &features)
		if err != nil {
			return Species{}, err
		}
		for _, feature := range features {
			flist = append(flist, Feature{Part: feature["part"], Dominant: feature["dominant"], Recessive: feature["recessive"], Mixed: feature["mixed"]})
		}
		allSpecies = append(allSpecies, Species{Id: id, Name: name, Features: flist})
	}
	speciesIndex := rand.Int() % len(allSpecies)
	species := allSpecies[speciesIndex]
	return species, nil
}
