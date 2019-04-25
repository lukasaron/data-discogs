package model

type Label struct {
	Id          string
	Name        string
	Images      []Image
	ContactInfo string
	Profile     string
	DataQuality string
	Urls        []string
	SubLabels   []Label
}
