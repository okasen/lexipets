package pets

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"math/rand"
)

func (pet *Pet) img() string {
	generatedString := fmt.Sprintf("%v-", pet.Species.Name)
	for i, gene := range pet.Genes {
		geneString := fmt.Sprintf("-%v", gene.Feature.Part)
		switch {
		case gene.Dominant == gene.Recessive:
			geneString += fmt.Sprintf("-%v", pet.Species.Features[i].Mixed)
		case gene.Dominant:
			geneString += fmt.Sprintf("-%v", pet.Species.Features[i].Dominant)
		case gene.Recessive:
			geneString += fmt.Sprintf("-%v", pet.Species.Features[i].Recessive)
		}
		generatedString += geneString
	}

	generatedString += ".png"

	return generatedString
}

func New(session *gocql.Session, gc *gin.Context) (Pet, error) {
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

	pet := Pet{Name: "Alex", Species: species, Genes: genes, Img: ""}

	pet.Img = pet.img()

	return pet, nil
}
