package model

type Element struct {
	Name    string     `json:"name"`
	Recipes [][]string `json:"recipes"`
	Tier    int        `json:"tier"`
}

type Data struct {
	Elements []Element `json:"elements"`
}
