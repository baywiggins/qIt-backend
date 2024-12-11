DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS State_to_auth;

CREATE TABLE IF NOT EXISTS Users (
  id                TEXT UNIQUE NOT NULL,
  username          VARCHAR(128) UNIQUE NOT NULL,
  pass              VARCHAR(255) NOT NULL,
  user_state        VARCHAR(128) UNIQUE NOT NULL,
  finished_creating INTEGER NOT NULL,
  PRIMARY KEY (`id`)
  FOREIGN KEY (user_state) REFERENCES State_to_Code(user_state)
);

CREATE TABLE IF NOT EXISTS State_to_auth (
    user_state       VARCHAR(128) NOT NULL,
    auth_token        VARCHAR(128) NOT NULL,
    refresh_token     VARCHAR(128) NOT NULL,
    expiration_time   TIME NOT NULL,
    PRIMARY KEY (`user_state`)
);

INSERT INTO Users VALUES(
  "123456789",
  "test",
  "$2a$04$TOJGmLmeq8/y9cyV5XHtnOEn307hqVx8xyNUXpBC3lCo0sZcePMWK",
  "user-state",
  1
)