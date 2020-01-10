// Copyright (c) 2020 Lukas Aron. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package write

import (
	"strings"
	"testing"
)

func TestSQLWriter_Options(t *testing.T) {
	s := NewSQLWriter(nil, nil)
	opt := s.Options()

	if opt.ExcludeImages {
		t.Error("exclude images should be false as a default")
	}
}

func TestSQLWriter_WriteArtist(t *testing.T) {
	b := &strings.Builder{}
	s := NewSQLWriter(b, nil)

	err := s.WriteArtist(artists[0])
	if err != nil {
		t.Error(err)
	}

	got := b.String()
	if expectedArtist != got {
		t.Error("sql output differs from what it's expected")
	}
}

func TestSQLWriter_WriteArtists(t *testing.T) {
	b := &strings.Builder{}
	s := NewSQLWriter(b, nil)

	err := s.WriteArtists(artists)
	if err != nil {
		t.Error(err)
	}

	got := b.String()
	if expectedArtist != got {
		t.Error("sql output differs from what it's expected")
	}
}

func TestSQLWriter_WriteLabel(t *testing.T) {
	b := &strings.Builder{}
	s := NewSQLWriter(b, nil)

	err := s.WriteLabel(labels[0])
	if err != nil {
		t.Error(err)
	}

	got := b.String()
	if expectedLabel != got {
		t.Error("sql output differs from what it's expected")
	}
}

func TestSQLWriter_WriteLabels(t *testing.T) {
	b := &strings.Builder{}
	s := NewSQLWriter(b, nil)

	err := s.WriteLabels(labels)
	if err != nil {
		t.Error(err)
	}

	got := b.String()
	if expectedLabel != got {
		t.Error("sql output differs from what it's expected")
	}
}

func TestSQLWriter_WriteMaster(t *testing.T) {
	b := &strings.Builder{}
	s := NewSQLWriter(b, nil)

	err := s.WriteMaster(masters[0])
	if err != nil {
		t.Error(err)
	}

	got := b.String()
	if expectedMaster != got {
		t.Error("sql output differs from what it's expected")
	}
}

func TestSQLWriter_WriteMasters(t *testing.T) {
	b := &strings.Builder{}
	s := NewSQLWriter(b, nil)

	err := s.WriteMasters(masters)
	if err != nil {
		t.Error(err)
	}

	got := b.String()
	if expectedMaster != got {
		t.Error("sql output differs from what it's expected")
	}
}

func TestSQLWriter_WriteRelease(t *testing.T) {
	b := &strings.Builder{}
	s := NewSQLWriter(b, nil)

	err := s.WriteRelease(releases[0])
	if err != nil {
		t.Error(err)
	}

	got := b.String()
	if expectedRelease != got {
		t.Error("sql output differs from what it's expected")
	}
}

func TestSQLWriter_WriteReleases(t *testing.T) {
	b := &strings.Builder{}
	s := NewSQLWriter(b, nil)

	err := s.WriteReleases(releases)
	if err != nil {
		t.Error(err)
	}

	got := b.String()
	if expectedRelease != got {
		t.Error("sql output differs from what it's expected")
	}
}

// ------------------------------------------------------- DATA -------------------------------------------------------

var expectedArtist = `BEGIN;
INSERT INTO artists (artist_id, name, real_name, profile, data_quality, name_variations, urls) VALUES ('2', 'Mr. James Barth & A.D.', 'Cari Lekebusch & Alexi Delano', '', 'Correct', ARRAY['Mr Barth & A.D.','MR JAMES BARTH & A. D.','Mr. Barth & A.D.','Mr. James Barth & A. D.'], ARRAY['']);
INSERT INTO artist_aliases (artist_id, alias_id, name) VALUES ('2', '2470', 'Puente Latino');
INSERT INTO artist_aliases (artist_id, alias_id, name) VALUES ('2', '19536', 'Yakari & Delano');
INSERT INTO artist_aliases (artist_id, alias_id, name) VALUES ('2', '103709', 'Crushed Insect & The Sick Puppy');
INSERT INTO artist_aliases (artist_id, alias_id, name) VALUES ('2', '384581', 'ADCL');
INSERT INTO artist_aliases (artist_id, alias_id, name) VALUES ('2', '1779857', 'Alexi Delano & Cari Lekebusch');
INSERT INTO artist_members (artist_id, member_id, name) VALUES ('2', '26', 'Alexi Delano');
INSERT INTO artist_members (artist_id, member_id, name) VALUES ('2', '27', 'Cari Lekebusch');
COMMIT;
`

var expectedLabel = `BEGIN;
INSERT INTO labels (label_id, name, contact_info, profile, data_quality, urls) VALUES ('1', 'Planet E', 'Planet E Communications', '[a=Carl Craig]''s classic techno label founded in 1991.', 'Correct', ARRAY['http://planet-e.net','http://planetecommunications.bandcamp.com','http://www.facebook.com/planetedetroit','http://www.flickr.com/photos/planetedetroit','http://plus.google.com/100841702106447505236','http://www.instagram.com/carlcraignet','http://myspace.com/planetecom','http://myspace.com/planetedetroit','http://soundcloud.com/planetedetroit','http://twitter.com/planetedetroit','http://vimeo.com/user1265384','http://en.wikipedia.org/wiki/Planet_E_Communications','http://www.youtube.com/user/planetedetroit']);
INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('1', '86537', 'Antidote (4)', 'false');
INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('1', '41841', 'Community Projects', 'false');
INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('1', '153760', 'Guilty Pleasures', 'false');
INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('1', '31405', 'I Ner Zon Sounds', 'false');
INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('1', '277579', 'Planet E Communications', 'false');
INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('1', '294738', 'Planet E Communications, Inc.', 'false');
INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('1', '1560615', 'Planet E Productions', 'false');
INSERT INTO label_labels (label_id, sub_label_id, name, parent) VALUES ('1', '488315', 'TWPENTY', 'false');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '1', '', '', '24', '132', 'primary', '', '');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '1', '', '', '126', '587', 'secondary', '', '');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '1', '', '', '196', '600', 'secondary', '', '');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '1', '', '', '121', '275', 'secondary', '', '');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '1', '', '', '720', '382', 'secondary', '', '');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '1', '', '', '398', '500', 'secondary', '', '');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '1', '', '', '189', '600', 'secondary', '', '');
COMMIT;
`
var expectedMaster = `BEGIN;
INSERT INTO masters (master_id, main_release, genres, styles, year, title, data_quality) VALUES ('18512', '33699', ARRAY['Electronic'], ARRAY['Tribal','Techno'], '2002', 'Psyche EP', 'Correct');
INSERT INTO release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ('18512', '', '212070', 'Samuel L Session', 'false', '', '', '', '');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '', '18512', '', '150', '150', 'primary', '', '');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '', '18512', '', '592', '600', 'secondary', '', '');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '', '18512', '', '592', '600', 'secondary', '', '');
INSERT INTO videos (master_id, release_id, duration, embed, src, title, description) VALUES ('18512', '', '118', 'true', 'https://www.youtube.com/watch?v=QYf4j0Pd2FU', 'Samuel L. Session - Arrival', 'Samuel L. Session - Arrival');
INSERT INTO videos (master_id, release_id, duration, embed, src, title, description) VALUES ('18512', '', '376', 'true', 'https://www.youtube.com/watch?v=c_AfLqTdncI', 'Samuel L. Session - Psyche Part 1', 'Samuel L. Session - Psyche Part 1');
INSERT INTO videos (master_id, release_id, duration, embed, src, title, description) VALUES ('18512', '', '419', 'true', 'https://www.youtube.com/watch?v=0nxvR8Zl9wY', 'Samuel L. Session - Psyche Part 2', 'Samuel L. Session - Psyche Part 2');
COMMIT;
`
var expectedRelease = `BEGIN;
INSERT INTO releases (release_id, status, title, genres, styles, country, released, notes, data_quality, master_id, main_release) VALUES ('2', 'Accepted', 'Knockin'' Boots Vol 2 Of 2', ARRAY['Electronic'], ARRAY['Broken Beat','Techno','Tech House'], 'Sweden', '1998-06-00', 'All joints recorded in NYC (Dec.97).', 'Correct', '713738', 'true');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '', '', '2', '394', '400', 'primary', '', '');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '', '', '2', '600', '600', 'secondary', '', '');
INSERT INTO images (artist_id, label_id, master_id, release_id, height, width, type, uri, uri_150) VALUES ('', '', '', '2', '600', '600', 'secondary', '', '');
INSERT INTO release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ('', '2', '2', 'Mr. James Barth & A.D.', 'false', '', '', '', '');
INSERT INTO release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ('', '2', '26', 'Alexi Delano', 'true', '', '', 'Producer, Recorded By', '');
INSERT INTO release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ('', '2', '27', 'Cari Lekebusch', 'true', '', '', 'Producer, Recorded By', '');
INSERT INTO release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ('', '2', '26', 'Alexi Delano', 'true', '', 'A. Delano', 'Written-By', '');
INSERT INTO release_artists (master_id, release_id, release_artist_id, name, extra, joiner, anv, role, tracks) VALUES ('', '2', '27', 'Cari Lekebusch', 'true', '', 'C. Lekebusch', 'Written-By', '');
INSERT INTO release_formats (release_id, name, quantity, text, descriptions) VALUES ('2', 'Vinyl', '1', '', ARRAY['12"','33 ⅓ RPM']);
INSERT INTO release_tracks (release_id, position, title, duration) VALUES ('2', 'A1', 'A Sea Apart', '5:08');
INSERT INTO release_tracks (release_id, position, title, duration) VALUES ('2', 'A2', 'Dutchmaster', '4:21');
INSERT INTO release_tracks (release_id, position, title, duration) VALUES ('2', 'B1', 'Inner City Lullaby', '4:22');
INSERT INTO release_tracks (release_id, position, title, duration) VALUES ('2', 'B2', 'Yeah Kid!', '4:46');
INSERT INTO release_identifiers (release_id, description, type, value) VALUES ('2', 'Side A Runout Etching', 'Matrix / Runout', 'MPO SK026-A -J.T.S.-');
INSERT INTO release_identifiers (release_id, description, type, value) VALUES ('2', 'Side B Runout Etching', 'Matrix / Runout', 'MPO SK026-B -J.T.S.-');
INSERT INTO release_labels (release_id, release_label_id, name, category) VALUES ('2', '5', 'Svek', 'SK 026');
INSERT INTO release_labels (release_id, release_label_id, name, category) VALUES ('2', '5', 'Svek', 'SK026');
INSERT INTO release_companies (release_id, release_company_id, name, category, entity_type, entity_type_name, resource_url) VALUES ('2', '266169', 'JTS Studios', '', '29', 'Mastered At', 'https://api.discogs.com/labels/266169');
INSERT INTO release_companies (release_id, release_company_id, name, category, entity_type, entity_type_name, resource_url) VALUES ('2', '56025', 'MPO', '', '17', 'Pressed By', 'https://api.discogs.com/labels/56025');
INSERT INTO videos (master_id, release_id, duration, embed, src, title, description) VALUES ('', '2', '310', 'true', 'https://www.youtube.com/watch?v=MIgQNVhYILA', 'Mr. James Barth & A.D. - A Sea Apart', 'Mr. James Barth & A.D. - A Sea Apart');
INSERT INTO videos (master_id, release_id, duration, embed, src, title, description) VALUES ('', '2', '265', 'true', 'https://www.youtube.com/watch?v=LgLchSRehhc', 'Mr. James Barth & A.D. - Dutchmaster', 'Mr. James Barth & A.D. - Dutchmaster');
INSERT INTO videos (master_id, release_id, duration, embed, src, title, description) VALUES ('', '2', '260', 'true', 'https://www.youtube.com/watch?v=iaqHaULlqqg', 'Mr. James Barth & A.D. - Inner City Lullaby', 'Mr. James Barth & A.D. - Inner City Lullaby');
INSERT INTO videos (master_id, release_id, duration, embed, src, title, description) VALUES ('', '2', '290', 'true', 'https://www.youtube.com/watch?v=x_Os7b-iWKs', 'Mr. James Barth & A.D. - Yeah Kid!', 'Mr. James Barth & A.D. - Yeah Kid!');
COMMIT;
`
