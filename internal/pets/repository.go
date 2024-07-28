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
	if len(allSpecies) == 0 {
		return Species{}, errors.New("no species found, is there a db error?")
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

	err := session.Query(`INSERT INTO lexipets.pets (owner_id, id, name, species_name, species_features, genes, img) VALUES (?, ?, ?, ?, ?, ?, ?)`, pet.OwnerId, petId, pet.Name, pet.SpeciesName, pet.SpeciesFeatures, genes, pet.Img).WithContext(ctx).Exec()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error inserting into Cassandra: %v", err))
	}
	return petId.String(), nil
}

func scan(session *gocql.Session, ctx *gin.Context, fieldName string, fieldValue string) ([]Pet, error) {
	var (
		ownerId         string
		id              string
		name            string
		speciesName     string
		speciesFeatures []Feature
		genes           []Gene
		img             string
	)
	statement := fmt.Sprintf(`SELECT owner_id, id, name, species_name, species_features, genes, img FROM lexipets.pets WHERE %v = ? ALLOW FILTERING`, fieldName)
	scanner := session.Query(statement, fieldValue).Iter().Scanner()

	var List []Pet
	for scanner.Next() {
		err := scanner.Scan(&ownerId, &id, &name, &speciesName, &speciesFeatures, &genes, &img)
		if err != nil {
			return []Pet{}, errors.New(fmt.Sprintf("Error while fetching pets. Original error: %v", err))
		}
		List = append(List, Pet{OwnerId: ownerId, Id: id, Name: name, SpeciesName: speciesName, SpeciesFeatures: speciesFeatures, Genes: genes, Img: img})
	}
	return List, nil
}
