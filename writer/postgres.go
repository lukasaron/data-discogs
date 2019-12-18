package writer

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Twyer/discogs-parser/model"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresWriter struct {
	options Options
	db      *sql.DB
	err     error
}

func NewPostgresWriter(host string, port int, dbName, user, password, sslMode string, options ...Options) Writer {
	pg := PostgresWriter{}

	connStr := fmt.Sprintf("host='%s' dbname='%s' user='%s' password='%s' port='%d' sslmode=%s",
		host,
		dbName,
		user,
		password,
		port,
		sslMode)

	pg.db, pg.err = sql.Open("postgres", connStr)

	// add options when available (only the first one)
	if options != nil && len(options) > 0 {
		pg.options = options[0]
	}

	return pg
}

func (pg PostgresWriter) Close() error {
	if pg.err != nil {
		return pg.err
	}

	return pg.db.Close()
}

func (pg PostgresWriter) WriteArtist(a model.Artist) error {
	if pg.err != nil {
		return pg.err
	}

	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	err = pg.writeArtist(tx, a)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteArtists(as []model.Artist) error {
	if pg.err != nil {
		return pg.err
	}

	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, a := range as {
		err = pg.writeArtist(tx, a)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteLabel(l model.Label) error {
	if pg.err != nil {
		return pg.err
	}

	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	err = pg.writeLabel(tx, l)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteLabels(ls []model.Label) error {
	if pg.err != nil {
		return pg.err
	}

	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, l := range ls {
		err = pg.writeLabel(tx, l)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteMaster(m model.Master) error {
	if pg.err != nil {
		return pg.err
	}

	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	err = pg.writeMaster(tx, m)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteMasters(ms []model.Master) error {
	if pg.err != nil {
		return pg.err
	}

	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, m := range ms {
		err = pg.writeMaster(tx, m)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteRelease(r model.Release) error {
	if pg.err != nil {
		return pg.err
	}

	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	err = pg.writeRelease(tx, r)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteReleases(rs []model.Release) error {
	if pg.err != nil {
		return pg.err
	}

	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, r := range rs {
		err = pg.writeRelease(tx, r)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (pg PostgresWriter) writeTransaction(tx *sql.Tx, query string, values ...interface{}) error {

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(values...)
	return err
}

func (pg PostgresWriter) writeLabel(tx *sql.Tx, l model.Label) (err error) {
	err = pg.writeTransaction(
		tx,
		"INSERT INTO public.labels (label_id, name, contact_info, profile, data_quality, urls) VALUES ($1, $2, $3, $4, $5, $6)",
		l.Id,
		l.Name,
		l.ContactInfo,
		l.Profile,
		l.DataQuality,
		pq.Array(l.Urls))

	if err != nil {
		return err
	}

	if l.ParentLabel != nil {
		err = pg.writeLabelLabel(tx, l.Id, "true", *l.ParentLabel)
		if err != nil {
			return err
		}
	}

	err = pg.writeLabelLabels(tx, l.Id, "false", l.SubLabels)
	if err != nil {
		return err
	}

	return pg.writeImages(tx, "", l.Id, "", "", l.Images)
}

func (pg PostgresWriter) writeLabelLabel(tx *sql.Tx, labelId, parent string, ll model.LabelLabel) error {
	return pg.writeTransaction(
		tx,
		"INSERT INTO public.label_labels (label_id, sub_label_id, name, parent) VALUES ($1, $2, $3, $4)",
		labelId,
		ll.Id,
		ll.Name,
		parent)
}

func (pg PostgresWriter) writeLabelLabels(tx *sql.Tx, labelId, parent string, lls []model.LabelLabel) (err error) {
	for _, ll := range lls {
		err = pg.writeLabelLabel(tx, labelId, parent, ll)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg PostgresWriter) writeMaster(tx *sql.Tx, m model.Master) (err error) {
	err = pg.writeTransaction(
		tx,
		"INSERT INTO public.masters (master_id, main_release, genres, styles, year, title, data_quality) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		m.Id,
		m.MainRelease,
		pq.Array(m.Genres),
		pq.Array(m.Styles),
		m.Year,
		m.Title,
		m.DataQuality)

	if err != nil {
		return err
	}

	err = pg.writeImages(tx, "", "", m.Id, "", m.Images)
	if err != nil {
		return err
	}

	err = pg.writeReleaseArtists(tx, m.Id, "", "false", m.Artists)
	if err != nil {
		return err
	}

	return pg.writeVideos(tx, m.Id, "", m.Videos)
}

func (pg PostgresWriter) writeRelease(tx *sql.Tx, r model.Release) (err error) {
	err = pg.writeTransaction(
		tx,
		"INSERT INTO public.releases (release_id, status, title, genres, styles, country, released, notes, data_quality, master_id, main_release) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		r.Id,
		r.Status,
		r.Title,
		pq.Array(r.Genres),
		pq.Array(r.Styles),
		r.Country,
		r.Released,
		r.Notes,
		r.DataQuality,
		r.MasterId,
		r.MainRelease)

	if err != nil {
		return err
	}

	err = pg.writeImages(tx, "", "", "", r.Id, r.Images)
	if err != nil {
		return err
	}

	err = pg.writeReleaseArtists(tx, "", r.Id, "false", r.Artists)
	if err != nil {
		return err
	}

	err = pg.writeReleaseArtists(tx, "", r.Id, "true", r.ExtraArtists)
	if err != nil {
		return err
	}

	err = pg.writeFormats(tx, r.Id, r.Formats)
	if err != nil {
		return err
	}

	err = pg.writeTrackList(tx, r.Id, r.TrackList)
	if err != nil {
		return err
	}

	err = pg.writeIdentifiers(tx, r.Id, r.Identifiers)
	if err != nil {
		return err
	}

	err = pg.writeVideos(tx, "", r.Id, r.Videos)
	if err != nil {
		return err
	}

	err = pg.writeReleaseLabels(tx, r.Id, r.Labels)
	if err != nil {
		return err
	}

	return pg.writeCompanies(tx, r.Id, r.Companies)
}

func (pg PostgresWriter) writeCompany(tx *sql.Tx, releaseId string, c model.Company) error {
	return pg.writeTransaction(
		tx,
		"INSERT INTO public.release_companies (release_id, release_company_id, name, category, entity_type, entity_type_name, resource_url) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		releaseId,
		c.Id,
		c.Name,
		c.Category,
		c.EntityType,
		c.EntityTypeName,
		c.ResourceUrl)
}

func (pg PostgresWriter) writeCompanies(tx *sql.Tx, releaseId string, cs []model.Company) (err error) {
	for _, c := range cs {
		err = pg.writeCompany(tx, releaseId, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg PostgresWriter) writeReleaseArtist(tx *sql.Tx, masterId, releaseId, extra string, ra model.ReleaseArtist) error {

	return pg.writeTransaction(
		tx,
		"INSERT INTO public.release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		masterId,
		releaseId,
		ra.Id,
		ra.Name,
		extra,
		ra.Join,
		ra.Anv,
		ra.Role,
		ra.Tracks)
}

func (pg PostgresWriter) writeReleaseArtists(tx *sql.Tx, masterId, releaseId, extra string, ras []model.ReleaseArtist) (err error) {
	for _, ra := range ras {
		err = pg.writeReleaseArtist(tx, masterId, releaseId, extra, ra)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg PostgresWriter) writeFormat(tx *sql.Tx, releaseId string, f model.Format) error {
	return pg.writeTransaction(
		tx,
		"INSERT INTO public.release_formats (release_id, name, quantity, text, descriptions) VALUES ($1, $2, $3, $4, $5)",
		releaseId,
		f.Name,
		f.Quantity,
		f.Text,
		pq.Array(f.Descriptions))
}

func (pg PostgresWriter) writeFormats(tx *sql.Tx, releaseId string, fs []model.Format) (err error) {
	for _, f := range fs {
		err = pg.writeFormat(tx, releaseId, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg PostgresWriter) writeTrack(tx *sql.Tx, releaseId string, t model.Track) error {
	return pg.writeTransaction(
		tx,
		"INSERT INTO public.release_tracks (release_id, position, title, duration) VALUES ($1, $2, $3, $4)",
		releaseId,
		t.Position,
		t.Title,
		t.Duration)
}

func (pg PostgresWriter) writeTrackList(tx *sql.Tx, releaseId string, tl []model.Track) (err error) {
	for _, t := range tl {
		err = pg.writeTrack(tx, releaseId, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg PostgresWriter) writeIdentifier(tx *sql.Tx, releaseId string, i model.Identifier) error {
	return pg.writeTransaction(
		tx,
		"INSERT INTO public.release_identifiers (release_id, description, type, value) VALUES ($1, $2, $3, $4)",
		releaseId,
		i.Description,
		i.Type,
		i.Value)
}

func (pg PostgresWriter) writeIdentifiers(tx *sql.Tx, releaseId string, is []model.Identifier) (err error) {
	for _, i := range is {
		err = pg.writeIdentifier(tx, releaseId, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg PostgresWriter) writeReleaseLabel(tx *sql.Tx, releaseId string, rl model.ReleaseLabel) error {
	return pg.writeTransaction(
		tx,
		"INSERT INTO public.release_labels (release_id, release_label_id, name, category) VALUES ($1, $2, $3, $4)",
		releaseId,
		rl.Id,
		rl.Name,
		rl.Category)
}

func (pg PostgresWriter) writeReleaseLabels(tx *sql.Tx, releaseId string, rls []model.ReleaseLabel) (err error) {
	for _, rl := range rls {
		err = pg.writeReleaseLabel(tx, releaseId, rl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg PostgresWriter) writeAlias(tx *sql.Tx, artistId string, a model.Alias) error {
	return pg.writeTransaction(
		tx,
		"INSERT INTO public.artist_aliases (artist_id, alias_id, name) VALUES ($1, $2, $3)",
		artistId,
		a.Id,
		a.Name)
}

func (pg PostgresWriter) writeAliases(tx *sql.Tx, artistId string, as []model.Alias) (err error) {
	for _, a := range as {
		err = pg.writeAlias(tx, artistId, a)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg PostgresWriter) writeImage(tx *sql.Tx, artistId, labelId, masterId, releaseId string, img model.Image) error {
	if !pg.options.ExcludeImages {
		return pg.writeTransaction(
			tx,
			"INSERT INTO public.images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
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

	return nil
}

func (pg PostgresWriter) writeImages(tx *sql.Tx, artistId, labelId, masterId, releaseId string, imgs []model.Image) (err error) {
	if !pg.options.ExcludeImages {
		for _, img := range imgs {
			err = pg.writeImage(tx, artistId, labelId, masterId, releaseId, img)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pg PostgresWriter) writeVideo(tx *sql.Tx, masterId, releaseId string, v model.Video) (err error) {
	return pg.writeTransaction(
		tx,
		"INSERT INTO public.videos (master_id, release_id, duration, embed, src, title, description) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		masterId,
		releaseId,
		v.Duration,
		v.Embed,
		v.Src,
		v.Title,
		v.Description)
}

func (pg PostgresWriter) writeVideos(tx *sql.Tx, masterId, releaseId string, vs []model.Video) (err error) {
	for _, v := range vs {
		err = pg.writeVideo(tx, masterId, releaseId, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg PostgresWriter) writeArtist(tx *sql.Tx, a model.Artist) (err error) {
	err = pg.writeTransaction(
		tx,
		"INSERT INTO public.artists (artist_id, name, real_name, profile, data_quality, name_variations, urls) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		a.Id,
		a.Name,
		a.RealName,
		a.Profile,
		a.DataQuality,
		pq.Array(a.NameVariations),
		pq.Array(a.Urls))

	if err != nil {
		return err
	}

	err = pg.writeAliases(tx, a.Id, a.Aliases)
	if err != nil {
		return err
	}

	return pg.writeImages(tx, a.Id, "", "", "", a.Images)
}
