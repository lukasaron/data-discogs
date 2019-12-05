CREATE TABLE artists (
    id SERIAL PRIMARY KEY,
    artist JSON NOT NULL
);

CREATE TABLE masters (
    id SERIAL PRIMARY KEY,
    master JSON NOT NULL
);

CREATE TABLE labels (
    id SERIAL PRIMARY KEY,
    label JSON NOT NULL
);

CREATE TABLE releases (
    id SERIAL PRIMARY KEY,
    release JSON NOT NULL
);
