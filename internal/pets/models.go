package pets

type Feature struct {
	Part      string `cql:"part" json:"part"`
	Dominant  string `cql:"dominant" json:"dominant"`
	Recessive string `cql:"recessive" json:"recessive"`
	Mixed     string `cql:"mixed" json:"mixed"`
}

type Species struct {
	Id       string    `cql:"id" json:"id"`
	Name     string    `cql:"name" json:"name"`
	Features []Feature `cql:"features" json:"feature"`
}

type Gene struct {
	Feature   Feature `cql:"feature" json:"feature"`
	Dominant  bool    `cql:"dominant" json:"dominant"`
	Recessive bool    `cql:"recessive" json:"recessive"`
}

type Pet struct {
	OwnerId         string    `cql:"owner_id" json:"owner_id"`
	Id              string    `cql:"id" json:"id"`
	Name            string    `cql:"name" json:"name"`
	SpeciesName     string    `cql:"species_name" json:"species_name"`
	SpeciesFeatures []Feature `cql:"species_features" json:"species_features"`
	Genes           []Gene    `cql:"genes" json:"genes"`
	Img             string    `cql:"img" json:"img"`
}
