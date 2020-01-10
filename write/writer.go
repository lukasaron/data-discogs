// Copyright (c) 2020 Lukas Aron. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package write integrates a few writer that could be used for storage the decoded Discogs data.
package write

import (
	"github.com/lukasaron/data-discogs/model"
)

/*
The Writer interface specify all necessary methods for writing Disocgs data that could be useful during processing Discogs dump.
*/
type Writer interface {
	WriteArtist(artist model.Artist) error
	WriteArtists(artists []model.Artist) error
	WriteLabel(label model.Label) error
	WriteLabels(labels []model.Label) error
	WriteMaster(master model.Master) error
	WriteMasters(masters []model.Master) error
	WriteRelease(release model.Release) error
	WriteReleases(releases []model.Release) error
	Options() Options
	Close() error
}

// Options related to writing settings. At this stage only one option is available - Exclude images. This specific option
// is in connection to the Discogs dump data and their politics to provide data without images.
// However, provided data dumps still contains XML tags with property values which are mostly empty.
type Options struct {
	ExcludeImages bool
}
