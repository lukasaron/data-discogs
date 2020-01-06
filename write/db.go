package write

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lukasaron/data-discogs/model"
)

type DbWriter struct {
	o   Options
	db  *sql.DB
	err error
}

func NewDbWriter(db *sql.DB, options *Options) Writer {
	w := DbWriter{db: db}

	// add o when available (only the first one)
	if options != nil {
		w.o = *options
	}

	return w
}

func (db DbWriter) Reset() error {
	db.err = nil
	return nil
}

func (db DbWriter) Close() error {
	return db.db.Close()
}

func (db DbWriter) Options() Options {
	return db.o
}

func (db DbWriter) WriteArtist(artist model.Artist) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	db.writeArtist(tx, artist)
	db.writeAliases(tx, artist.Id, artist.Aliases)
	db.writeImages(tx, artist.Id, "", "", "", artist.Images)
	db.writeArtistMembers(tx, artist.Id, artist.Members)

	if db.err != nil {
		_ = tx.Rollback()
		return db.err
	}

	return tx.Commit()
}

func (db DbWriter) WriteArtists(artists []model.Artist) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, a := range artists {
		db.writeArtist(tx, a)
		db.writeAliases(tx, a.Id, a.Aliases)
		db.writeImages(tx, a.Id, "", "", "", a.Images)
		db.writeArtistMembers(tx, a.Id, a.Members)

		if db.err != nil {
			_ = tx.Rollback()
			return db.err
		}
	}

	return tx.Commit()
}

func (db DbWriter) WriteLabel(label model.Label) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	db.writeLabel(tx, label)
	db.writeLabelLabels(tx, label.Id, "false", label.SubLabels)
	db.writeImages(tx, "", label.Id, "", "", label.Images)

	if label.ParentLabel != nil {
		db.writeLabelLabel(tx, label.Id, "true", *label.ParentLabel)
	}

	if db.err != nil {
		_ = tx.Rollback()
		return db.err
	}

	return tx.Commit()
}

func (db DbWriter) WriteLabels(labels []model.Label) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, l := range labels {
		db.writeLabel(tx, l)
		db.writeLabelLabels(tx, l.Id, "false", l.SubLabels)
		db.writeImages(tx, "", l.Id, "", "", l.Images)

		if l.ParentLabel != nil {
			db.writeLabelLabel(tx, l.Id, "true", *l.ParentLabel)
		}

		if db.err != nil {
			_ = tx.Rollback()
			return db.err
		}
	}

	return tx.Commit()
}

func (db DbWriter) WriteMaster(master model.Master) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	db.writeMaster(tx, master)
	db.writeImages(tx, "", "", master.Id, "", master.Images)
	db.writeReleaseArtists(tx, master.Id, "", "false", master.Artists)
	db.writeVideos(tx, master.Id, "", master.Videos)

	if db.err != nil {
		_ = tx.Rollback()
		return db.err
	}

	return tx.Commit()
}

func (db DbWriter) WriteMasters(masters []model.Master) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, m := range masters {
		db.writeMaster(tx, m)
		db.writeImages(tx, "", "", m.Id, "", m.Images)
		db.writeReleaseArtists(tx, m.Id, "", "false", m.Artists)
		db.writeVideos(tx, m.Id, "", m.Videos)

		if db.err != nil {
			_ = tx.Rollback()
			return db.err
		}
	}

	return tx.Commit()
}

func (db DbWriter) WriteRelease(release model.Release) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	db.writeRelease(tx, release)
	db.writeImages(tx, "", "", "", release.Id, release.Images)
	db.writeReleaseArtists(tx, "", release.Id, "false", release.Artists)
	db.writeReleaseArtists(tx, "", release.Id, "true", release.ExtraArtists)
	db.writeFormats(tx, release.Id, release.Formats)
	db.writeTrackList(tx, release.Id, release.TrackList)
	db.writeIdentifiers(tx, release.Id, release.Identifiers)
	db.writeVideos(tx, "", release.Id, release.Videos)
	db.writeReleaseLabels(tx, release.Id, release.Labels)
	db.writeCompanies(tx, release.Id, release.Companies)
	if db.err != nil {
		_ = tx.Rollback()
		return db.err
	}

	return tx.Commit()
}

func (db DbWriter) WriteReleases(releases []model.Release) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, r := range releases {
		db.writeRelease(tx, r)
		db.writeImages(tx, "", "", "", r.Id, r.Images)
		db.writeReleaseArtists(tx, "", r.Id, "false", r.Artists)
		db.writeReleaseArtists(tx, "", r.Id, "true", r.ExtraArtists)
		db.writeFormats(tx, r.Id, r.Formats)
		db.writeTrackList(tx, r.Id, r.TrackList)
		db.writeIdentifiers(tx, r.Id, r.Identifiers)
		db.writeVideos(tx, "", r.Id, r.Videos)
		db.writeReleaseLabels(tx, r.Id, r.Labels)
		db.writeCompanies(tx, r.Id, r.Companies)

		if db.err != nil {
			_ = tx.Rollback()
			return db.err
		}
	}

	return tx.Commit()
}

func (db DbWriter) writeLabel(tx *sql.Tx, l model.Label) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO labels (label_id, name, contact_info, profile, data_quality, urls) VALUES ('%s', '%s', '%s', '%s', '%s', ARRAY[%s])",
		l.Id,
		cleanText(l.Name),
		cleanText(l.ContactInfo),
		cleanText(l.Profile),
		l.DataQuality,
		array(l.Urls))
}

func (db DbWriter) writeLabelLabel(tx *sql.Tx, labelId, parent string, ll model.LabelLabel) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('%s', '%s', '%s', '%s')",
		labelId,
		ll.Id,
		cleanText(ll.Name),
		parent)
}

func (db DbWriter) writeLabelLabels(tx *sql.Tx, labelId, parent string, lls []model.LabelLabel) {
	if db.err != nil {
		return
	}

	for _, ll := range lls {
		db.writeLabelLabel(tx, labelId, parent, ll)
		if db.err != nil {
			return
		}
	}
}

func (db DbWriter) writeMaster(tx *sql.Tx, m model.Master) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO masters (master_id, main_release, genres, styles, year, title, data_quality) VALUES ('%s', '%s', ARRAY[%s], ARRAY[%s], '%s', '%s', '%s')",
		m.Id,
		m.MainRelease,
		array(m.Genres),
		array(m.Styles),
		m.Year,
		cleanText(m.Title),
		m.DataQuality)
}

func (db DbWriter) writeRelease(tx *sql.Tx, r model.Release) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO releases (release_id, status, title, genres, styles, country, released, notes, data_quality, master_id, main_release) VALUES ('%s', '%s', '%s', ARRAY[%s], ARRAY[%s], '%s', '%s', '%s', '%s', '%s', '%s')",
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
		r.MainRelease)
}

func (db DbWriter) writeCompany(tx *sql.Tx, releaseId string, c model.Company) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO release_companies (release_id, release_company_id, name, category, entity_type, entity_type_name, resource_url) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s')",
		releaseId,
		c.Id,
		cleanText(c.Name),
		cleanText(c.Category),
		cleanText(c.EntityType),
		cleanText(c.EntityTypeName),
		cleanText(c.ResourceUrl))
}

func (db DbWriter) writeCompanies(tx *sql.Tx, releaseId string, cs []model.Company) {
	if db.err != nil {
		return
	}

	for _, c := range cs {
		db.writeCompany(tx, releaseId, c)
		if db.err != nil {
			return
		}
	}
}

func (db DbWriter) writeReleaseArtist(tx *sql.Tx, masterId, releaseId, extra string, ra model.ReleaseArtist) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
		masterId,
		releaseId,
		ra.Id,
		cleanText(ra.Name),
		cleanText(extra),
		cleanText(ra.Join),
		cleanText(ra.Anv),
		cleanText(ra.Role),
		cleanText(ra.Tracks))
}

func (db DbWriter) writeReleaseArtists(tx *sql.Tx, masterId, releaseId, extra string, ras []model.ReleaseArtist) {
	if db.err != nil {
		return
	}

	for _, ra := range ras {
		db.writeReleaseArtist(tx, masterId, releaseId, extra, ra)
		if db.err != nil {
			return
		}
	}
}

func (db DbWriter) writeFormat(tx *sql.Tx, releaseId string, f model.Format) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO release_formats (release_id, name, quantity, text, descriptions) VALUES ('%s', '%s', '%s', '%s', ARRAY[%s])",
		releaseId,
		cleanText(f.Name),
		f.Quantity,
		cleanText(f.Text),
		array(f.Descriptions))
}

func (db DbWriter) writeFormats(tx *sql.Tx, releaseId string, fs []model.Format) {
	if db.err != nil {
		return
	}

	for _, f := range fs {
		db.writeFormat(tx, releaseId, f)
		if db.err != nil {
			return
		}
	}
}

func (db DbWriter) writeTrack(tx *sql.Tx, releaseId string, t model.Track) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO release_tracks (release_id, position, title, duration) VALUES ('%s', '%s', '%s', '%s')",
		releaseId,
		cleanText(t.Position),
		cleanText(t.Title),
		cleanText(t.Duration))
}

func (db DbWriter) writeTrackList(tx *sql.Tx, releaseId string, tl []model.Track) {
	if db.err != nil {
		return
	}

	for _, t := range tl {
		db.writeTrack(tx, releaseId, t)
		if db.err != nil {
			return
		}
	}
}

func (db DbWriter) writeIdentifier(tx *sql.Tx, releaseId string, i model.Identifier) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO release_identifiers (release_id, description, type, value) VALUES ('%s', '%s', '%s', '%s')",
		releaseId,
		cleanText(i.Description),
		cleanText(i.Type),
		cleanText(i.Value))
}

func (db DbWriter) writeIdentifiers(tx *sql.Tx, releaseId string, is []model.Identifier) {
	if db.err != nil {
		return
	}

	for _, i := range is {
		db.writeIdentifier(tx, releaseId, i)
		if db.err != nil {
			return
		}
	}
}

func (db DbWriter) writeReleaseLabel(tx *sql.Tx, releaseId string, rl model.ReleaseLabel) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO release_labels (release_id, release_label_id, name, category) VALUES ('%s', '%s', '%s', '%s')",
		releaseId,
		rl.Id,
		cleanText(rl.Name),
		cleanText(rl.Category))
}

func (db DbWriter) writeReleaseLabels(tx *sql.Tx, releaseId string, rls []model.ReleaseLabel) {
	if db.err != nil {
		return
	}

	for _, rl := range rls {
		db.writeReleaseLabel(tx, releaseId, rl)
		if db.err != nil {
			return
		}
	}
}

func (db DbWriter) writeAlias(tx *sql.Tx, artistId string, a model.Alias) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO artist_aliases (artist_id, alias_id, name) VALUES ('%s', '%s', '%s')",
		artistId,
		a.Id,
		cleanText(a.Name))
}

func (db DbWriter) writeAliases(tx *sql.Tx, artistId string, as []model.Alias) {
	if db.err != nil {
		return
	}

	for _, a := range as {
		db.writeAlias(tx, artistId, a)
		if db.err != nil {
			return
		}
	}
}

func (db DbWriter) writeImage(tx *sql.Tx, artistId, labelId, masterId, releaseId string, img model.Image) {
	if db.err == nil && !db.o.ExcludeImages {
		db.writeTransaction(
			tx,
			"INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
			artistId,
			labelId,
			masterId,
			releaseId,
			img.Height,
			img.Width,
			img.Type,
			img.Uri,
			img.Uri150)
	}
}

func (db DbWriter) writeImages(tx *sql.Tx, artistId, labelId, masterId, releaseId string, imgs []model.Image) {
	if db.err == nil && !db.o.ExcludeImages {
		for _, img := range imgs {
			db.writeImage(tx, artistId, labelId, masterId, releaseId, img)
			if db.err != nil {
				return
			}
		}
	}
}

func (db DbWriter) writeVideo(tx *sql.Tx, masterId, releaseId string, v model.Video) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO videos (master_id, release_id, duration, embed, src, title, description) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s')",
		masterId,
		releaseId,
		cleanText(v.Duration),
		v.Embed,
		cleanText(v.Src),
		cleanText(v.Title),
		cleanText(v.Description))
}

func (db DbWriter) writeVideos(tx *sql.Tx, masterId, releaseId string, vs []model.Video) {
	if db.err != nil {
		return
	}

	for _, v := range vs {
		db.writeVideo(tx, masterId, releaseId, v)
		if db.err != nil {
			return
		}
	}
}

func (db DbWriter) writeArtist(tx *sql.Tx, a model.Artist) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO artists (artist_id, name, real_name, profile, data_quality, name_variations, urls) VALUES ('%s', '%s', '%s', '%s', '%s', ARRAY[%s], ARRAY[%s])",
		a.Id,
		cleanText(a.Name),
		cleanText(a.RealName),
		cleanText(a.Profile),
		a.DataQuality,
		array(a.NameVariations),
		array(a.Urls))
}

func (db DbWriter) writeArtistMember(tx *sql.Tx, artistId string, m model.Member) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO artist_members (artist_id, member_id, name) VALUES ('%s', '%s', '%s')",
		artistId,
		m.Id,
		cleanText(m.Name))
}

func (db DbWriter) writeArtistMembers(tx *sql.Tx, artistId string, ms []model.Member) {
	if db.err != nil {
		return
	}

	for _, m := range ms {
		db.writeArtistMember(tx, artistId, m)
		if db.err != nil {
			return
		}
	}
}

func (db DbWriter) writeTransaction(tx *sql.Tx, query string, values ...interface{}) {
	if db.err != nil {
		return
	}

	_, db.err = tx.Exec(fmt.Sprintf(query, values...))
}
