
CREATE TABLE artists (
    artist_id VARCHAR(10),
    name VARCHAR(1024),
    real_name VARCHAR(1024),
    profile TEXT,
    data_quality VARCHAR(20),
    name_variations VARCHAR(1024)[],
    urls VARCHAR(1024)[]
);

CREATE TABLE artist_aliases (
    artist_id VARCHAR(10),
    alias_id VARCHAR(10),
    name VARCHAR(1024)
);

CREATE TABLE artist_members (
    artist_id VARCHAR(10),
    member_id VARCHAR(10),
    name VARCHAR(1024)
);

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

CREATE TABLE labels (
    label_id VARCHAR(10),
    name VARCHAR(1024),
    contact_info TEXT,
    profile TEXT,
    data_quality VARCHAR(20),
    urls VARCHAR(1024)[]
);

CREATE TABLE label_labels (
    label_id VARCHAR(10),
    sub_label_id VARCHAR(10),
    name VARCHAR(1024),
    parent VARCHAR(5)
);

CREATE TABLE masters (
    master_id VARCHAR(10),
    main_release VARCHAR(10),
    genres VARCHAR(1024)[],
    styles VARCHAR(1024)[],
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
    release_id VARCHAR(10),
    status VARCHAR(20),
    title VARCHAR(100),
    genres VARCHAR(1024)[],
    styles VARCHAR(1024)[],
    country VARCHAR(50),
    released VARCHAR(10),
    notes TEXT,
    data_quality VARCHAR(20),
    master_id VARCHAR(10),
    main_release VARCHAR(10)
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
    descriptions TEXT[]
);

CREATE TABLE release_companies (
    release_id VARCHAR(10),
    release_company_id VARCHAR(10),
    name VARCHAR(1024),
    category VARCHAR(100),
    entity_type VARCHAR(1024),
    entity_type_name VARCHAR(1024),
    resource_url VARCHAR(1024)
);

CREATE TABLE release_tracks (
    release_id VARCHAR(10),
    position VARCHAR(10),
    title VARCHAR(100),
    duration VARCHAR(10)
);
