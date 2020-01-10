// Copyright (c) 2020 Lukas Aron. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package write

import "github.com/lukasaron/data-discogs/model"

// ------------------------------------------------------- DATA -------------------------------------------------------

var artists = []model.Artist{
	{
		ID:          "2",
		Name:        "Mr. James Barth & A.D.",
		RealName:    "Cari Lekebusch & Alexi Delano",
		DataQuality: "Correct",
		NameVariations: []string{"Mr Barth & A.D.", "MR JAMES BARTH & A. D.",
			"Mr. Barth & A.D.", "Mr. James Barth & A. D."},
		Aliases: []model.Alias{
			{
				ID:   "2470",
				Name: "Puente Latino",
			},
			{
				ID:   "19536",
				Name: "Yakari & Delano",
			},
			{
				ID:   "103709",
				Name: "Crushed Insect & The Sick Puppy",
			},
			{
				ID:   "384581",
				Name: "ADCL",
			},
			{
				ID:   "1779857",
				Name: "Alexi Delano & Cari Lekebusch",
			},
		},
		Members: []model.Member{
			{
				ID:   "26",
				Name: "Alexi Delano",
			},
			{
				ID:   "27",
				Name: "Cari Lekebusch",
			},
		},
	},
}

var labels = []model.Label{
	{
		ID:   "1",
		Name: "Planet E",
		Images: []model.Image{
			{
				Height: "24",
				Width:  "132",
				Type:   "primary",
			},
			{
				Height: "126",
				Width:  "587",
				Type:   "secondary",
			},
			{
				Height: "196",
				Width:  "600",
				Type:   "secondary",
			},
			{
				Height: "121",
				Width:  "275",
				Type:   "secondary",
			},
			{
				Height: "720",
				Width:  "382",
				Type:   "secondary",
			},
			{
				Height: "398",
				Width:  "500",
				Type:   "secondary",
			},
			{
				Height: "189",
				Width:  "600",
				Type:   "secondary",
			},
		},
		ContactInfo: "Planet E Communications",
		Profile:     "[a=Carl Craig]'s classic techno label founded in 1991.",
		DataQuality: "Correct",
		Urls: []string{"http://planet-e.net", "http://planetecommunications.bandcamp.com",
			"http://www.facebook.com/planetedetroit", "http://www.flickr.com/photos/planetedetroit",
			"http://plus.google.com/100841702106447505236", "http://www.instagram.com/carlcraignet",
			"http://myspace.com/planetecom", "http://myspace.com/planetedetroit",
			"http://soundcloud.com/planetedetroit", "http://twitter.com/planetedetroit", "http://vimeo.com/user1265384",
			"http://en.wikipedia.org/wiki/Planet_E_Communications", "http://www.youtube.com/user/planetedetroit"},
		SubLabels: []model.LabelLabel{
			{
				ID:   "86537",
				Name: "Antidote (4)",
			},
			{
				ID:   "41841",
				Name: "Community Projects",
			},
			{
				ID:   "153760",
				Name: "Guilty Pleasures",
			},
			{
				ID:   "31405",
				Name: "I Ner Zon Sounds",
			},
			{
				ID:   "277579",
				Name: "Planet E Communications",
			},
			{
				ID:   "294738",
				Name: "Planet E Communications, Inc.",
			},
			{
				ID:   "1560615",
				Name: "Planet E Productions",
			},
			{
				ID:   "488315",
				Name: "TWPENTY",
			},
		},
	},
}

var masters = []model.Master{
	{
		ID:          "18512",
		MainRelease: "33699",
		Images: []model.Image{
			{
				Height: "150",
				Width:  "150",
				Type:   "primary",
			},
			{
				Height: "592",
				Width:  "600",
				Type:   "secondary",
			},
			{
				Height: "592",
				Width:  "600",
				Type:   "secondary",
			},
		},
		Artists: []model.ReleaseArtist{
			{
				ID:   "212070",
				Name: "Samuel L Session",
			},
		},
		Genres:      []string{"Electronic"},
		Styles:      []string{"Tribal", "Techno"},
		Year:        "2002",
		Title:       "Psyche EP",
		DataQuality: "Correct",
		Videos: []model.Video{
			{
				Duration:    "118",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=QYf4j0Pd2FU",
				Title:       "Samuel L. Session - Arrival",
				Description: "Samuel L. Session - Arrival",
			},
			{
				Duration:    "376",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=c_AfLqTdncI",
				Title:       "Samuel L. Session - Psyche Part 1",
				Description: "Samuel L. Session - Psyche Part 1",
			},
			{
				Duration:    "419",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=0nxvR8Zl9wY",
				Title:       "Samuel L. Session - Psyche Part 2",
				Description: "Samuel L. Session - Psyche Part 2",
			},
		},
	},
}

var releases = []model.Release{
	{
		ID:     "2",
		Status: "Accepted",
		Images: []model.Image{
			{
				Height: "394",
				Width:  "400",
				Type:   "primary",
			},
			{
				Height: "600",
				Width:  "600",
				Type:   "secondary",
			},
			{
				Height: "600",
				Width:  "600",
				Type:   "secondary",
			},
		},
		Artists: []model.ReleaseArtist{
			{
				ID:   "2",
				Name: "Mr. James Barth & A.D.",
			},
		},
		ExtraArtists: []model.ReleaseArtist{
			{
				ID:   "26",
				Name: "Alexi Delano",
				Role: "Producer, Recorded By",
			},
			{
				ID:   "27",
				Name: "Cari Lekebusch",
				Role: "Producer, Recorded By",
			},
			{
				ID:   "26",
				Name: "Alexi Delano",
				Anv:  "A. Delano",
				Role: "Written-By",
			},
			{
				ID:   "27",
				Name: "Cari Lekebusch",
				Anv:  "C. Lekebusch",
				Role: "Written-By",
			},
		},
		Title: "Knockin' Boots Vol 2 Of 2",
		Formats: []model.Format{
			{
				Name:         "Vinyl",
				Quantity:     "1",
				Descriptions: []string{"12\"", "33 ⅓ RPM"},
			},
		},
		Genres:      []string{"Electronic"},
		Styles:      []string{"Broken Beat", "Techno", "Tech House"},
		Country:     "Sweden",
		Released:    "1998-06-00",
		Notes:       "All joints recorded in NYC (Dec.97).",
		DataQuality: "Correct",
		MasterID:    "713738",
		MainRelease: "true",
		TrackList: []model.Track{
			{
				Position: "A1",
				Title:    "A Sea Apart",
				Duration: "5:08",
			},
			{
				Position: "A2",
				Title:    "Dutchmaster",
				Duration: "4:21",
			},
			{
				Position: "B1",
				Title:    "Inner City Lullaby",
				Duration: "4:22",
			},
			{
				Position: "B2",
				Title:    "Yeah Kid!",
				Duration: "4:46",
			},
		},
		Identifiers: []model.Identifier{
			{
				Description: "Side A Runout Etching",
				Type:        "Matrix / Runout",
				Value:       "MPO SK026-A -J.T.S.-",
			},
			{
				Description: "Side B Runout Etching",
				Type:        "Matrix / Runout",
				Value:       "MPO SK026-B -J.T.S.-",
			},
		},
		Videos: []model.Video{
			{
				Duration:    "310",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=MIgQNVhYILA",
				Title:       "Mr. James Barth & A.D. - A Sea Apart",
				Description: "Mr. James Barth & A.D. - A Sea Apart",
			},
			{
				Duration:    "265",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=LgLchSRehhc",
				Title:       "Mr. James Barth & A.D. - Dutchmaster",
				Description: "Mr. James Barth & A.D. - Dutchmaster",
			},
			{
				Duration:    "260",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=iaqHaULlqqg",
				Title:       "Mr. James Barth & A.D. - Inner City Lullaby",
				Description: "Mr. James Barth & A.D. - Inner City Lullaby",
			},
			{
				Duration:    "290",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=x_Os7b-iWKs",
				Title:       "Mr. James Barth & A.D. - Yeah Kid!",
				Description: "Mr. James Barth & A.D. - Yeah Kid!",
			},
		},
		Labels: []model.ReleaseLabel{
			{
				ID:       "5",
				Name:     "Svek",
				Category: "SK 026",
			},
			{
				ID:       "5",
				Name:     "Svek",
				Category: "SK026",
			},
		},
		Companies: []model.Company{
			{
				ID:             "266169",
				Name:           "JTS Studios",
				EntityType:     "29",
				EntityTypeName: "Mastered At",
				ResourceURL:    "https://api.discogs.com/labels/266169",
			},
			{
				ID:             "56025",
				Name:           "MPO",
				EntityType:     "17",
				EntityTypeName: "Pressed By",
				ResourceURL:    "https://api.discogs.com/labels/56025",
			},
		},
	},
}
