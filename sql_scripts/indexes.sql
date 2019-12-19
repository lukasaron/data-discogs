
CREATE INDEX artists_artist_id ON artists(artist_id);
CREATE INDEX artists_name ON artists(name);
CREATE INDEX artists_real_name ON artists(real_name);
CREATE INDEX artists_data_quality ON artists(data_quality);

CREATE INDEX artist_aliases_artist_id ON artist_aliases(artist_id);
CREATE INDEX artist_aliases_alias_id ON artist_aliases(alias_id);

CREATE INDEX artist_members_artist_id ON artist_members(artist_id);
CREATE INDEX artist_members_member_id ON artist_members(member_id);

CREATE INDEX images_artist_id ON images(artist_id);
CREATE INDEX images_label_id ON images(label_id);
CREATE INDEX images_master_id ON images(master_id);
CREATE INDEX images_release_id ON images(release_id);

CREATE INDEX labels_label_id ON labels(label_id);
CREATE INDEX labels_name ON labels(name);
CREATE INDEX labels_data_quality ON labels(data_quality);

CREATE INDEX label_labels_label_id ON label_labels(label_id);
CREATE INDEX label_labels_sub_label_id ON label_labels(sub_label_id);
CREATE INDEX label_labels_name ON label_labels(name);

CREATE INDEX masters_master_id ON masters(master_id);
CREATE INDEX masters_data_quality ON masters(data_quality);

CREATE INDEX videos_master_id ON videos(master_id);
CREATE INDEX videos_release_id ON videos(release_id);
CREATE INDEX videos_title ON videos(title);

CREATE INDEX releases_release_id ON releases(release_id);
CREATE INDEX releases_status ON releases(status);
CREATE INDEX releases_title ON releases(title);
CREATE INDEX releases_country ON releases(country);
CREATE INDEX releases_released ON releases(released);
CREATE INDEX releases_master_id ON releases(master_id);

CREATE INDEX release_artists_master_id ON release_artists(master_id);
CREATE INDEX release_artists_release_id ON release_artists(release_id);
CREATE INDEX release_artists_name ON release_artists(name);

CREATE INDEX release_labels_release_id ON release_labels(release_id);
CREATE INDEX release_labels_release_label_id ON release_labels(release_label_id);
CREATE INDEX release_labels_name ON release_labels(name);
CREATE INDEX release_labels_category ON release_labels(category);

CREATE INDEX release_identifiers_release_id ON release_identifiers(release_id);

CREATE INDEX release_formats_release_id ON release_formats(release_id);
CREATE INDEX release_formats_name ON release_formats(name);

CREATE INDEX release_companies_release_id ON release_companies(release_id);
CREATE INDEX release_companies_release_company_id ON release_companies(release_company_id);
CREATE INDEX release_companies_name ON release_companies(name);
CREATE INDEX release_companies_category ON release_companies(category);

CREATE INDEX release_tracks_release_id ON release_tracks(release_id);