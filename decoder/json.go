package decoder

import (
	"encoding/json"
)

func (d *Decoder) JsonArtists(limit int) (int, []byte, error) {
	if d.Error == nil {
		n, a, err := d.Artists(limit)
		if err != nil {
			return n, nil, err
		}

		b, err := json.Marshal(a)
		if err != nil {
			return 0, b, err
		}

		return n, b, nil
	}

	return 0, nil, d.Error
}

func (d *Decoder) JsonLabels(limit int) (int, []byte, error) {
	if d.Error == nil {
		n, l, err := d.Labels(limit)
		if err != nil {
			return n, nil, err
		}

		b, err := json.Marshal(l)
		if err != nil {
			return 0, b, err
		}

		return n, b, nil
	}

	return 0, nil, d.Error
}

func (d *Decoder) JsonMasters(limit int) (int, []byte, error) {
	if d.Error == nil {
		n, m, err := d.Masters(limit)
		if err != nil {
			return n, nil, err
		}

		b, err := json.Marshal(m)
		if err != nil {
			return 0, nil, err
		}

		return n, b, nil
	}

	return 0, nil, d.Error
}

func (d *Decoder) JsonReleases(limit int) (int, []byte, error) {
	if d.Error == nil {
		n, r, err := d.Releases(limit)
		if err != nil {
			return n, nil, err
		}

		b, err := json.Marshal(r)
		if err != nil {
			return 0, nil, err
		}

		return n, b, nil
	}

	return 0, nil, d.Error
}
