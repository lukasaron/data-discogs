package writer

import (
	"bytes"
	"fmt"
	"github.com/Twyer/discogs-parser/model"
	"os"
	"strings"
)

type SqlWriter struct {
	option Options
	f      *os.File
	buffer bytes.Buffer
	err    error
}

func NewSqlWriter(fileName string, options ...Options) Writer {
	s := SqlWriter{buffer: bytes.Buffer{}}

	s.f, s.err = os.Create(fileName)

	if options != nil && len(options) > 0 {
		s.option = options[0]
	}

	return s
}

func (s SqlWriter) WriteArtist(artist model.Artist) error {
	s.beginTransaction()

	s.writeArtist(artist)
	s.writeImages(artist.Id, "", "", "", artist.Images)
	s.writeAliases(artist.Id, artist.Aliases)

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

func (s SqlWriter) WriteArtists(artists []model.Artist) error {
	s.beginTransaction()

	for _, a := range artists {
		s.writeArtist(a)
		s.writeImages(a.Id, "", "", "", a.Images)
		s.writeAliases(a.Id, a.Aliases)
		if s.err != nil {
			return s.err
		}
	}

	s.commitTransaction()
	s.flush()
	s.clean()

	return s.err
}

func (s SqlWriter) WriteLabel(label model.Label) error {
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

func (s SqlWriter) WriteLabels(labels []model.Label) error {
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

func (s SqlWriter) WriteMaster(master model.Master) (err error) {
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

func (s SqlWriter) WriteMasters(masters []model.Master) error {
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

func (s SqlWriter) WriteRelease(release model.Release) error {
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
func (s SqlWriter) WriteReleases(releases []model.Release) error {
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

func (s SqlWriter) Reset() error {
	s.err = nil
	s.clean()
	return nil
}

func (s SqlWriter) Close() error {
	return s.f.Close()
}

func (s SqlWriter) beginTransaction() {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString("BEGIN;\n")
}

func (s SqlWriter) commitTransaction() {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString("COMMIT;\n")
}

func (s SqlWriter) writeArtist(a model.Artist) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(
		fmt.Sprintf("INSERT INTO artists (artist_id, name, real_name, profile, data_quality, name_variations, urls) VALUES ('%s', '%s', '%s', '%s', '%s', ARRAY[%s], ARRAY[%s]);\n",
			a.Id,
			s.cleanText(a.Name),
			s.cleanText(a.RealName),
			s.cleanText(a.Profile),
			a.DataQuality,
			s.array(a.NameVariations),
			s.array(a.Urls)),
	)
}

func (s SqlWriter) writeImage(artistId, labelId, masterId, releaseId string, img model.Image) {
	if !s.option.ExcludeImages {
		_, s.err = s.buffer.WriteString(
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

func (s SqlWriter) writeImages(artistId, labelId, masterId, releaseId string, imgs []model.Image) {
	if s.err == nil && !s.option.ExcludeImages {
		for _, img := range imgs {
			s.writeImage(artistId, labelId, masterId, releaseId, img)
			if s.err != nil {
				return
			}
		}
	}
}

func (s SqlWriter) writeAlias(artistId string, a model.Alias) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(
		fmt.Sprintf("INSERT INTO artist_aliases (artist_id, alias_id, name) VALUES ('%s', '%s', '%s');\n",
			artistId,
			a.Id,
			s.cleanText(a.Name)),
	)
}

func (s SqlWriter) writeAliases(artistId string, as []model.Alias) {
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

func (s SqlWriter) writeLabel(l model.Label) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(
		fmt.Sprintf("INSERT INTO labels (label_id, name, contact_info, profile, data_quality, urls) VALUES ('%s', '%s', '%s', '%s', '%s', ARRAY[%s]);\n",
			l.Id,
			s.cleanText(l.Name),
			s.cleanText(l.ContactInfo),
			s.cleanText(l.Profile),
			l.DataQuality,
			s.array(l.Urls),
		),
	)
}

func (s SqlWriter) writeLabelLabel(labelId, parent string, l model.LabelLabel) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(
		fmt.Sprintf("INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('%s', '%s', '%s', '%s');\n",
			labelId,
			l.Id,
			s.cleanText(l.Name),
			parent,
		),
	)
}

func (s SqlWriter) writeLabelLabels(labelId, parent string, lls []model.LabelLabel) {
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

func (s SqlWriter) writeMaster(m model.Master) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(fmt.Sprintf("INSERT INTO masters (master_id, main_release, genres, styles, year, title, data_quality) VALUES ('%s', '%s', ARRAY[%s], ARRAY[%s], '%s', '%s', '%s');\n",
		m.Id,
		m.MainRelease,
		s.array(m.Genres),
		s.array(m.Styles),
		m.Year,
		s.cleanText(m.Title),
		m.DataQuality),
	)
}

func (s SqlWriter) writeRelease(r model.Release) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(fmt.Sprintf("INSERT INTO releases (release_id, status, title, genres, styles, country, released, notes, data_quality, master_id, main_release) VALUES ('%s', '%s', '%s', ARRAY[%s], ARRAY[%s], '%s', '%s', '%s', '%s', '%s', '%s');\n",
		r.Id,
		s.cleanText(r.Status),
		s.cleanText(r.Title),
		s.array(r.Genres),
		s.array(r.Styles),
		s.cleanText(r.Country),
		r.Released,
		s.cleanText(r.Notes),
		r.DataQuality,
		r.MasterId,
		r.MainRelease),
	)
}

func (s SqlWriter) writeCompany(releaseId string, c model.Company) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(fmt.Sprintf("INSERT INTO release_companies (release_id, release_company_id, name, category, entity_type, entity_type_name, resource_url) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
		releaseId,
		c.Id,
		s.cleanText(c.Name),
		s.cleanText(c.Category),
		s.cleanText(c.EntityType),
		s.cleanText(c.EntityTypeName),
		s.cleanText(c.ResourceUrl)),
	)
}

func (s SqlWriter) writeCompanies(releaseId string, cs []model.Company) {
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

func (s SqlWriter) writeReleaseLabel(releaseId string, rl model.ReleaseLabel) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(fmt.Sprintf("INSERT INTO release_labels (release_id, release_label_id, name, category) VALUES ('%s', '%s', '%s', '%s');\n",
		releaseId,
		rl.Id,
		s.cleanText(rl.Name),
		s.cleanText(rl.Category)),
	)
}

func (s SqlWriter) writeReleaseLabels(releaseId string, rls []model.ReleaseLabel) {
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

func (s SqlWriter) writeIdentifier(releaseId string, i model.Identifier) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(fmt.Sprintf("INSERT INTO release_identifiers (release_id, description, type, value) VALUES ('%s', '%s', '%s', '%s');\n",
		releaseId,
		s.cleanText(i.Description),
		s.cleanText(i.Type),
		s.cleanText(i.Value)),
	)
}

func (s SqlWriter) writeIdentifiers(releaseId string, is []model.Identifier) {
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

func (s SqlWriter) writeTrack(releaseId string, t model.Track) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(fmt.Sprintf("INSERT INTO release_tracks (release_id, position, title, duration) VALUES ('%s', '%s', '%s', '%s');\n",
		releaseId,
		s.cleanText(t.Position),
		s.cleanText(t.Title),
		s.cleanText(t.Duration)),
	)
}

func (s SqlWriter) writeTrackList(releaseId string, tl []model.Track) {
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

func (s SqlWriter) writeFormat(releaseId string, f model.Format) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(fmt.Sprintf("INSERT INTO release_formats (release_id, name, quantity, text, descriptions) VALUES ('%s', '%s', '%s', '%s', ARRAY[%s]);\n",
		releaseId,
		s.cleanText(f.Name),
		f.Quantity,
		s.cleanText(f.Text),
		s.array(f.Descriptions)),
	)
}

func (s SqlWriter) writeFormats(releaseId string, fs []model.Format) {
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

func (s SqlWriter) writeReleaseArtist(masterId, releaseId, extra string, ra model.ReleaseArtist) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(fmt.Sprintf("INSERT INTO release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
		masterId,
		releaseId,
		ra.Id,
		s.cleanText(ra.Name),
		extra,
		s.cleanText(ra.Join),
		s.cleanText(ra.Anv),
		s.cleanText(ra.Role),
		s.cleanText(ra.Tracks)),
	)
}

func (s SqlWriter) writeReleaseArtists(masterId, releaseId, extra string, ras []model.ReleaseArtist) {
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

func (s SqlWriter) writeVideo(masterId, releaseId string, v model.Video) {
	if s.err != nil {
		return
	}

	_, s.err = s.buffer.WriteString(fmt.Sprintf("INSERT INTO videos (master_id, release_id, duration, embed, src, title, description) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
		masterId,
		releaseId,
		s.cleanText(v.Duration),
		v.Embed,
		s.cleanText(v.Src),
		s.cleanText(v.Title),
		s.cleanText(v.Description)),
	)
}

func (s SqlWriter) writeVideos(masterId, releaseId string, vs []model.Video) {
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

func (s SqlWriter) flush() {
	if s.err != nil {
		return
	}

	_, s.err = s.f.Write(s.buffer.Bytes())
}

func (s SqlWriter) clean() {
	s.buffer.Reset()
}

//----------------------------------------------- HELPER FUNCTIONS -----------------------------------------------

func (s SqlWriter) cleanText(str string) string {
	return strings.ReplaceAll(str, "'", "''")
}

func (s SqlWriter) array(str []string) string {
	sb := strings.Builder{}
	sb.WriteString("'")
	sb.WriteString(strings.ReplaceAll(s.cleanText(strings.Join(str, ",")), ",", "','"))
	sb.WriteString("'")
	return sb.String()
}
