package writer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Twyer/discogs/decoder"
	"github.com/Twyer/discogs/model"
	_ "github.com/lib/pq"
)

type Postgres struct {
	d     *decoder.Decoder
	db    *sql.DB
	Error error
}

func NewPostgres(host, dbName, user, password, sslMode string, port int) *Postgres {
	connStr := fmt.Sprintf("host='%s' dbname='%s' user='%s' password='%s' port='%d' sslmode=%s",
		host,
		dbName,
		user,
		password,
		port,
		sslMode)

	pg := Postgres{}
	pg.db, pg.Error = sql.Open("postgres", connStr)
	return &pg
}

func (pg *Postgres) Ping() error {
	return pg.db.Ping()
}

func (pg *Postgres) Close() error {
	return pg.db.Close()
}

func (pg *Postgres) writeObject(obj interface{}, query string) error {
	tx, err := pg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	b, err := json.Marshal(obj)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = stmt.Exec(b)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg *Postgres) WriteArtist(artist model.Artist) error {
	return pg.writeObject(artist, "INSERT INTO public.artists (artist) VALUES ($1)")
}

func (pg *Postgres) WriteLabel(label model.Label) error {
	return pg.writeObject(label, "INSERT INTO public.labels (label) VALUES ($1)")
}

func (pg *Postgres) WriteMaster(master model.Master) error {
	return pg.writeObject(master, "INSERT INTO public.masters (master) VALUES ($1)")
}

func (pg *Postgres) WriteRelease(release model.Release) error {
	return pg.writeObject(release, "INSERT INTO public.releases (release) VALUES ($1)")
}
