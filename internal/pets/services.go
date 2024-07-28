package pets

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"math/rand"
)

func (pet *Pet) img() string {
	generatedString := fmt.Sprintf("%v-", pet.SpeciesName)
	for i, gene := range pet.Genes {
		geneString := fmt.Sprintf("-%v", gene.Feature.Part)
		switch {
		case gene.Dominant == gene.Recessive:
			geneString += fmt.Sprintf("-%v", pet.SpeciesFeatures[i].Mixed)
		case gene.Dominant:
			geneString += fmt.Sprintf("-%v", pet.SpeciesFeatures[i].Dominant)
		case gene.Recessive:
			geneString += fmt.Sprintf("-%v", pet.SpeciesFeatures[i].Recessive)
		}
		generatedString += geneString
	}

	generatedString += ".png"

	return generatedString
}

func New(session *gocql.Session, gc *gin.Context, name string) (Pet, error) {
	if name == "" {
		name = "test!"
	}

	species, err := singleSpecies(session, gc)
	if err != nil {
		return Pet{}, err
	}

	var genes []Gene
	for _, feature := range species.Features {
		dominant := rand.Int()%2 == 1
		recessive := rand.Int()%2 == 1
		newGene := Gene{Feature: feature, Dominant: dominant, Recessive: recessive}
		genes = append(genes, newGene)
	}

	pet := Pet{Name: name, SpeciesName: species.Name, SpeciesFeatures: species.Features, Genes: genes, Img: ""}

	pet.Img = pet.img()

	return pet, nil
}

func Save(session *gocql.Session, gc *gin.Context, petJson []byte) (string, error) {
	var pet Pet
	err := json.Unmarshal(petJson, &pet)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Cannot unmarshal JSON: %v", err))
	}
	petId, err := pet.persist(session, gc)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Cannot persist pet: %v", err))
	}
	return petId, nil
}

func List(session *gocql.Session, gc *gin.Context, ownerId string) ([]Pet, error) {
	petList, err := scan(session, gc, "owner_id", ownerId)

	if err != nil || petList == nil {
		return []Pet{}, err
	}

	return petList, nil
}
