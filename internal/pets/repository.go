package pets

import (
	"encoding/json"
	"errors"
	"fmt"
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
			features []Feature
			flist    []Feature
		)
		err := scanner.Scan(&id, &name, &features)
		if err != nil {
			return Species{}, err
		}
		for _, feature := range features {
			flist = append(flist, Feature{Part: feature.Part, Dominant: feature.Dominant, Recessive: feature.Recessive, Mixed: feature.Mixed})
		}
		allSpecies = append(allSpecies, Species{Id: id, Name: name, Features: flist})
	}
	speciesIndex := rand.Int() % len(allSpecies)
	species := allSpecies[speciesIndex]
	return species, nil
}

func toMap(s interface{}) (map[string]interface{}, error) {
	var newMap map[string]interface{}
	marshalled, err := json.Marshal(s)
	err = json.Unmarshal(marshalled, &newMap)
	if err != nil {
		return newMap, errors.New(fmt.Sprintf("Error converting struct to map: %v", err))
	}
	return newMap, nil
}

func (pet *Pet) persist(session *gocql.Session, ctx *gin.Context) (string, error) {
	petId := gocql.TimeUUID()
	var genes []map[string]interface{}
	for _, gene := range pet.Genes {
		readyGene, err := toMap(gene)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Error converting Genes for Cassandra: %v", err))
		}
		genes = append(genes, readyGene)
	}

	err := session.Query(`INSERT INTO lexipets.pets (id, name, species_name, species_features, genes, img) VALUES (?, ?, ?, ?, ?, ?)`, petId, pet.Name, pet.SpeciesName, pet.SpeciesFeatures, genes, pet.Img).WithContext(ctx).Exec()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error inserting into Cassandra: %v", err))
	}
	return petId.String(), nil
}
