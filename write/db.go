package write

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lukasaron/data-discogs/model"
)

// DBWriter is one of few provided writers that implements the Writer interface and provides the ability to save
// decoded data directly into SQL Database.
type DBWriter struct {
	o   Options
	db  *sql.DB
	err error
}

// Creates a new Writer based on the connection to SQL database. Options with ExcludeImages can be set when we don't
// want images as part of the final solution.
// When this is not the case and we want images in the database table the Option has to be passed as a second argument.
func NewDBWriter(db *sql.DB, options *Options) Writer {

	if options == nil {
		options = &Options{}
	}

	return DBWriter{
		db: db,
		o:  *options,
	}
}

// Reset function removes the state error. The reason to have this feature is when a user want's to try recover from
// previous mistake/error.
func (db DBWriter) Reset() error {
	db.err = nil
	return nil
}

// Closes opened database connection
func (db DBWriter) Close() error {
	return db.db.Close()
}

// Getter for options. Can be used to get the default values.
func (db DBWriter) Options() Options {
	return db.o
}

// Writes an artist to the provided database within a transaction
func (db DBWriter) WriteArtist(artist model.Artist) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	db.writeArtist(tx, artist)
	db.writeAliases(tx, artist.ID, artist.Aliases)
	db.writeImages(tx, artist.ID, "", "", "", artist.Images)
	db.writeArtistMembers(tx, artist.ID, artist.Members)

	if db.err != nil {
		_ = tx.Rollback()
		return db.err
	}

	return tx.Commit()
}

// Writes a slice of artists to the provided database within a transaction
func (db DBWriter) WriteArtists(artists []model.Artist) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, a := range artists {
		db.writeArtist(tx, a)
		db.writeAliases(tx, a.ID, a.Aliases)
		db.writeImages(tx, a.ID, "", "", "", a.Images)
		db.writeArtistMembers(tx, a.ID, a.Members)

		if db.err != nil {
			_ = tx.Rollback()
			return db.err
		}
	}

	return tx.Commit()
}

// Writes a label to the provided database within a transaction
func (db DBWriter) WriteLabel(label model.Label) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	db.writeLabel(tx, label)
	db.writeLabelLabels(tx, label.ID, "false", label.SubLabels)
	db.writeImages(tx, "", label.ID, "", "", label.Images)

	if label.ParentLabel != nil {
		db.writeLabelLabel(tx, label.ID, "true", *label.ParentLabel)
	}

	if db.err != nil {
		_ = tx.Rollback()
		return db.err
	}

	return tx.Commit()
}

// Writes a slice of labels to the provided database within a transaction
func (db DBWriter) WriteLabels(labels []model.Label) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, l := range labels {
		db.writeLabel(tx, l)
		db.writeLabelLabels(tx, l.ID, "false", l.SubLabels)
		db.writeImages(tx, "", l.ID, "", "", l.Images)

		if l.ParentLabel != nil {
			db.writeLabelLabel(tx, l.ID, "true", *l.ParentLabel)
		}

		if db.err != nil {
			_ = tx.Rollback()
			return db.err
		}
	}

	return tx.Commit()
}

// Writes a master to the provided database within a transaction
func (db DBWriter) WriteMaster(master model.Master) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	db.writeMaster(tx, master)
	db.writeImages(tx, "", "", master.ID, "", master.Images)
	db.writeReleaseArtists(tx, master.ID, "", "false", master.Artists)
	db.writeVideos(tx, master.ID, "", master.Videos)

	if db.err != nil {
		_ = tx.Rollback()
		return db.err
	}

	return tx.Commit()
}

// Writes a slice of masters to the provided database within a transaction
func (db DBWriter) WriteMasters(masters []model.Master) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, m := range masters {
		db.writeMaster(tx, m)
		db.writeImages(tx, "", "", m.ID, "", m.Images)
		db.writeReleaseArtists(tx, m.ID, "", "false", m.Artists)
		db.writeVideos(tx, m.ID, "", m.Videos)

		if db.err != nil {
			_ = tx.Rollback()
			return db.err
		}
	}

	return tx.Commit()
}

// Writes a release to the provided database within a transaction
func (db DBWriter) WriteRelease(release model.Release) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	db.writeRelease(tx, release)
	db.writeImages(tx, "", "", "", release.ID, release.Images)
	db.writeReleaseArtists(tx, "", release.ID, "false", release.Artists)
	db.writeReleaseArtists(tx, "", release.ID, "true", release.ExtraArtists)
	db.writeFormats(tx, release.ID, release.Formats)
	db.writeTrackList(tx, release.ID, release.TrackList)
	db.writeIdentifiers(tx, release.ID, release.Identifiers)
	db.writeVideos(tx, "", release.ID, release.Videos)
	db.writeReleaseLabels(tx, release.ID, release.Labels)
	db.writeCompanies(tx, release.ID, release.Companies)
	if db.err != nil {
		_ = tx.Rollback()
		return db.err
	}

	return tx.Commit()
}

// Writes a slice of releases to the provided database within a transaction
func (db DBWriter) WriteReleases(releases []model.Release) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, r := range releases {
		db.writeRelease(tx, r)
		db.writeImages(tx, "", "", "", r.ID, r.Images)
		db.writeReleaseArtists(tx, "", r.ID, "false", r.Artists)
		db.writeReleaseArtists(tx, "", r.ID, "true", r.ExtraArtists)
		db.writeFormats(tx, r.ID, r.Formats)
		db.writeTrackList(tx, r.ID, r.TrackList)
		db.writeIdentifiers(tx, r.ID, r.Identifiers)
		db.writeVideos(tx, "", r.ID, r.Videos)
		db.writeReleaseLabels(tx, r.ID, r.Labels)
		db.writeCompanies(tx, r.ID, r.Companies)

		if db.err != nil {
			_ = tx.Rollback()
			return db.err
		}
	}

	return tx.Commit()
}

func (db DBWriter) writeLabel(tx *sql.Tx, l model.Label) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO labels (label_id, name, contact_info, profile, data_quality, urls) VALUES ('%s', '%s', '%s', '%s', '%s', ARRAY[%s])",
		l.ID,
		cleanText(l.Name),
		cleanText(l.ContactInfo),
		cleanText(l.Profile),
		l.DataQuality,
		array(l.Urls))
}

func (db DBWriter) writeLabelLabel(tx *sql.Tx, labelId, parent string, ll model.LabelLabel) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('%s', '%s', '%s', '%s')",
		labelId,
		ll.ID,
		cleanText(ll.Name),
		parent)
}

func (db DBWriter) writeLabelLabels(tx *sql.Tx, labelId, parent string, lls []model.LabelLabel) {
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

func (db DBWriter) writeMaster(tx *sql.Tx, m model.Master) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO masters (master_id, main_release, genres, styles, year, title, data_quality) VALUES ('%s', '%s', ARRAY[%s], ARRAY[%s], '%s', '%s', '%s')",
		m.ID,
		m.MainRelease,
		array(m.Genres),
		array(m.Styles),
		m.Year,
		cleanText(m.Title),
		m.DataQuality)
}

func (db DBWriter) writeRelease(tx *sql.Tx, r model.Release) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO releases (release_id, status, title, genres, styles, country, released, notes, data_quality, master_id, main_release) VALUES ('%s', '%s', '%s', ARRAY[%s], ARRAY[%s], '%s', '%s', '%s', '%s', '%s', '%s')",
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
		r.MainRelease)
}

func (db DBWriter) writeCompany(tx *sql.Tx, releaseId string, c model.Company) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO release_companies (release_id, release_company_id, name, category, entity_type, entity_type_name, resource_url) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s')",
		releaseId,
		c.ID,
		cleanText(c.Name),
		cleanText(c.Category),
		cleanText(c.EntityType),
		cleanText(c.EntityTypeName),
		cleanText(c.ResourceURL))
}

func (db DBWriter) writeCompanies(tx *sql.Tx, releaseId string, cs []model.Company) {
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

func (db DBWriter) writeReleaseArtist(tx *sql.Tx, masterId, releaseId, extra string, ra model.ReleaseArtist) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
		masterId,
		releaseId,
		ra.ID,
		cleanText(ra.Name),
		cleanText(extra),
		cleanText(ra.Join),
		cleanText(ra.Anv),
		cleanText(ra.Role),
		cleanText(ra.Tracks))
}

func (db DBWriter) writeReleaseArtists(tx *sql.Tx, masterId, releaseId, extra string, ras []model.ReleaseArtist) {
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

func (db DBWriter) writeFormat(tx *sql.Tx, releaseId string, f model.Format) {
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

func (db DBWriter) writeFormats(tx *sql.Tx, releaseId string, fs []model.Format) {
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

func (db DBWriter) writeTrack(tx *sql.Tx, releaseId string, t model.Track) {
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

func (db DBWriter) writeTrackList(tx *sql.Tx, releaseId string, tl []model.Track) {
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

func (db DBWriter) writeIdentifier(tx *sql.Tx, releaseId string, i model.Identifier) {
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

func (db DBWriter) writeIdentifiers(tx *sql.Tx, releaseId string, is []model.Identifier) {
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

func (db DBWriter) writeReleaseLabel(tx *sql.Tx, releaseId string, rl model.ReleaseLabel) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO release_labels (release_id, release_label_id, name, category) VALUES ('%s', '%s', '%s', '%s')",
		releaseId,
		rl.ID,
		cleanText(rl.Name),
		cleanText(rl.Category))
}

func (db DBWriter) writeReleaseLabels(tx *sql.Tx, releaseId string, rls []model.ReleaseLabel) {
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

func (db DBWriter) writeAlias(tx *sql.Tx, artistId string, a model.Alias) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO artist_aliases (artist_id, alias_id, name) VALUES ('%s', '%s', '%s')",
		artistId,
		a.ID,
		cleanText(a.Name))
}

func (db DBWriter) writeAliases(tx *sql.Tx, artistId string, as []model.Alias) {
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

func (db DBWriter) writeImage(tx *sql.Tx, artistId, labelId, masterId, releaseId string, img model.Image) {
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
			img.URI,
			img.URI150)
	}
}

func (db DBWriter) writeImages(tx *sql.Tx, artistId, labelId, masterId, releaseId string, imgs []model.Image) {
	if db.err == nil && !db.o.ExcludeImages {
		for _, img := range imgs {
			db.writeImage(tx, artistId, labelId, masterId, releaseId, img)
			if db.err != nil {
				return
			}
		}
	}
}

func (db DBWriter) writeVideo(tx *sql.Tx, masterId, releaseId string, v model.Video) {
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

func (db DBWriter) writeVideos(tx *sql.Tx, masterId, releaseId string, vs []model.Video) {
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

func (db DBWriter) writeArtist(tx *sql.Tx, a model.Artist) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO artists (artist_id, name, real_name, profile, data_quality, name_variations, urls) VALUES ('%s', '%s', '%s', '%s', '%s', ARRAY[%s], ARRAY[%s])",
		a.ID,
		cleanText(a.Name),
		cleanText(a.RealName),
		cleanText(a.Profile),
		a.DataQuality,
		array(a.NameVariations),
		array(a.Urls))
}

func (db DBWriter) writeArtistMember(tx *sql.Tx, artistId string, m model.Member) {
	if db.err != nil {
		return
	}

	db.writeTransaction(
		tx,
		"INSERT INTO artist_members (artist_id, member_id, name) VALUES ('%s', '%s', '%s')",
		artistId,
		m.ID,
		cleanText(m.Name))
}

func (db DBWriter) writeArtistMembers(tx *sql.Tx, artistId string, ms []model.Member) {
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

func (db DBWriter) writeTransaction(tx *sql.Tx, query string, values ...interface{}) {
	if db.err != nil {
		return
	}

	_, db.err = tx.Exec(fmt.Sprintf(query, values...))
}
