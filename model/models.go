// Copyright (C) 2020  Lukas Aron
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package model expresses data structures based on XML files provided by Discogs.
package model

//--------------------------------------------------- Artist ---------------------------------------------------

//Artist is one of the main structure from Discogs:
type Artist struct {
	ID             string   `json:"id"`
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
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Member is a sub structure of Artist:
type Member struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//--------------------------------------------------- Company ---------------------------------------------------

//Company structure:
type Company struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	EntityType     string `json:"entity_type"`
	EntityTypeName string `json:"entity_type_name"`
	ResourceURL    string `json:"resource_url"`
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
	URI    string `json:"uri"`
	URI150 string `json:"uri_150"`
}

//--------------------------------------------------- Label ---------------------------------------------------

// Label is one of the main structure from Discogs:
type Label struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Images      []Image      `json:"images,omitempty"`
	ContactInfo string       `json:"contact_info"`
	Profile     string       `json:"profile"`
	DataQuality string       `json:"data_quality"`
	Urls        []string     `json:"urls"`
	ParentLabel *LabelLabel  `json:"parent_label,omitempty"`
	SubLabels   []LabelLabel `json:"sub_labels,omitempty"`
}

// LabelLabel is a sub structure of Label
type LabelLabel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//--------------------------------------------------- Master ---------------------------------------------------

// Master is one of the main structure from Discogs:
type Master struct {
	ID          string          `json:"id"`
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
	ID           string          `json:"id"`
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
	MasterID     string          `json:"master_id"`
	MainRelease  string          `json:"main_release"`
	TrackList    []Track         `json:"track_list,omitempty"`
	Identifiers  []Identifier    `json:"identifiers,omitempty"`
	Videos       []Video         `json:"videos,omitempty"`
	Labels       []ReleaseLabel  `json:"labels,omitempty"`
	Companies    []Company       `json:"companies,omitempty"`
}

// ReleaseArtist is a sub structure of Release
type ReleaseArtist struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Join   string `json:"join"`
	Anv    string `json:"anv"`
	Role   string `json:"role"`
	Tracks string `json:"tracks"`
}

// ReleaseLabel is a sub structure of Release
type ReleaseLabel struct {
	ID       string `json:"id"`
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
