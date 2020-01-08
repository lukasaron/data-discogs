package write

import (
	"bytes"
	"fmt"
	"github.com/lukasaron/data-discogs/model"
	"io"
	"strings"
)

type SQLWriter struct {
	o   Options
	w   io.Writer
	b   *bytes.Buffer
	err error
}

// Creates a new Writer based on the provided output writer (for instance a file).
// Options with ExcludeImages can be set when we don't want images as part of the final solution.
// When this is not the case and we want images in the result SQL commands the Option has to be passed as a
// second argument.
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

// Writes an artist as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteArtist(artist model.Artist) error {
	s.beginTransaction()

	s.writeArtist(artist)
	s.writeImages(artist.Id, "", "", "", artist.Images)
	s.writeAliases(artist.Id, artist.Aliases)
	s.writeArtistMembers(artist.Id, artist.Members)

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// Writes a slice of artists as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteArtists(artists []model.Artist) error {
	s.beginTransaction()

	for _, a := range artists {
		s.writeArtist(a)
		s.writeImages(a.Id, "", "", "", a.Images)
		s.writeAliases(a.Id, a.Aliases)
		s.writeArtistMembers(a.Id, a.Members)

		if s.err != nil {
			return s.err
		}
	}

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// Writes a label as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteLabel(label model.Label) error {
	s.beginTransaction()

	s.writeLabel(label)
	s.writeLabelLabels(label.Id, "false", label.SubLabels)
	s.writeImages("", label.Id, "", "", label.Images)
	if label.ParentLabel != nil {
		s.writeLabelLabel(label.Id, "true", *label.ParentLabel)
	}

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// Writes a slice of labels as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteLabels(labels []model.Label) error {
	s.beginTransaction()

	for _, l := range labels {
		s.writeLabel(l)
		s.writeLabelLabels(l.Id, "false", l.SubLabels)
		s.writeImages("", l.Id, "", "", l.Images)
		if l.ParentLabel != nil {
			s.writeLabelLabel(l.Id, "true", *l.ParentLabel)
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

// Writes a master as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteMaster(master model.Master) (err error) {
	s.beginTransaction()

	s.writeMaster(master)
	s.writeReleaseArtists(master.Id, "", "false", master.Artists)
	s.writeImages("", "", master.Id, "", master.Images)
	s.writeVideos(master.Id, "", master.Videos)

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// Writes a slice of masters as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteMasters(masters []model.Master) error {
	s.beginTransaction()

	for _, m := range masters {
		s.writeMaster(m)
		s.writeReleaseArtists(m.Id, "", "false", m.Artists)
		s.writeImages("", "", m.Id, "", m.Images)
		s.writeVideos(m.Id, "", m.Videos)
		if s.err != nil {
			return s.err
		}
	}

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// Writes a release as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteRelease(release model.Release) error {
	s.beginTransaction()

	s.writeRelease(release)
	s.writeImages("", "", "", release.Id, release.Images)
	s.writeReleaseArtists("", release.Id, "false", release.Artists)
	s.writeReleaseArtists("", release.Id, "true", release.ExtraArtists)
	s.writeFormats(release.Id, release.Formats)
	s.writeTrackList(release.Id, release.TrackList)
	s.writeIdentifiers(release.Id, release.Identifiers)
	s.writeReleaseLabels(release.Id, release.Labels)
	s.writeCompanies(release.Id, release.Companies)
	s.writeVideos("", release.Id, release.Videos)

	s.commitTransaction()

	s.flush()
	s.clean()

	return s.err
}

// Writes a slice of releases as a set of SQL insert commands into the SQL output.
func (s SQLWriter) WriteReleases(releases []model.Release) error {
	s.beginTransaction()

	for _, r := range releases {
		s.writeRelease(r)
		s.writeImages("", "", "", r.Id, r.Images)
		s.writeReleaseArtists("", r.Id, "false", r.Artists)
		s.writeReleaseArtists("", r.Id, "true", r.ExtraArtists)
		s.writeFormats(r.Id, r.Formats)
		s.writeTrackList(r.Id, r.TrackList)
		s.writeIdentifiers(r.Id, r.Identifiers)
		s.writeReleaseLabels(r.Id, r.Labels)
		s.writeCompanies(r.Id, r.Companies)
		s.writeVideos("", r.Id, r.Videos)
		if s.err != nil {
			return s.err
		}
	}

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

// Removes the inner state error and also reset the inner buffer.
func (s SQLWriter) Reset() error {
	s.err = nil
	s.clean()
	return nil
}

// No behavior implemented.
func (s SQLWriter) Close() error {
	return nil
}

// Returns the current options. It could be useful to get the default options.
func (s SQLWriter) Options() Options {
	return s.o
}

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
			a.Id,
			cleanText(a.Name),
			cleanText(a.RealName),
			cleanText(a.Profile),
			a.DataQuality,
			array(a.NameVariations),
			array(a.Urls)),
	)
}

func (s SQLWriter) writeImage(artistId, labelId, masterId, releaseId string, img model.Image) {
	if !s.o.ExcludeImages {
		_, s.err = s.b.WriteString(
			fmt.Sprintf("INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
				artistId,
				labelId,
				masterId,
				releaseId,
				img.Height,
				img.Width,
				img.Type,
				img.Uri,
				img.Uri150,
			),
		)
	}
}

func (s SQLWriter) writeImages(artistId, labelId, masterId, releaseId string, imgs []model.Image) {
	if s.err == nil && !s.o.ExcludeImages {
		for _, img := range imgs {
			s.writeImage(artistId, labelId, masterId, releaseId, img)
			if s.err != nil {
				return
			}
		}
	}
}

func (s SQLWriter) writeAlias(artistId string, a model.Alias) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(
		fmt.Sprintf("INSERT INTO artist_aliases (artist_id, alias_id, name) VALUES ('%s', '%s', '%s');\n",
			artistId,
			a.Id,
			cleanText(a.Name)),
	)
}

func (s SQLWriter) writeAliases(artistId string, as []model.Alias) {
	if s.err != nil {
		return
	}

	for _, a := range as {
		s.writeAlias(artistId, a)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeArtistMember(artistId string, m model.Member) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO artist_members (artist_id, member_id, name) VALUES ('%s', '%s', '%s');\n",
		artistId,
		m.Id,
		cleanText(m.Name)),
	)
}

func (s SQLWriter) writeArtistMembers(artistId string, ms []model.Member) {
	if s.err != nil {
		return
	}

	for _, m := range ms {
		s.writeArtistMember(artistId, m)
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
			l.Id,
			cleanText(l.Name),
			cleanText(l.ContactInfo),
			cleanText(l.Profile),
			l.DataQuality,
			array(l.Urls),
		),
	)
}

func (s SQLWriter) writeLabelLabel(labelId, parent string, l model.LabelLabel) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(
		fmt.Sprintf("INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('%s', '%s', '%s', '%s');\n",
			labelId,
			l.Id,
			cleanText(l.Name),
			parent,
		),
	)
}

func (s SQLWriter) writeLabelLabels(labelId, parent string, lls []model.LabelLabel) {
	if s.err != nil {
		return
	}

	for _, ll := range lls {
		s.writeLabelLabel(labelId, parent, ll)
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
		m.Id,
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
		r.Id,
		cleanText(r.Status),
		cleanText(r.Title),
		array(r.Genres),
		array(r.Styles),
		cleanText(r.Country),
		r.Released,
		cleanText(r.Notes),
		r.DataQuality,
		r.MasterId,
		r.MainRelease),
	)
}

func (s SQLWriter) writeCompany(releaseId string, c model.Company) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO release_companies (release_id, release_company_id, name, category, entity_type, entity_type_name, resource_url) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
		releaseId,
		c.Id,
		cleanText(c.Name),
		cleanText(c.Category),
		cleanText(c.EntityType),
		cleanText(c.EntityTypeName),
		cleanText(c.ResourceUrl)),
	)
}

func (s SQLWriter) writeCompanies(releaseId string, cs []model.Company) {
	if s.err != nil {
		return
	}

	for _, c := range cs {
		s.writeCompany(releaseId, c)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeReleaseLabel(releaseId string, rl model.ReleaseLabel) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO release_labels (release_id, release_label_id, name, category) VALUES ('%s', '%s', '%s', '%s');\n",
		releaseId,
		rl.Id,
		cleanText(rl.Name),
		cleanText(rl.Category)),
	)
}

func (s SQLWriter) writeReleaseLabels(releaseId string, rls []model.ReleaseLabel) {
	if s.err != nil {
		return
	}

	for _, rl := range rls {
		s.writeReleaseLabel(releaseId, rl)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeIdentifier(releaseId string, i model.Identifier) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO release_identifiers (release_id, description, type, value) VALUES ('%s', '%s', '%s', '%s');\n",
		releaseId,
		cleanText(i.Description),
		cleanText(i.Type),
		cleanText(i.Value)),
	)
}

func (s SQLWriter) writeIdentifiers(releaseId string, is []model.Identifier) {
	if s.err != nil {
		return
	}

	for _, i := range is {
		s.writeIdentifier(releaseId, i)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeTrack(releaseId string, t model.Track) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO release_tracks (release_id, position, title, duration) VALUES ('%s', '%s', '%s', '%s');\n",
		releaseId,
		cleanText(t.Position),
		cleanText(t.Title),
		cleanText(t.Duration)),
	)
}

func (s SQLWriter) writeTrackList(releaseId string, tl []model.Track) {
	if s.err != nil {
		return
	}

	for _, t := range tl {
		s.writeTrack(releaseId, t)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeFormat(releaseId string, f model.Format) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO release_formats (release_id, name, quantity, text, descriptions) VALUES ('%s', '%s', '%s', '%s', ARRAY[%s]);\n",
		releaseId,
		cleanText(f.Name),
		f.Quantity,
		cleanText(f.Text),
		array(f.Descriptions)),
	)
}

func (s SQLWriter) writeFormats(releaseId string, fs []model.Format) {
	if s.err != nil {
		return
	}

	for _, f := range fs {
		s.writeFormat(releaseId, f)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeReleaseArtist(masterId, releaseId, extra string, ra model.ReleaseArtist) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
		masterId,
		releaseId,
		ra.Id,
		cleanText(ra.Name),
		extra,
		cleanText(ra.Join),
		cleanText(ra.Anv),
		cleanText(ra.Role),
		cleanText(ra.Tracks)),
	)
}

func (s SQLWriter) writeReleaseArtists(masterId, releaseId, extra string, ras []model.ReleaseArtist) {
	if s.err != nil {
		return
	}

	for _, ra := range ras {
		s.writeReleaseArtist(masterId, releaseId, extra, ra)
		if s.err != nil {
			return
		}
	}
}

func (s SQLWriter) writeVideo(masterId, releaseId string, v model.Video) {
	if s.err != nil {
		return
	}

	_, s.err = s.b.WriteString(fmt.Sprintf("INSERT INTO videos (master_id, release_id, duration, embed, src, title, description) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
		masterId,
		releaseId,
		cleanText(v.Duration),
		v.Embed,
		cleanText(v.Src),
		cleanText(v.Title),
		cleanText(v.Description)),
	)
}

func (s SQLWriter) writeVideos(masterId, releaseId string, vs []model.Video) {
	if s.err != nil {
		return
	}

	for _, v := range vs {
		s.writeVideo(masterId, releaseId, v)
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

//----------------------------------------------- HELPER FUNCTIONS -----------------------------------------------

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
