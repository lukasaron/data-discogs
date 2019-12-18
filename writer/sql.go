package writer

import (
	"fmt"
	"github.com/Twyer/discogs-parser/model"
	"os"
	"strings"
)

type SqlWriter struct {
	option Options
	f      *os.File
	err    error
}

func NewSqlWriter(fileName string, options ...Options) Writer {
	j := SqlWriter{}

	j.f, j.err = os.Create(fileName)

	if options != nil && len(options) > 0 {
		j.option = options[0]
	}

	return j
}

func (s SqlWriter) WriteArtist(artist model.Artist) (err error) {
	if s.err != nil {
		return s.err
	}

	err = s.beginTransaction()
	if err != nil {
		return err
	}

	err = s.writeArtist(artist)
	if err != nil {
		return err
	}

	return s.commitTransaction()
}

func (s SqlWriter) WriteArtists(artists []model.Artist) (err error) {
	err = s.beginTransaction()
	if err != nil {
		return err
	}

	for _, a := range artists {
		err = s.writeArtist(a)
		if err != nil {
			return err
		}
	}

	return s.commitTransaction()
}

func (s SqlWriter) WriteLabel(label model.Label) (err error) {
	err = s.beginTransaction()
	if err != nil {
		return err
	}

	err = s.writeLabel(label)
	if err != nil {
		return err
	}

	return s.commitTransaction()
}

func (s SqlWriter) WriteLabels(labels []model.Label) (err error) {
	err = s.beginTransaction()
	if err != nil {
		return err
	}

	for _, l := range labels {
		err = s.writeLabel(l)
		if err != nil {
			return err
		}
	}

	return s.commitTransaction()
}

func (s SqlWriter) WriteMaster(master model.Master) (err error) {
	err = s.beginTransaction()
	if err != nil {
		return err
	}

	err = s.writeMaster(master)
	if err != nil {
		return err
	}

	return s.commitTransaction()
}

func (s SqlWriter) WriteMasters(masters []model.Master) (err error) {
	err = s.beginTransaction()
	if err != nil {
		return err
	}

	for _, m := range masters {
		err = s.writeMaster(m)
		if err != nil {
			return err
		}
	}

	return s.commitTransaction()
}

func (s SqlWriter) WriteRelease(release model.Release) (err error) {
	err = s.beginTransaction()
	if err != nil {
		return err
	}

	err = s.writeRelease(release)
	if err != nil {
		return err
	}

	return s.commitTransaction()
}
func (s SqlWriter) WriteReleases(releases []model.Release) (err error) {
	err = s.beginTransaction()
	if err != nil {
		return err
	}

	for _, r := range releases {
		err = s.writeRelease(r)
		if err != nil {
			return err
		}
	}

	return s.commitTransaction()
}

func (s SqlWriter) Close() error {
	return s.f.Close()
}

func (s SqlWriter) beginTransaction() (err error) {
	_, err = s.f.WriteString("BEGIN;\n")
	return err
}

func (s SqlWriter) commitTransaction() (err error) {
	_, err = s.f.WriteString("COMMIT;\n")
	return err
}

func (s SqlWriter) writeArtist(a model.Artist) (err error) {
	_, err = s.f.WriteString(
		fmt.Sprintf("INSERT INTO artists (artist_id, name, real_name, profile, data_quality, name_variations, urls) VALUES ('%s', '%s', '%s', '%s', '%s', ARRAY[%s], ARRAY[%s]);\n",
			a.Id,
			s.cleanText(a.Name),
			s.cleanText(a.RealName),
			s.cleanText(a.Profile),
			a.DataQuality,
			s.array(a.NameVariations),
			s.array(a.Urls)),
	)
	if err != nil {
		return err
	}

	err = s.writeImages(a.Id, "", "", "", a.Images)
	if err != nil {
		return err
	}

	return s.writeAliases(a.Id, a.Aliases)
}

func (s SqlWriter) writeImage(artistId, labelId, masterId, releaseId string, img model.Image) error {
	if !s.option.ExcludeImages {
		_, err := s.f.WriteString(
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
		return err
	}

	return nil
}

func (s SqlWriter) writeImages(artistId, labelId, masterId, releaseId string, imgs []model.Image) (err error) {
	if !s.option.ExcludeImages {
		for _, img := range imgs {
			err = s.writeImage(artistId, labelId, masterId, releaseId, img)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s SqlWriter) writeAlias(artistId string, a model.Alias) (err error) {
	_, err = s.f.WriteString(
		fmt.Sprintf("INSERT INTO artist_aliases (artist_id, alias_id, name) VALUES ('%s', '%s', '%s');\n",
			artistId,
			a.Id,
			s.cleanText(a.Name)),
	)

	return err
}

func (s SqlWriter) writeAliases(artistId string, as []model.Alias) (err error) {
	for _, a := range as {
		err = s.writeAlias(artistId, a)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SqlWriter) writeLabel(l model.Label) (err error) {
	_, err = s.f.WriteString(
		fmt.Sprintf("INSERT INTO labels (label_id, name, contact_info, profile, data_quality, urls) VALUES ('%s', '%s', '%s', '%s', '%s', ARRAY[%s]);\n",
			l.Id,
			s.cleanText(l.Name),
			s.cleanText(l.ContactInfo),
			s.cleanText(l.Profile),
			l.DataQuality,
			s.array(l.Urls),
		),
	)
	if err != nil {
		return err
	}

	if l.ParentLabel != nil {
		err = s.writeLabelLabel(l.Id, "true", *l.ParentLabel)
		if err != nil {
			return err
		}
	}

	err = s.writeLabelLabels(l.Id, "false", l.SubLabels)
	if err != nil {
		return err
	}

	return s.writeImages("", l.Id, "", "", l.Images)
}

func (s SqlWriter) writeLabelLabel(labelId, parent string, l model.LabelLabel) (err error) {
	_, err = s.f.WriteString(
		fmt.Sprintf("INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('%s', '%s', '%s', '%s');\n",
			labelId,
			l.Id,
			s.cleanText(l.Name),
			parent,
		),
	)

	return err
}

func (s SqlWriter) writeLabelLabels(labelId, parent string, lls []model.LabelLabel) (err error) {
	for _, ll := range lls {
		err = s.writeLabelLabel(labelId, parent, ll)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SqlWriter) writeMaster(m model.Master) (err error) {
	_, err = s.f.WriteString(fmt.Sprintf("INSERT INTO masters (master_id, main_release, genres, styles, year, title, data_quality) VALUES ('%s', '%s', ARRAY[%s], ARRAY[%s], '%s', '%s', '%s');\n",
		m.Id,
		m.MainRelease,
		s.array(m.Genres),
		s.array(m.Styles),
		m.Year,
		s.cleanText(m.Title),
		m.DataQuality),
	)

	if err != nil {
		return err
	}

	err = s.writeReleaseArtists(m.Id, "", "false", m.Artists)
	if err != nil {
		return err
	}

	err = s.writeImages("", "", m.Id, "", m.Images)
	if err != nil {
		return err
	}

	return s.writeVideos(m.Id, "", m.Videos)
}

func (s SqlWriter) writeRelease(r model.Release) (err error) {
	_, err = s.f.WriteString(fmt.Sprintf("INSERT INTO releases (release_id, status, title, genres, styles, country, released, notes, data_quality, master_id, main_release) VALUES ('%s', '%s', '%s', ARRAY[%s], ARRAY[%s], '%s', '%s', '%s', '%s', '%s', '%s');\n",
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

	if err != nil {
		return err
	}

	err = s.writeImages("", "", "", r.Id, r.Images)
	if err != nil {
		return err
	}

	err = s.writeReleaseArtists("", r.Id, "false", r.Artists)
	if err != nil {
		return err
	}

	err = s.writeReleaseArtists("", r.Id, "true", r.ExtraArtists)
	if err != nil {
		return err
	}

	err = s.writeFormats(r.Id, r.Formats)
	if err != nil {
		return err
	}

	err = s.writeTrackList(r.Id, r.TrackList)
	if err != nil {
		return err
	}

	err = s.writeIdentifiers(r.Id, r.Identifiers)
	if err != nil {
		return err
	}

	err = s.writeReleaseLabels(r.Id, r.Labels)
	if err != nil {
		return err
	}

	err = s.writeCompanies(r.Id, r.Companies)
	if err != nil {
		return err
	}

	return s.writeVideos("", r.Id, r.Videos)
}

func (s SqlWriter) writeCompany(releaseId string, c model.Company) (err error) {
	_, err = s.f.WriteString(fmt.Sprintf("INSERT INTO release_companies (release_id, release_company_id, name, category, entity_type, entity_type_name, resource_url) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
		releaseId,
		c.Id,
		s.cleanText(c.Name),
		s.cleanText(c.Category),
		s.cleanText(c.EntityType),
		s.cleanText(c.EntityTypeName),
		s.cleanText(c.ResourceUrl)),
	)

	return err
}

func (s SqlWriter) writeCompanies(releaseId string, cs []model.Company) (err error) {
	for _, c := range cs {
		err = s.writeCompany(releaseId, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SqlWriter) writeReleaseLabel(releaseId string, rl model.ReleaseLabel) (err error) {
	_, err = s.f.WriteString(fmt.Sprintf("INSERT INTO release_labels (release_id, release_label_id, name, category) VALUES ('%s', '%s', '%s', '%s');\n",
		releaseId,
		rl.Id,
		s.cleanText(rl.Name),
		s.cleanText(rl.Category)),
	)

	return err
}

func (s SqlWriter) writeReleaseLabels(releaseId string, rls []model.ReleaseLabel) (err error) {
	for _, rl := range rls {
		err = s.writeReleaseLabel(releaseId, rl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SqlWriter) writeIdentifier(releaseId string, i model.Identifier) (err error) {
	_, err = s.f.WriteString(fmt.Sprintf("INSERT INTO release_identifiers (release_id, description, type, value) VALUES ('%s', '%s', '%s', '%s');\n",
		releaseId,
		s.cleanText(i.Description),
		s.cleanText(i.Type),
		s.cleanText(i.Value)),
	)

	return err
}

func (s SqlWriter) writeIdentifiers(releaseId string, is []model.Identifier) (err error) {
	for _, i := range is {
		err = s.writeIdentifier(releaseId, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SqlWriter) writeTrack(releaseId string, t model.Track) (err error) {
	_, err = s.f.WriteString(fmt.Sprintf("INSERT INTO release_tracks (release_id, position, title, duration) VALUES ('%s', '%s', '%s', '%s');\n",
		releaseId,
		s.cleanText(t.Position),
		s.cleanText(t.Title),
		s.cleanText(t.Duration)),
	)

	return err
}

func (s SqlWriter) writeTrackList(releaseId string, tl []model.Track) (err error) {
	for _, t := range tl {
		err = s.writeTrack(releaseId, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SqlWriter) writeFormat(releaseId string, f model.Format) (err error) {
	_, err = s.f.WriteString(fmt.Sprintf("INSERT INTO release_formats (release_id, name, quantity, text, descriptions) VALUES ('%s', '%s', '%s', '%s', ARRAY[%s]);\n",
		releaseId,
		s.cleanText(f.Name),
		f.Quantity,
		s.cleanText(f.Text),
		s.array(f.Descriptions)),
	)

	return err
}

func (s SqlWriter) writeFormats(releaseId string, fs []model.Format) (err error) {
	for _, f := range fs {
		err = s.writeFormat(releaseId, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SqlWriter) writeReleaseArtist(masterId, releaseId, extra string, ra model.ReleaseArtist) (err error) {
	_, err = s.f.WriteString(fmt.Sprintf("INSERT INTO release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
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

	return err
}

func (s SqlWriter) writeReleaseArtists(masterId, releaseId, extra string, ras []model.ReleaseArtist) (err error) {
	for _, ra := range ras {
		err = s.writeReleaseArtist(masterId, releaseId, extra, ra)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SqlWriter) writeVideo(masterId, releaseId string, v model.Video) (err error) {
	_, err = s.f.WriteString(fmt.Sprintf("INSERT INTO videos (master_id, release_id, duration, embed, src, title, description) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
		masterId,
		releaseId,
		s.cleanText(v.Duration),
		v.Embed,
		s.cleanText(v.Src),
		s.cleanText(v.Title),
		s.cleanText(v.Description)),
	)

	return err
}

func (s SqlWriter) writeVideos(masterId, releaseId string, vs []model.Video) (err error) {
	for _, v := range vs {
		err = s.writeVideo(masterId, releaseId, v)
		if err != nil {
			return err
		}
	}

	return nil
}

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
