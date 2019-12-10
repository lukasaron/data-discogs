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
		SslMode  string `default:"disable" env:"DB_SSL_MODE"`
	}
	File struct {
		Name string `env:"FILE_NAME"`
		Type string `env:"FILE_TYPE"`
	}
	Block struct {
		Size  int `default:"10000" env:"BLOCK_SIZE"`
		Skip  int `default:"0" env:"BLOCK_SKIP"`
		Limit int `default:"2147483647" env:"BLOCK_LIMIT"`
	}
}

func Start() (err error) {
	err = configor.Load(&Config)
	if err != nil {
		return err
	}

	ft := getDecoderFileType(Config.File.Type)
	if ft == decoder.Unknown {
		return wrongTypeSpecified
	}

	d := decoder.NewDecoder(Config.File.Name)
	defer d.Close()

	pg := writer.NewPostgres(
		Config.DB.Host,
		Config.DB.Name,
		Config.DB.User,
		Config.DB.Password,
		Config.DB.SslMode,
		Config.DB.Port)

	defer pg.Close()

	return decodeData(d, pg, ft)
}

func getDecoderFileType(fileType string) (ft decoder.FileType) {
	switch fileType {
	case "artists":
		ft = decoder.Artists
	case "labels":
		ft = decoder.Labels
	case "masters":
		ft = decoder.Masters
	case "releases":
		ft = decoder.Releases
	default:
		ft = decoder.Unknown
	}

	return ft
}

func decodeData(d decoder.Decoder, w writer.Writer, ft decoder.FileType) error {
	fn, err := getDecodeFunction(ft)
	if err != nil {
		return err
	}

	blockCount := 1
	for ; blockCount <= Config.Block.Limit; blockCount++ {
		err = fn(d, w, blockCount > Config.Block.Skip)
		if err != nil {
			_ = fmt.Errorf("Block %d failed\n", blockCount)
			return err
		}

		if blockCount > Config.Block.Skip {
			fmt.Printf("Block %d written\n", blockCount)
		} else {
			fmt.Printf("Block %d skipped\n", blockCount)
		}

	}

	return nil
}

func getDecodeFunction(ft decoder.FileType) (func(decoder.Decoder, writer.Writer, bool) error, error) {
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

func decodeArtists(d decoder.Decoder, w writer.Writer, write bool) error {
	num, a, err := d.Artists(Config.Block.Size)
	if err != nil || num == 0 {
		return err
	}

	if write {
		return w.WriteArtists(a)
	}

	return nil
}

func decodeLabels(d decoder.Decoder, w writer.Writer, write bool) error {
	num, l, err := d.Labels(Config.Block.Size)
	if err != nil || num == 0 {
		return err
	}

	if write {
		return w.WriteLabels(l)
	}

	return nil
}

func decodeMasters(d decoder.Decoder, w writer.Writer, write bool) error {
	num, m, err := d.Masters(Config.Block.Size)
	if err != nil || num == 0 {
		return err
	}

	if write {
		return w.WriteMasters(m)
	}

	return nil
}

func decodeReleases(d decoder.Decoder, w writer.Writer, write bool) error {
	num, r, err := d.Releases(Config.Block.Size)
	if err != nil || num == 0 {
		return err
	}

	if write {
		return w.WriteReleases(r)
	}

	return nil
}
