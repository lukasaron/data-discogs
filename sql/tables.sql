
CREATE TABLE artists (
    artist_id VARCHAR(10),
    name VARCHAR(1024),
    real_name VARCHAR(1024),
    profile TEXT,
    data_quality VARCHAR(20),
    name_variations TEXT,
    urls TEXT
);

CREATE INDEX artists_artist_id ON artists(artist_id);
CREATE INDEX artists_name ON artists(name);
CREATE INDEX artists_real_name ON artists(real_name);
CREATE INDEX artists_data_quality ON artists(data_quality);

CREATE TABLE artist_aliases (
    artist_id VARCHAR(10),
    alias_id VARCHAR(10),
    name VARCHAR(1024)
);

CREATE INDEX artist_aliases_artist_id ON artist_aliases(artist_id);
CREATE INDEX artist_aliases_alias_id ON artist_aliases(alias_id);

CREATE TABLE images (
  artist_id VARCHAR(10),
  label_id VARCHAR(10),
  master_id VARCHAR(10),
  release_id VARCHAR(10),
  height VARCHAR(10),
  width VARCHAR(10),
  type VARCHAR(10),
  uri VARCHAR(1024),
  uri_150 VARCHAR(1024)
);

CREATE INDEX images_artist_id ON images(artist_id);
CREATE INDEX images_label_id ON images(label_id);
CREATE INDEX images_master_id ON images(master_id);
CREATE INDEX images_release_id ON images(release_id);

CREATE TABLE labels (
    label_id VARCHAR(10),
    name VARCHAR(1024),
    contact_info TEXT,
    profile TEXT,
    data_quality VARCHAR(20),
    urls TEXT
);

CREATE INDEX labels_label_id ON labels(label_id);
CREATE INDEX labels_name ON labels(name);
CREATE INDEX labels_data_quality ON labels(data_quality);

CREATE TABLE label_labels (
    label_id VARCHAR(10),
    sub_label_id VARCHAR(10),
    name VARCHAR(1024),
    parent VARCHAR(5)
);

CREATE INDEX label_labels_label_id ON label_labels(label_id);
CREATE INDEX label_labels_sub_label_id ON label_labels(sub_label_id);
CREATE INDEX label_labels_name ON label_labels(name);

CREATE TABLE masters (
    master_id VARCHAR(10),
    main_release VARCHAR(10),
    genres TEXT,
    styles TEXT,
    year VARCHAR(4),
    title VARCHAR(1024),
    data_quality VARCHAR(20)
);

CREATE INDEX masters_master_id ON masters(master_id);
CREATE INDEX masters_data_quality ON masters(data_quality);

CREATE TABLE videos (
    master_id VARCHAR(10),
    release_id VARCHAR(10),
    duration VARCHAR(10),
    embed VARCHAR(5),
    src VARCHAR(1024),
    title VARCHAR(1024),
    description TEXT
);

CREATE INDEX videos_master_id ON videos(master_id);
CREATE INDEX videos_release_id ON videos(release_id);
CREATE INDEX videos_title ON videos(title);

CREATE TABLE releases (
    release_id VARCHAR(10),
    status VARCHAR(20),
    title VARCHAR(100),
    genres TEXT,
    styles TEXT,
    country VARCHAR(50),
    released VARCHAR(10),
    notes TEXT,
    data_quality VARCHAR(20),
    master_id VARCHAR(10),
    main_release VARCHAR(5)
);

CREATE INDEX releases_release_id ON releases(release_id);
CREATE INDEX releases_status ON releases(status);
CREATE INDEX releases_title ON releases(title);
CREATE INDEX releases_country ON releases(country);
CREATE INDEX releases_released ON releases(released);
CREATE INDEX releases_master_id ON releases(master_id);

CREATE TABLE release_artists (
    master_id VARCHAR(10),
    release_id VARCHAR(10),
    release_artist_id VARCHAR(10),
    name VARCHAR(1024),
    extra VARCHAR(5),
    joiner TEXT,
    anv TEXT,
    role TEXT,
    tracks TEXT
);

CREATE INDEX release_artists_master_id ON release_artists(master_id);
CREATE INDEX release_artists_release_id ON release_artists(release_id);
CREATE INDEX release_artists_name ON release_artists(name);

CREATE TABLE release_labels (
    release_id VARCHAR(10),
    release_label_id VARCHAR(10),
    name VARCHAR(1024),
    category VARCHAR(100)
);

CREATE INDEX release_labels_release_id ON release_labels(release_id);
CREATE INDEX release_labels_release_label_id ON release_labels(release_label_id);
CREATE INDEX release_labels_name ON release_labels(name);
CREATE INDEX release_labels_category ON release_labels(category);

CREATE TABLE release_identifiers (
    release_id VARCHAR(10),
    description TEXT,
    type TEXT,
    value TEXT
);

CREATE INDEX release_identifiers_release_id ON release_identifiers(release_id);

CREATE TABLE release_formats (
    release_id VARCHAR(10),
    name VARCHAR(1024),
    quantity VARCHAR(10),
    text TEXT,
    description TEXT
);

CREATE INDEX release_formats_release_id ON release_formats(release_id);
CREATE INDEX release_formats_name ON release_formats(name);

CREATE TABLE release_companies (
    release_id VARCHAR(10),
    release_company_id VARCHAR(10),
    name VARCHAR(1024),
    category VARCHAR(100),
    entity_name VARCHAR(1024),
    entity_type_name VARCHAR(1024),
    resource_url VARCHAR(1024)
);

CREATE INDEX release_companies_release_id ON release_companies(release_id);
CREATE INDEX release_companies_release_company_id ON release_companies(release_company_id);
CREATE INDEX release_companies_name ON release_companies(name);
CREATE INDEX release_companies_category ON release_companies(category);

CREATE TABLE release_tracks (
    release_id VARCHAR(10),
    position VARCHAR(5),
    title VARCHAR(100),
    duration VARCHAR(10)
);

CREATE INDEX release_tracks_release_id ON release_tracks(release_id);