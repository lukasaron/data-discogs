package cmd

import (
	"errors"
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
	FileName string `env:"FILE_NAME"`
	FileType string `env:"FILE_TYPE"`
	Num      int    `default:"10000" env:"NUM"`
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
	switch ft {
	case decoder.Artists:
		return decodeArtists(d, pg)
	case decoder.Labels:
		return decodeLabels(d, pg)
	case decoder.Masters:
		return decodeMasters(d, pg)
	case decoder.Releases:
		return decodeReleases(d, pg)
	default:
		return wrongTypeSpecified
	}

}

func decodeArtists(d *decoder.Decoder, pg *writer.Postgres) error {
	for {
		num, a, err := d.Artists(Config.Num)
		if err != nil || num == 0 {
			return err
		}

		err = pg.WriteArtists(a)
		if err != nil {
			return err
		}
	}
}

func decodeLabels(d *decoder.Decoder, pg *writer.Postgres) error {
	for {
		num, l, err := d.Labels(Config.Num)
		if err != nil || num == 0 {
			return err
		}

		err = pg.WriteLabels(l)
		if err != nil {
			return err
		}
	}
}

func decodeMasters(d *decoder.Decoder, pg *writer.Postgres) error {
	for {
		num, m, err := d.Masters(Config.Num)
		if err != nil || num == 0 {
			return err
		}

		err = pg.WriteMasters(m)
		if err != nil {
			return err
		}
	}
}

func decodeReleases(d *decoder.Decoder, pg *writer.Postgres) error {
	for {
		num, r, err := d.Releases(Config.Num)
		if err != nil || num == 0 {
			return err
		}

		err = pg.WriteReleases(r)
		if err != nil {
			return err
		}
	}
}
