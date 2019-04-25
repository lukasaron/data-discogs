package model

type Artist struct {
	Id             string
	Name           string
	RealName       string
	Images         []Image
	Profile        string
	DataQuality    string
	NameVariations []string
	Urls           []string
	Aliases        []Artist
}
