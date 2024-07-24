package pets

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImageStringReferencesCombinedGenes(t *testing.T) {
	featuresList := []Feature{
		Feature{
			Part:      "tail",
			Dominant:  "curly",
			Recessive: "straight",
			Mixed:     "curved",
		},
		Feature{
			Part:      "ears",
			Dominant:  "pointy",
			Recessive: "floppy",
			Mixed:     "perky",
		},
		Feature{
			Part:      "color",
			Dominant:  "brown",
			Recessive: "yellow",
			Mixed:     "mottled",
		},
		Feature{
			Part:      "paws",
			Dominant:  "big",
			Recessive: "small",
			Mixed:     "medium",
		},
	}
	species := Species{Id: "1", Name: "Waggler", Features: featuresList}
	pet := Pet{Species: species, Genes: []Gene{
		Gene{
			featuresList[0],
			true,
			true,
		},
		Gene{
			featuresList[1],
			true,
			false,
		},
		Gene{
			featuresList[2],
			false,
			true,
		},
		Gene{
			featuresList[3],
			false,
			false,
		},
	}}
	expectedUrl := "Waggler--tail-curved-ears-pointy-color-yellow-paws-medium.png"
	resultUrl := pet.img()

	assert.Equal(t, expectedUrl, resultUrl)
}
