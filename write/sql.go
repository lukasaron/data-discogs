// Copyright (c) 2020 Lukas Aron. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package write

import (
	"bytes"
	"fmt"
	"github.com/lukasaron/data-discogs/model"
	"io"
	"strings"
)

// SQLWriter is one of few provided writers that implements the Writer interface and provides the ability to save
// decoded data in the format of SQL insert commands.
type SQLWriter struct {
	o   Options
	w   io.Writer
	b   *bytes.Buffer
	err error
}

// NewSQLWriter creates a new Writer instance based on the provided output writer (for instance a file).
// Options with ExcludeImages can be set when we don't want images as part of the final solution.
// When this is not the case and we want images in the result SQL commands the Option can be omitted.
func NewSQLWriter(output io.Writer, options *Options) Writer {

	if options == nil {
		options = &Options{}
	}

	return &SQLWriter{
		b: &bytes.Buffer{},
		o: *options,
		w: output,
	}
}

// WriteArtist function writes an artist as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteArtist(artist model.Artist) error {
	s.beginTransaction()

	s.writeArtist(artist)
	s.writeImages(artist.ID, "", "", "", artist.Images)
	s.writeAliases(artist.ID, artist.Aliases)
	s.writeArtistMembers(artist.ID, artist.Members)

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// WriteArtists function writes a slice of artists as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteArtists(artists []model.Artist) error {
	s.beginTransaction()

	for _, a := range artists {
		s.writeArtist(a)
		s.writeImages(a.ID, "", "", "", a.Images)
		s.writeAliases(a.ID, a.Aliases)
		s.writeArtistMembers(a.ID, a.Members)

		if s.err != nil {
			return s.err
		}
	}

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// WriteLabel function writes a label as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteLabel(label model.Label) error {
	s.beginTransaction()

	s.writeLabel(label)
	s.writeLabelLabels(label.ID, "false", label.SubLabels)
	s.writeImages("", label.ID, "", "", label.Images)
	if label.ParentLabel != nil {
		s.writeLabelLabel(label.ID, "true", *label.ParentLabel)
	}

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// WriteLabels function writes a slice of labels as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteLabels(labels []model.Label) error {
	s.beginTransaction()

	for _, l := range labels {
		s.writeLabel(l)
		s.writeLabelLabels(l.ID, "false", l.SubLabels)
		s.writeImages("", l.ID, "", "", l.Images)
		if l.ParentLabel != nil {
			s.writeLabelLabel(l.ID, "true", *l.ParentLabel)
		}
		if s.err != nil {
			return s.err
		}
	}

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// WriteMaster function writes a master as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteMaster(master model.Master) (err error) {
	s.beginTransaction()

	s.writeMaster(master)
	s.writeReleaseArtists(master.ID, "", "false", master.Artists)
	s.writeImages("", "", master.ID, "", master.Images)
	s.writeVideos(master.ID, "", master.Videos)

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// WriteMasters function writes a slice of masters as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteMasters(masters []model.Master) error {
	s.beginTransaction()

	for _, m := range masters {
		s.writeMaster(m)
		s.writeReleaseArtists(m.ID, "", "false", m.Artists)
		s.writeImages("", "", m.ID, "", m.Images)
		s.writeVideos(m.ID, "", m.Videos)
		if s.err != nil {
			return s.err
		}
	}

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// WriteRelease function writes a release as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteRelease(release model.Release) error {
	s.beginTransaction()

	s.writeRelease(release)
	s.writeImages("", "", "", release.ID, release.Images)
	s.writeReleaseArtists("", release.ID, "false", release.Artists)
	s.writeReleaseArtists("", release.ID, "true", release.ExtraArtists)
	s.writeFormats(release.ID, release.Formats)
	s.writeTrackList(release.ID, release.TrackList)
	s.writeIdentifiers(release.ID, release.Identifiers)
	s.writeReleaseLabels(release.ID, release.Labels)
	s.writeCompanies(release.ID, release.Companies)
	s.writeVideos("", release.ID, release.Videos)

	s.commitTransaction()

	s.flush()
	s.clean()

	return s.err
}

// WriteReleases function writes a slice of releases as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteReleases(releases []model.Release) error {
	s.beginTransaction()

	for _, r := range releases {
		s.writeRelease(r)
		s.writeImages("", "", "", r.ID, r.Images)
		s.writeReleaseArtists("", r.ID, "false", r.Artists)
		s.writeReleaseArtists("", r.ID, "true", r.ExtraArtists)
		s.writeFormats(r.ID, r.Formats)
		s.writeTrackList(r.ID, r.TrackList)
		s.writeIdentifiers(r.ID, r.Identifiers)
		s.writeReleaseLabels(r.ID, r.Labels)
		s.writeCompanies(r.ID, r.Companies)
		s.writeVideos("", r.ID, r.Videos)
		if s.err != nil {
			return s.err
		}
	}

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// Options function returns the current options. It could be useful to get the default options.
func (s SQLWriter) Options() Options {
	return s.o
}

// ----------------------------------------------- UNPUBLISHED FUNCTIONS -----------------------------------------------

func (s SQLWriter) beginTransaction() {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString("BEGIN;\n")
}

func (s SQLWriter) commitTransaction() {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString("COMMIT;\n")
}

func (s SQLWriter) writeArtist(a model.Artist) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(
		fmt.Sprintf("INSERT INTO artists (artist_id, name, real_name, profile, data_quality, name_variations, urls) VALUES ('%s', '%s', '%s', '%s', '%s', ARRAY[%s], ARRAY[%s]);\n",
			a.ID,
			cleanText(a.Name),
			cleanText(a.RealName),
			cleanText(a.Profile),
			a.DataQuality,
			array(a.NameVariations),
			array(a.Urls)),
	)
}

func (s SQLWriter) writeImage(artistID, labelID, masterID, releaseID string, img model.Image) {
	if !s.o.ExcludeImages {
		_, s.err = s.b.WriteString(
			fmt.Sprintf("INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
				artistID,
				labelID,
				masterID,
				releaseID,
				img.Height,
				img.Width,
				img.Type,
				img.URI,
				img.URI150,
			),
		)
	}
}

func (s SQLWriter) writeImages(artistID, labelID, masterID, releaseID string, imgs []model.Image) {
	if s.err == nil && !s.o.ExcludeImages {
		for _, img := range imgs {
			s.writeImage(artistID, labelID, masterID, releaseID, img)
			if s.err != nil {
				return
			}
		}
	}
}

func (s SQLWriter) writeAlias(artistID string, a model.Alias) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(
		fmt.Sprintf("INSERT INTO artist_aliases (artist_id, alias_id, name) VALUES ('%s', '%s', '%s');\n",
			artistID,
			a.ID,
			cleanText(a.Name)),
	)
}

func (s SQLWriter) writeAliases(artistID string, as []model.Alias) {
	if s.err != nil {
		return
	}

	for _, a := range as {
		s.writeAlias(artistID, a)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeArtistMember(artistID string, m model.Member) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO artist_members (artist_id, member_id, name) VALUES ('%s', '%s', '%s');\n",
		artistID,
		m.ID,
		cleanText(m.Name)),
	)
}

func (s SQLWriter) writeArtistMembers(artistID string, ms []model.Member) {
	if s.err != nil {
		return
	}

	for _, m := range ms {
		s.writeArtistMember(artistID, m)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeLabel(l model.Label) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(
		fmt.Sprintf("INSERT INTO labels (label_id, name, contact_info, profile, data_quality, urls) VALUES ('%s', '%s', '%s', '%s', '%s', ARRAY[%s]);\n",
			l.ID,
			cleanText(l.Name),
			cleanText(l.ContactInfo),
			cleanText(l.Profile),
			l.DataQuality,
			array(l.Urls),
		),
	)
}

func (s SQLWriter) writeLabelLabel(labelID, parent string, l model.LabelLabel) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(
		fmt.Sprintf("INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('%s', '%s', '%s', '%s');\n",
			labelID,
			l.ID,
			cleanText(l.Name),
			parent,
		),
	)
}

func (s SQLWriter) writeLabelLabels(labelID, parent string, lls []model.LabelLabel) {
	if s.err != nil {
		return
	}

	for _, ll := range lls {
		s.writeLabelLabel(labelID, parent, ll)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeMaster(m model.Master) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO masters (master_id, main_release, genres, styles, year, title, data_quality) VALUES ('%s', '%s', ARRAY[%s], ARRAY[%s], '%s', '%s', '%s');\n",
		m.ID,
		m.MainRelease,
		array(m.Genres),
		array(m.Styles),
		m.Year,
		cleanText(m.Title),
		m.DataQuality),
	)
}

func (s SQLWriter) writeRelease(r model.Release) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO releases (release_id, status, title, genres, styles, country, released, notes, data_quality, master_id, main_release) VALUES ('%s', '%s', '%s', ARRAY[%s], ARRAY[%s], '%s', '%s', '%s', '%s', '%s', '%s');\n",
		r.ID,
		cleanText(r.Status),
		cleanText(r.Title),
		array(r.Genres),
		array(r.Styles),
		cleanText(r.Country),
		r.Released,
		cleanText(r.Notes),
		r.DataQuality,
		r.MasterID,
		r.MainRelease),
	)
}

func (s SQLWriter) writeCompany(releaseID string, c model.Company) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO release_companies (release_id, release_company_id, name, category, entity_type, entity_type_name, resource_url) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
		releaseID,
		c.ID,
		cleanText(c.Name),
		cleanText(c.Category),
		cleanText(c.EntityType),
		cleanText(c.EntityTypeName),
		cleanText(c.ResourceURL)),
	)
}

func (s SQLWriter) writeCompanies(releaseID string, cs []model.Company) {
	if s.err != nil {
		return
	}

	for _, c := range cs {
		s.writeCompany(releaseID, c)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeReleaseLabel(releaseID string, rl model.ReleaseLabel) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO release_labels (release_id, release_label_id, name, category) VALUES ('%s', '%s', '%s', '%s');\n",
		releaseID,
		rl.ID,
		cleanText(rl.Name),
		cleanText(rl.Category)),
	)
}

func (s SQLWriter) writeReleaseLabels(releaseID string, rls []model.ReleaseLabel) {
	if s.err != nil {
		return
	}

	for _, rl := range rls {
		s.writeReleaseLabel(releaseID, rl)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeIdentifier(releaseID string, i model.Identifier) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO release_identifiers (release_id, description, type, value) VALUES ('%s', '%s', '%s', '%s');\n",
		releaseID,
		cleanText(i.Description),
		cleanText(i.Type),
		cleanText(i.Value)),
	)
}

func (s SQLWriter) writeIdentifiers(releaseID string, is []model.Identifier) {
	if s.err != nil {
		return
	}

	for _, i := range is {
		s.writeIdentifier(releaseID, i)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeTrack(releaseID string, t model.Track) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO release_tracks (release_id, position, title, duration) VALUES ('%s', '%s', '%s', '%s');\n",
		releaseID,
		cleanText(t.Position),
		cleanText(t.Title),
		cleanText(t.Duration)),
	)
}

func (s SQLWriter) writeTrackList(releaseID string, tl []model.Track) {
	if s.err != nil {
		return
	}

	for _, t := range tl {
		s.writeTrack(releaseID, t)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeFormat(releaseID string, f model.Format) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO release_formats (release_id, name, quantity, text, descriptions) VALUES ('%s', '%s', '%s', '%s', ARRAY[%s]);\n",
		releaseID,
		cleanText(f.Name),
		f.Quantity,
		cleanText(f.Text),
		array(f.Descriptions)),
	)
}

func (s SQLWriter) writeFormats(releaseID string, fs []model.Format) {
	if s.err != nil {
		return
	}

	for _, f := range fs {
		s.writeFormat(releaseID, f)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeReleaseArtist(masterID, releaseID, extra string, ra model.ReleaseArtist) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
		masterID,
		releaseID,
		ra.ID,
		cleanText(ra.Name),
		extra,
		cleanText(ra.Join),
		cleanText(ra.Anv),
		cleanText(ra.Role),
		cleanText(ra.Tracks)),
	)
}

func (s SQLWriter) writeReleaseArtists(masterID, releaseID, extra string, ras []model.ReleaseArtist) {
	if s.err != nil {
		return
	}

	for _, ra := range ras {
		s.writeReleaseArtist(masterID, releaseID, extra, ra)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeVideo(masterID, releaseID string, v model.Video) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO videos (master_id, release_id, duration, embed, src, title, description) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
		masterID,
		releaseID,
		cleanText(v.Duration),
		v.Embed,
		cleanText(v.Src),
		cleanText(v.Title),
		cleanText(v.Description)),
	)
}

func (s SQLWriter) writeVideos(masterID, releaseID string, vs []model.Video) {
	if s.err != nil {
		return
	}

	for _, v := range vs {
		s.writeVideo(masterID, releaseID, v)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) flush() {
	if s.err != nil {
		return
	}

	_, s.err = s.w.Write(s.b.Bytes())
}

func (s SQLWriter) clean() {
	s.b.Reset()
}

// ----------------------------------------------- HELPER FUNCTIONS -----------------------------------------------

func cleanText(str string) string {
	return strings.ReplaceAll(str, "'", "''")
}

func array(str []string) string {
	sb := strings.Builder{}
	sb.WriteString("'")
	sb.WriteString(strings.ReplaceAll(cleanText(strings.Join(str, ",")), ",", "','"))
	sb.WriteString("'")
	return sb.String()
}
