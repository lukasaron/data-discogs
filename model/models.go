package model

//--------------------------------------------------- Artist ---------------------------------------------------

//Artist is one of the main structure from Discogs:
type Artist struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	RealName       string   `json:"realName"`
	Images         []Image  `json:"images,omitempty"`
	Profile        string   `json:"profile"`
	DataQuality    string   `json:"data_quality"`
	NameVariations []string `json:"name_variations"`
	Urls           []string `json:"urls"`
	Aliases        []Alias  `json:"aliases,omitempty"`
	Members        []Member `json:"members,omitempty"`
}

// Alias is a sub structure of Artist:
type Alias struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Member is a sub structure of Artist:
type Member struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//--------------------------------------------------- Company ---------------------------------------------------

//Company structure:
type Company struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	EntityType     string `json:"entity_type"`
	EntityTypeName string `json:"entity_type_name"`
	ResourceUrl    string `json:"resource_url"`
}

//--------------------------------------------------- Format ---------------------------------------------------

// Format structure:
type Format struct {
	Name         string   `json:"name"`
	Quantity     string   `json:"quantity"`
	Text         string   `json:"text"`
	Descriptions []string `json:"description"`
}

//--------------------------------------------------- Image ---------------------------------------------------

// Image structure
type Image struct {
	Height string `json:"height"`
	Width  string `json:"width"`
	Type   string `json:"type"`
	Uri    string `json:"uri"`
	Uri150 string `json:"uri_150"`
}

//--------------------------------------------------- Label ---------------------------------------------------

// Label is one of the main structure from Discogs:
type Label struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Images      []Image      `json:"images,omitempty"`
	ContactInfo string       `json:"contact_info"`
	Profile     string       `json:"profile"`
	DataQuality string       `json:"data_quality"`
	Urls        []string     `json:"urls"`
	ParentLabel *LabelLabel  `json:"parent_label,omitempty"`
	SubLabels   []LabelLabel `json:"sub_labels,omitempty"`
}

type LabelLabel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//--------------------------------------------------- Master ---------------------------------------------------

// Master is one of the main structure from Discogs:
type Master struct {
	Id          string          `json:"id"`
	MainRelease string          `json:"main_release"`
	Images      []Image         `json:"images,omitempty"`
	Artists     []ReleaseArtist `json:"artists,omitempty"`
	Genres      []string        `json:"genres"`
	Styles      []string        `json:"styles"`
	Year        string          `json:"year"`
	Title       string          `json:"title"`
	DataQuality string          `json:"data_quality"`
	Videos      []Video         `json:"videos,omitempty"`
}

//--------------------------------------------------- Release ---------------------------------------------------

// Release is one of the main structure from Discogs:
type Release struct {
	Id           string          `json:"id"`
	Status       string          `json:"status"`
	Images       []Image         `json:"images,omitempty"`
	Artists      []ReleaseArtist `json:"artists,omitempty"`
	ExtraArtists []ReleaseArtist `json:"extra_artists,omitempty"`
	Title        string          `json:"title"`
	Formats      []Format        `json:"formats,omitempty"`
	Genres       []string        `json:"genres"`
	Styles       []string        `json:"styles"`
	Country      string          `json:"country"`
	Released     string          `json:"released"`
	Notes        string          `json:"notes"`
	DataQuality  string          `json:"data_quality"`
	MasterId     string          `json:"master_id"`
	MainRelease  string          `json:"main_release"`
	TrackList    []Track         `json:"track_list,omitempty"`
	Identifiers  []Identifier    `json:"identifiers,omitempty"`
	Videos       []Video         `json:"videos,omitempty"`
	Labels       []ReleaseLabel  `json:"labels,omitempty"`
	Companies    []Company       `json:"companies,omitempty"`
}

// Release artist is a sub structure of Release
type ReleaseArtist struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Join   string `json:"join"`
	Anv    string `json:"anv"`
	Role   string `json:"role"`
	Tracks string `json:"tracks"`
}

// Release label is a sub structure of Release
type ReleaseLabel struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

// Identifier is a sub structure of Release
type Identifier struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	Value       string `json:"value"`
}

//--------------------------------------------------- TrackList ---------------------------------------------------

// Track structure is usually part of a slice resulting into track list:
type Track struct {
	Position string `json:"position"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
}

//--------------------------------------------------- Video ---------------------------------------------------

// Video is usually part of a slice combining more videos:
type Video struct {
	Duration    string `json:"duration"`
	Embed       string `json:"embed"`
	Src         string `json:"src"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
