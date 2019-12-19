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

func (pg PostgresWriter) Reset() error {
	pg.err = nil
	return nil
}

func (pg PostgresWriter) Close() error {
	return pg.db.Close()
}

func (pg PostgresWriter) WriteArtist(artist model.Artist) error {
	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	pg.writeArtist(tx, artist)
	pg.writeAliases(tx, artist.Id, artist.Aliases)
	pg.writeImages(tx, artist.Id, "", "", "", artist.Images)

	if pg.err != nil {
		_ = tx.Rollback()
		return pg.err
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteArtists(artists []model.Artist) error {
	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, a := range artists {
		pg.writeArtist(tx, a)
		pg.writeAliases(tx, a.Id, a.Aliases)
		pg.writeImages(tx, a.Id, "", "", "", a.Images)

		if pg.err != nil {
			_ = tx.Rollback()
			return pg.err
		}
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteLabel(label model.Label) error {
	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	pg.writeLabel(tx, label)
	pg.writeLabelLabels(tx, label.Id, "false", label.SubLabels)
	pg.writeImages(tx, "", label.Id, "", "", label.Images)

	if label.ParentLabel != nil {
		pg.writeLabelLabel(tx, label.Id, "true", *label.ParentLabel)
	}

	if pg.err != nil {
		_ = tx.Rollback()
		return pg.err
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteLabels(labels []model.Label) error {
	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, l := range labels {
		pg.writeLabel(tx, l)
		pg.writeLabelLabels(tx, l.Id, "false", l.SubLabels)
		pg.writeImages(tx, "", l.Id, "", "", l.Images)

		if l.ParentLabel != nil {
			pg.writeLabelLabel(tx, l.Id, "true", *l.ParentLabel)
		}

		if pg.err != nil {
			_ = tx.Rollback()
			return pg.err
		}
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteMaster(master model.Master) error {
	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	pg.writeMaster(tx, master)
	pg.writeImages(tx, "", "", master.Id, "", master.Images)
	pg.writeReleaseArtists(tx, master.Id, "", "false", master.Artists)
	pg.writeVideos(tx, master.Id, "", master.Videos)

	if pg.err != nil {
		_ = tx.Rollback()
		return pg.err
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteMasters(masters []model.Master) error {
	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, m := range masters {
		pg.writeMaster(tx, m)
		pg.writeImages(tx, "", "", m.Id, "", m.Images)
		pg.writeReleaseArtists(tx, m.Id, "", "false", m.Artists)
		pg.writeVideos(tx, m.Id, "", m.Videos)

		if pg.err != nil {
			_ = tx.Rollback()
			return pg.err
		}
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteRelease(release model.Release) error {
	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	pg.writeRelease(tx, release)
	pg.writeImages(tx, "", "", "", release.Id, release.Images)
	pg.writeReleaseArtists(tx, "", release.Id, "false", release.Artists)
	pg.writeReleaseArtists(tx, "", release.Id, "true", release.ExtraArtists)
	pg.writeFormats(tx, release.Id, release.Formats)
	pg.writeTrackList(tx, release.Id, release.TrackList)
	pg.writeIdentifiers(tx, release.Id, release.Identifiers)
	pg.writeVideos(tx, "", release.Id, release.Videos)
	pg.writeReleaseLabels(tx, release.Id, release.Labels)
	pg.writeCompanies(tx, release.Id, release.Companies)
	if pg.err != nil {
		_ = tx.Rollback()
		return pg.err
	}

	return tx.Commit()
}

func (pg PostgresWriter) WriteReleases(releases []model.Release) error {
	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, r := range releases {
		pg.writeRelease(tx, r)
		pg.writeImages(tx, "", "", "", r.Id, r.Images)
		pg.writeReleaseArtists(tx, "", r.Id, "false", r.Artists)
		pg.writeReleaseArtists(tx, "", r.Id, "true", r.ExtraArtists)
		pg.writeFormats(tx, r.Id, r.Formats)
		pg.writeTrackList(tx, r.Id, r.TrackList)
		pg.writeIdentifiers(tx, r.Id, r.Identifiers)
		pg.writeVideos(tx, "", r.Id, r.Videos)
		pg.writeReleaseLabels(tx, r.Id, r.Labels)
		pg.writeCompanies(tx, r.Id, r.Companies)

		if pg.err != nil {
			_ = tx.Rollback()
			return pg.err
		}
	}

	return tx.Commit()
}

func (pg PostgresWriter) writeTransaction(tx *sql.Tx, query string, values ...interface{}) {
	if pg.err != nil {
		return
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		pg.err = err
		return
	}

	_, pg.err = stmt.Exec(values...)
}

func (pg PostgresWriter) writeLabel(tx *sql.Tx, l model.Label) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
		tx,
		"INSERT INTO public.labels (label_id, name, contact_info, profile, data_quality, urls) VALUES ($1, $2, $3, $4, $5, $6)",
		l.Id,
		l.Name,
		l.ContactInfo,
		l.Profile,
		l.DataQuality,
		pq.Array(l.Urls))
}

func (pg PostgresWriter) writeLabelLabel(tx *sql.Tx, labelId, parent string, ll model.LabelLabel) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
		tx,
		"INSERT INTO public.label_labels (label_id, sub_label_id, name, parent) VALUES ($1, $2, $3, $4)",
		labelId,
		ll.Id,
		ll.Name,
		parent)
}

func (pg PostgresWriter) writeLabelLabels(tx *sql.Tx, labelId, parent string, lls []model.LabelLabel) {
	if pg.err != nil {
		return
	}

	for _, ll := range lls {
		pg.writeLabelLabel(tx, labelId, parent, ll)
		if pg.err != nil {
			return
		}
	}
}

func (pg PostgresWriter) writeMaster(tx *sql.Tx, m model.Master) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
		tx,
		"INSERT INTO public.masters (master_id, main_release, genres, styles, year, title, data_quality) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		m.Id,
		m.MainRelease,
		pq.Array(m.Genres),
		pq.Array(m.Styles),
		m.Year,
		m.Title,
		m.DataQuality)
}

func (pg PostgresWriter) writeRelease(tx *sql.Tx, r model.Release) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
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
}

func (pg PostgresWriter) writeCompany(tx *sql.Tx, releaseId string, c model.Company) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
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

func (pg PostgresWriter) writeCompanies(tx *sql.Tx, releaseId string, cs []model.Company) {
	if pg.err != nil {
		return
	}

	for _, c := range cs {
		pg.writeCompany(tx, releaseId, c)
		if pg.err != nil {
			return
		}
	}
}

func (pg PostgresWriter) writeReleaseArtist(tx *sql.Tx, masterId, releaseId, extra string, ra model.ReleaseArtist) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
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

func (pg PostgresWriter) writeReleaseArtists(tx *sql.Tx, masterId, releaseId, extra string, ras []model.ReleaseArtist) {
	if pg.err != nil {
		return
	}

	for _, ra := range ras {
		pg.writeReleaseArtist(tx, masterId, releaseId, extra, ra)
		if pg.err != nil {
			return
		}
	}
}

func (pg PostgresWriter) writeFormat(tx *sql.Tx, releaseId string, f model.Format) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
		tx,
		"INSERT INTO public.release_formats (release_id, name, quantity, text, descriptions) VALUES ($1, $2, $3, $4, $5)",
		releaseId,
		f.Name,
		f.Quantity,
		f.Text,
		pq.Array(f.Descriptions))
}

func (pg PostgresWriter) writeFormats(tx *sql.Tx, releaseId string, fs []model.Format) {
	if pg.err != nil {
		return
	}

	for _, f := range fs {
		pg.writeFormat(tx, releaseId, f)
		if pg.err != nil {
			return
		}
	}
}

func (pg PostgresWriter) writeTrack(tx *sql.Tx, releaseId string, t model.Track) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
		tx,
		"INSERT INTO public.release_tracks (release_id, position, title, duration) VALUES ($1, $2, $3, $4)",
		releaseId,
		t.Position,
		t.Title,
		t.Duration)
}

func (pg PostgresWriter) writeTrackList(tx *sql.Tx, releaseId string, tl []model.Track) {
	if pg.err != nil {
		return
	}

	for _, t := range tl {
		pg.writeTrack(tx, releaseId, t)
		if pg.err != nil {
			return
		}
	}
}

func (pg PostgresWriter) writeIdentifier(tx *sql.Tx, releaseId string, i model.Identifier) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
		tx,
		"INSERT INTO public.release_identifiers (release_id, description, type, value) VALUES ($1, $2, $3, $4)",
		releaseId,
		i.Description,
		i.Type,
		i.Value)
}

func (pg PostgresWriter) writeIdentifiers(tx *sql.Tx, releaseId string, is []model.Identifier) {
	if pg.err != nil {
		return
	}

	for _, i := range is {
		pg.writeIdentifier(tx, releaseId, i)
		if pg.err != nil {
			return
		}
	}
}

func (pg PostgresWriter) writeReleaseLabel(tx *sql.Tx, releaseId string, rl model.ReleaseLabel) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
		tx,
		"INSERT INTO public.release_labels (release_id, release_label_id, name, category) VALUES ($1, $2, $3, $4)",
		releaseId,
		rl.Id,
		rl.Name,
		rl.Category)
}

func (pg PostgresWriter) writeReleaseLabels(tx *sql.Tx, releaseId string, rls []model.ReleaseLabel) {
	if pg.err != nil {
		return
	}

	for _, rl := range rls {
		pg.writeReleaseLabel(tx, releaseId, rl)
		if pg.err != nil {
			return
		}
	}
}

func (pg PostgresWriter) writeAlias(tx *sql.Tx, artistId string, a model.Alias) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
		tx,
		"INSERT INTO public.artist_aliases (artist_id, alias_id, name) VALUES ($1, $2, $3)",
		artistId,
		a.Id,
		a.Name)
}

func (pg PostgresWriter) writeAliases(tx *sql.Tx, artistId string, as []model.Alias) {
	if pg.err != nil {
		return
	}

	for _, a := range as {
		pg.writeAlias(tx, artistId, a)
		if pg.err != nil {
			return
		}
	}
}

func (pg PostgresWriter) writeImage(tx *sql.Tx, artistId, labelId, masterId, releaseId string, img model.Image) {
	if pg.err == nil && !pg.options.ExcludeImages {
		pg.writeTransaction(
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
}

func (pg PostgresWriter) writeImages(tx *sql.Tx, artistId, labelId, masterId, releaseId string, imgs []model.Image) {
	if pg.err == nil && !pg.options.ExcludeImages {
		for _, img := range imgs {
			pg.writeImage(tx, artistId, labelId, masterId, releaseId, img)
			if pg.err != nil {
				return
			}
		}
	}
}

func (pg PostgresWriter) writeVideo(tx *sql.Tx, masterId, releaseId string, v model.Video) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
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

func (pg PostgresWriter) writeVideos(tx *sql.Tx, masterId, releaseId string, vs []model.Video) {
	if pg.err != nil {
		return
	}

	for _, v := range vs {
		pg.writeVideo(tx, masterId, releaseId, v)
		if pg.err != nil {
			return
		}
	}
}

func (pg PostgresWriter) writeArtist(tx *sql.Tx, a model.Artist) {
	if pg.err != nil {
		return
	}

	pg.writeTransaction(
		tx,
		"INSERT INTO public.artists (artist_id, name, real_name, profile, data_quality, name_variations, urls) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		a.Id,
		a.Name,
		a.RealName,
		a.Profile,
		a.DataQuality,
		pq.Array(a.NameVariations),
		pq.Array(a.Urls))
}
