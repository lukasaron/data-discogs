package decoder

import (
	"errors"
	"github.com/lukasaron/discogs-parser/model"
	"github.com/lukasaron/discogs-parser/writer"
	"github.com/sirupsen/logrus"
	"io"
)

type Decoder interface {
	Close() error
	Artists(limit int) (int, []model.Artist, error)
	Labels(limit int) (int, []model.Label, error)
	Masters(limit int) (int, []model.Master, error)
	Releases(limit int) (int, []model.Release, error)
}

type Options struct {
	QualityLevel QualityLevel
}

var wrongTypeSpecified = errors.New("wrong file type specified")

func DecodeData(decoder Decoder, fileType FileType, writer writer.Writer, blockSize, blockLimit, blockSkip int) error {
	if blockLimit < 0 {
		blockLimit = int(^uint(0) >> 1)
	}

	fn, err := getDecodeFunction(fileType)
	if err != nil {
		return err
	}

	blockCount := 1
	for ; blockCount <= blockLimit; blockCount++ {
		// call appropriate decoder function
		num, err := fn(decoder, writer, blockSize, blockCount > blockSkip)
		if err != nil && err != io.EOF {
			logrus.Errorf("Block %d failed [%d]\n", blockCount, num)
			return err
		}

		if num == 0 && err == io.EOF {
			break
		}

		if blockCount > blockSkip {
			logrus.Infof("Block %d written [%d]\n", blockCount, num)
		} else {
			logrus.Infof("Block %d skipped [%d]\n", blockCount, num)
		}
	}

	return nil
}

func getDecodeFunction(ft FileType) (func(Decoder, writer.Writer, int, bool) (int, error), error) {
	switch ft {
	case Artists:
		return decodeArtists, nil
	case Labels:
		return decodeLabels, nil
	case Masters:
		return decodeMasters, nil
	case Releases:
		return decodeReleases, nil
	default:
		return nil, wrongTypeSpecified
	}
}

func decodeArtists(d Decoder, w writer.Writer, blockSize int, write bool) (int, error) {
	num, a, err := d.Artists(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteArtists(a)
	}

	return num, err
}

func decodeLabels(d Decoder, w writer.Writer, blockSize int, write bool) (int, error) {
	num, l, err := d.Labels(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteLabels(l)
	}

	return num, err
}

func decodeMasters(d Decoder, w writer.Writer, blockSize int, write bool) (int, error) {
	num, m, err := d.Masters(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteMasters(m)
	}

	return num, err
}

func decodeReleases(d Decoder, w writer.Writer, blockSize int, write bool) (int, error) {
	num, r, err := d.Releases(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteReleases(r)
	}

	return num, err
}
