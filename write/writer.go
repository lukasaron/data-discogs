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
