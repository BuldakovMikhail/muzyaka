CREATE TYPE ALBUM_TYPE AS ENUM ('single', 'LP', 'EP');
CREATE TYPE ROLE_TYPE AS ENUM ('user', 'musician', 'admin');

CREATE TABLE IF NOT EXISTS albums (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    cover_file BYTEA,
    type ALBUM_TYPE NOT NULL
);

CREATE TABLE IF NOT EXISTS genres(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(254) NOT NULL UNIQUE
);


CREATE TABLE IF NOT EXISTS tracks(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    source VARCHAR(254) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    genre INT REFERENCES genres(id)
);


CREATE TABLE IF NOT EXISTS merch(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(254) NOT NULL,
    description TEXT,
    link VARCHAR(254) NOT NULL
);

CREATE TABLE IF NOT EXISTS merch_photos(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    photo_file BYTEA,
    merch_id INT NOT NULL REFERENCES merch(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS musicians(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(254) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS musicians_photos(
                                               photo_src VARCHAR(254) NOT NULL,
                                               musician_id INT NOT NULL REFERENCES musicians(id) ON DELETE CASCADE,
                                               PRIMARY KEY (photo_src, musician_id)
);

CREATE TABLE IF NOT EXISTS playlists(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(254) NOT NULL,
    cover VARCHAR(254) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS users(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    role ROLE_TYPE DEFAULT 'user'
);

-- ----------------- LINKS ----------------------------------

CREATE TABLE IF NOT EXISTS albums_tracks(
    album_id INT NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    track_id INT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    PRIMARY KEY (album_id, track_id)
);

CREATE TABLE IF NOT EXISTS musicians_albums(
    album_id INT NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    musician_id INT NOT NULL REFERENCES musicians(id) ON DELETE CASCADE,
    PRIMARY KEY (album_id, musician_id)
);

CREATE TABLE IF NOT EXISTS merch_musicians(
                                              merch_id INT NOT NULL REFERENCES merch(id) ON DELETE CASCADE,
                                              musician_id INT NOT NULL REFERENCES musicians(id) ON DELETE CASCADE,
                                              PRIMARY KEY (merch_id, musician_id)
);

CREATE TABLE IF NOT EXISTS playlists_tracks(
                                              track_id INT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
                                              playlist_id INT NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
                                              PRIMARY KEY (track_id, playlist_id)
);

CREATE TABLE IF NOT EXISTS users_playlists(
                                               user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                               playlist_id INT NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
                                               is_favourite BOOLEAN NOT NULL DEFAULT FALSE,
                                               is_public BOOLEAN NOT NULL DEFAULT FALSE,
                                               PRIMARY KEY (user_id, playlist_id)
);

