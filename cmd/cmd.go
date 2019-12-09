package cmd

import (
	"errors"
	"fmt"
	"github.com/Twyer/discogs/decoder"
	"github.com/Twyer/discogs/writer"
	"github.com/jinzhu/configor"
)

var wrongTypeSpecified = errors.New("wrong file type specified")

var Config struct {
	DB struct {
		Host     string `default:"localhost" env:"DB_HOST"`
		Name     string `default:"discogs" env:"DB_NAME"`
		User     string `default:"user" env:"DB_USERNAME"`
		Password string `default:"password" env:"DB_PASSWORD"`
		Port     int    `default:"5432" env:"DB_PORT"`
	}
	FileName   string `env:"FILE_NAME"`
	FileType   string `env:"FILE_TYPE"`
	BlockSize  int    `default:"10000" env:"BLOCK_SIZE"`
	DropBlocks int    `default:"0" env:"DROP_BLOCKS"`
}

func Start() (err error) {
	err = configor.Load(&Config)
	if err != nil {
		return err
	}

	ft := getDecoderFileType(Config.FileType)
	if ft == decoder.Unknown {
		return wrongTypeSpecified
	}

	d := decoder.NewDecoder(Config.FileName)
	if d.Error != nil {
		return err
	}
	defer d.Close()

	pg := writer.NewPostgres(
		Config.DB.Host,
		Config.DB.Name,
		Config.DB.User,
		Config.DB.Password,
		"disable",
		Config.DB.Port)
	if pg.Error != nil {
		return pg.Error
	}
	defer pg.Close()

	return decodeData(d, pg, ft)
}

func getDecoderFileType(fileType string) decoder.FileType {
	switch fileType {
	case "artists":
		return decoder.Artists
	case "labels":
		return decoder.Labels
	case "masters":
		return decoder.Masters
	case "releases":
		return decoder.Releases
	default:
		return decoder.Unknown
	}
}

func decodeData(d *decoder.Decoder, pg *writer.Postgres, ft decoder.FileType) error {
	fn, err := getDecodeFunction(ft)
	if err != nil {
		return err
	}

	blockCount := 1
	for ; ; blockCount++ {
		err = fn(d, pg, blockCount >= Config.DropBlocks)
		if err != nil {
			fmt.Errorf("Block %d failed\n", blockCount)
			return err
		}

		fmt.Printf("Block %d written\n", blockCount)
	}
}

func getDecodeFunction(ft decoder.FileType) (func(*decoder.Decoder, *writer.Postgres, bool) error, error) {
	switch ft {
	case decoder.Artists:
		return decodeArtists, nil
	case decoder.Labels:
		return decodeLabels, nil
	case decoder.Masters:
		return decodeMasters, nil
	case decoder.Releases:
		return decodeReleases, nil
	default:
		return nil, wrongTypeSpecified
	}
}

func decodeArtists(d *decoder.Decoder, pg *writer.Postgres, write bool) error {
	num, a, err := d.Artists(Config.BlockSize)
	if err != nil || num == 0 {
		return err
	}

	if write {
		return pg.WriteArtists(a)
	}

	return nil
}

func decodeLabels(d *decoder.Decoder, pg *writer.Postgres, write bool) error {
	num, l, err := d.Labels(Config.BlockSize)
	if err != nil || num == 0 {
		return err
	}

	if write {
		return pg.WriteLabels(l)
	}

	return nil
}

func decodeMasters(d *decoder.Decoder, pg *writer.Postgres, write bool) error {
	num, m, err := d.Masters(Config.BlockSize)
	if err != nil || num == 0 {
		return err
	}

	if write {
		return pg.WriteMasters(m)
	}

	return nil
}

func decodeReleases(d *decoder.Decoder, pg *writer.Postgres, write bool) error {
	num, r, err := d.Releases(Config.BlockSize)
	if err != nil || num == 0 {
		return err
	}

	if write {
		return pg.WriteReleases(r)
	}

	return nil
}
