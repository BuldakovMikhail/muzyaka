CREATE TYPE ALBUM_TYPE AS ENUM ('single', 'LP', 'EP');
CREATE TYPE ROLE_TYPE AS ENUM ('user', 'musician', 'admin');

CREATE TABLE IF NOT EXISTS musicians
(
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name        VARCHAR(254) NOT NULL,
    description TEXT,
    CHECK ( name <> '' )
);

CREATE TABLE IF NOT EXISTS albums
(
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    cover_file  BYTEA        NOT NULL,
    type        ALBUM_TYPE   NOT NULL,
    musician_id INT          NOT NULL
        REFERENCES musicians (id)
            ON DELETE CASCADE,
    CHECK ( name <> '' ),
    CHECK ( length(cover_file) > 0 )
);


CREATE TABLE IF NOT EXISTS genres
(
    id   INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    CHECK ( name <> '' )
);


CREATE TABLE IF NOT EXISTS tracks
(
    id       INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    source   VARCHAR(254) NOT NULL,
    name     VARCHAR(100) NOT NULL,
    genre    INT REFERENCES genres (id),
    album_id INT          NOT NULL
        REFERENCES albums (id)
            ON DELETE CASCADE,
    CHECK ( source <> '' ),
    CHECK ( name <> '' )
);

CREATE TABLE IF NOT EXISTS merch
(
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name        VARCHAR(254) NOT NULL,
    description TEXT,
    link        VARCHAR(254) NOT NULL,
    musician_id INT          NOT NULL
        REFERENCES musicians (id)
            ON DELETE CASCADE,
    CHECK ( name <> '' ),
    CHECK ( link <> '' ),
    UNIQUE (name, link, musician_id)
);


CREATE TABLE IF NOT EXISTS merch_photos
(
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    photo_file BYTEA NOT NULL,
    merch_id   INT   NOT NULL
        REFERENCES merch (id)
            ON DELETE CASCADE,
    CHECK ( length(photo_file) > 0 )
);


CREATE TABLE IF NOT EXISTS musicians_photos
(
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    photo_file  BYTEA NOT NULL,
    musician_id INT   NOT NULL
        REFERENCES musicians (id)
            ON DELETE CASCADE,
    CHECK ( length(photo_file) > 0 )
);

CREATE TABLE IF NOT EXISTS users
(
    id       INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name     VARCHAR(100) NOT NULL,
    email    VARCHAR(254) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    role     ROLE_TYPE DEFAULT 'user',
    CHECK ( name <> '' ),
    CHECK ( email <> '' ),
    CHECK ( password <> '' )
);

CREATE TABLE IF NOT EXISTS playlists
(
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name        VARCHAR(254) NOT NULL,
    cover_file  BYTEA        NOT NULL,
    description TEXT,
    user_id     INT          NOT NULL
        REFERENCES users (id)
            ON DELETE CASCADE,
    CHECK ( name <> '' ),
    CHECK ( length(cover_file) > 0 )
);

CREATE TABLE IF NOT EXISTS users_musicians
(
    user_id     INT NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    musician_id INT NOT NULL UNIQUE REFERENCES musicians (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, musician_id)
);

CREATE TABLE IF NOT EXISTS outbox
(
    id       SERIAL PRIMARY KEY,
    event_id TEXT                  NOT NULL,
    track_id INTEGER               NOT NULL,
    source   VARCHAR(254),
    name     VARCHAR(100),
    genre    INT,
    type     VARCHAR(100)          NOT NULL,
    sent     BOOLEAN DEFAULT FALSE NOT NULL
);


-- ----------------- LINKS ----------------------------------

CREATE TABLE IF NOT EXISTS track_playlist
(
    track_id    INT NOT NULL REFERENCES tracks (id) ON DELETE CASCADE,
    playlist_id INT NOT NULL REFERENCES playlists (id) ON DELETE CASCADE,
    PRIMARY KEY (track_id, playlist_id)
);

CREATE TABLE IF NOT EXISTS user_track
(
    track_id INT NOT NULL REFERENCES tracks (id) ON DELETE CASCADE,
    user_id  INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (track_id, user_id)
);
