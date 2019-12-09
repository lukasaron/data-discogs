
CREATE TABLE artists (
    id SERIAL PRIMARY KEY,
    artist_id VARCHAR(10),
    name VARCHAR(1024),
    real_name VARCHAR(1024),
    profile TEXT,
    data_quality VARCHAR(20),
    name_variations TEXT,
    urls TEXT
);

CREATE TABLE artist_aliases (
    artist_id VARCHAR(10),
    alias_id VARCHAR(10),
    name VARCHAR(1024)
);

CREATE TABLE images (
  id SERIAL PRIMARY KEY,
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

CREATE TABLE labels (
    id SERIAL PRIMARY KEY,
    label_id VARCHAR(10),
    name VARCHAR(1024),
    contact_info TEXT,
    profile TEXT,
    data_quality VARCHAR(20),
    urls TEXT
);

CREATE TABLE label_labels (
    label_id VARCHAR(10),
    sub_label_id VARCHAR(10),
    name VARCHAR(1024),
    parent VARCHAR(5)
);

CREATE TABLE masters (
    id SERIAL PRIMARY KEY,
    master_id VARCHAR(10),
    main_release VARCHAR(10),
    genres TEXT,
    styles TEXT,
    year VARCHAR(4),
    title VARCHAR(1024),
    data_quality VARCHAR(20)
);

CREATE TABLE videos (
    master_id VARCHAR(10),
    release_id VARCHAR(10),
    duration VARCHAR(10),
    embed VARCHAR(5),
    src VARCHAR(1024),
    title VARCHAR(1024),
    description TEXT
);

CREATE TABLE releases (
    id SERIAL PRIMARY KEY,
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

CREATE TABLE release_labels (
    release_id VARCHAR(10),
    release_label_id VARCHAR(10),
    name VARCHAR(1024),
    category VARCHAR(100)
);

CREATE TABLE release_identifiers (
    release_id VARCHAR(10),
    description TEXT,
    type TEXT,
    value TEXT
);

CREATE TABLE release_formats (
    release_id VARCHAR(10),
    name VARCHAR(1024),
    quantity VARCHAR(10),
    text TEXT,
    description TEXT
);

CREATE TABLE release_companies (
    release_id VARCHAR(10),
    release_company_id VARCHAR(10),
    name VARCHAR(1024),
    category VARCHAR(100),
    entity_name VARCHAR(1024),
    entity_type_name VARCHAR(1024),
    resource_url VARCHAR(1024)
);


CREATE TABLE release_tracks (
    release_id VARCHAR(10),
    position VARCHAR(5),
    title VARCHAR(100),
    duration VARCHAR(10)
);