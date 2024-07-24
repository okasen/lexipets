package pets

type Feature struct {
	Part      string `json:"part"`
	Dominant  string `json:"dominant"`
	Recessive string `json:"recessive"`
	Mixed     string `json:"mixed"`
}

type Species struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Features []Feature `json:"features"`
}

type Gene struct {
	Feature   Feature `json:"feature"`
	Dominant  bool    `json:"dominant"`
	Recessive bool    `json:"recessive"`
}

type Pet struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Species Species `json:"species"`
	Genes   []Gene  `json:"genes"`
	Img     string  `json:"img"`
}
