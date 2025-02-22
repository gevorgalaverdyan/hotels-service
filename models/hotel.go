package models

type Hotel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	WikiLink    string `json:"wikiLink"`
	City        string `json:"city"`
	Province    string `json:"province"`
	Image       string `json:"image"`
	Coordinates string `json:"coordinates"`
	Website     string `json:"website"`
}
