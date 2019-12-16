package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Twyer/discogs/decoder"
	"github.com/Twyer/discogs/writer"
	"github.com/jinzhu/configor"
	"io"
	"regexp"
)

var (
	ftRegex = regexp.MustCompile(".*(artists|labels|masters|releases).*\\.xml")

	wrongTypeSpecified = errors.New("wrong file type specified")

	Config struct {
		DB struct {
			Host     string `default:"localhost" env:"DB_HOST"`
			Name     string `default:"discogs" env:"DB_NAME"`
			User     string `default:"user" env:"DB_USERNAME"`
			Password string `default:"password" env:"DB_PASSWORD"`
			Port     int    `default:"5432" env:"DB_PORT"`
			SslMode  string `default:"disable" env:"DB_SSL_MODE"`
		}
		File struct {
			Name string `default:"" env:"FILE_NAME"`
		}
		Block struct {
			Size  int `default:"1000" env:"BLOCK_SIZE"`
			Skip  int `default:"0" env:"BLOCK_SKIP"`
			Limit int `default:"2147483647" env:"BLOCK_LIMIT"`
		}
		Filter struct {
			Quality string `default:"All" env:"FILTER_QUALITY"`
		}
		Writer struct {
			Type   string `default:"json" env:"WRITER_TYPE"`
			Output string `default:"" env:"WRITER_OUTPUT"`
		}
	}
)

func Run() (err error) {
	flag.StringVar(&Config.File.Name, "filename", "", "input file")
	flag.StringVar(&Config.Filter.Quality, "quality", "All", "quality filter")
	flag.StringVar(&Config.Writer.Type, "writer-type", "json", "writer type")
	flag.StringVar(&Config.Writer.Output, "output", "", "writer output")
	flag.IntVar(&Config.Block.Size, "block-size", 1000, "block size")
	flag.IntVar(&Config.Block.Skip, "block-skip", 0, "block skip")
	flag.IntVar(&Config.Block.Limit, "block-limit", 2147483647, "block limit")
	flag.Parse()

	err = configor.Load(&Config)
	if err != nil {
		return err
	}

	ft := getFileTypeFromFileName(Config.File.Name)
	if ft == decoder.Unknown {
		return wrongTypeSpecified
	}

	d := decoder.NewDecoder(
		Config.File.Name,
		decoder.Options{
			QualityLevel: decoder.StrToQualityLevel(Config.Filter.Quality),
		},
	)
	defer d.Close()

	w := getWriter()
	defer w.Close()

	return decodeData(d, w, ft)
}

func getWriter() (w writer.Writer) {
	wt := writer.StrToWriterType(Config.Writer.Type)
	switch wt {
	case writer.PostgresType:
		w = writer.NewPostgres(
			Config.DB.Host,
			Config.DB.Port,
			Config.DB.Name,
			Config.DB.User,
			Config.DB.Password,
			Config.DB.SslMode,
			writer.Options{ExcludeImages: true},
		)
	case writer.JsonType:
		w = writer.NewJson(
			Config.Writer.Output,
			writer.Options{ExcludeImages: true},
		)
	}

	return w
}

func getFileTypeFromFileName(fileName string) decoder.FileType {
	ftStr := ""
	ftSubMatches := ftRegex.FindStringSubmatch(fileName)
	if len(ftSubMatches) > 1 {
		ftStr = ftSubMatches[1]
	}

	return getDecoderFileType(ftStr)
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
		num, err := fn(d, w, blockCount > Config.Block.Skip)
		if err != nil && err != io.EOF {
			_ = fmt.Errorf("Block %d failed [%d]\n", blockCount, num)
			return err
		}

		if num == 0 && err == io.EOF {
			break
		}

		if blockCount > Config.Block.Skip {
			fmt.Printf("Block %d written [%d]\n", blockCount, num)
		} else {
			fmt.Printf("Block %d skipped [%d]\n", blockCount, num)
		}
	}

	return nil
}

func getDecodeFunction(ft decoder.FileType) (func(decoder.Decoder, writer.Writer, bool) (int, error), error) {
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

func decodeArtists(d decoder.Decoder, w writer.Writer, write bool) (int, error) {
	num, a, err := d.Artists(Config.Block.Size)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteArtists(a)
	}

	return num, err
}

func decodeLabels(d decoder.Decoder, w writer.Writer, write bool) (int, error) {
	num, l, err := d.Labels(Config.Block.Size)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteLabels(l)
	}

	return num, err
}

func decodeMasters(d decoder.Decoder, w writer.Writer, write bool) (int, error) {
	num, m, err := d.Masters(Config.Block.Size)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteMasters(m)
	}

	return num, err
}

func decodeReleases(d decoder.Decoder, w writer.Writer, write bool) (int, error) {
	num, r, err := d.Releases(Config.Block.Size)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteReleases(r)
	}

	return num, err
}
