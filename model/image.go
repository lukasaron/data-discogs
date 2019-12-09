package model

type Image struct {
	Height string `json:"height"`
	Width  string `json:"width"`
	Type   string `json:"type"`
	Uri    string `json:"uri"`
	Uri150 string `json:"uri150"`
}

/*
<image height="337" type="primary" uri="" uri150="" width="600" />
*/
