DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS State_to_auth;
DROP TABLE IF EXISTS Rooms;
DROP TABLE IF EXISTS Votes;
DROP TABLE IF EXISTS Song_Suggestion;

CREATE TABLE IF NOT EXISTS Users (
  id                UUID UNIQUE NOT NULL,
  username          VARCHAR(128) UNIQUE NOT NULL,
  pass              VARCHAR(255) NOT NULL,
  user_state        UUID UNIQUE NOT NULL,
  finished_creating INTEGER NOT NULL,
  PRIMARY KEY (`id`)
  FOREIGN KEY (user_state) REFERENCES State_to_Code(user_state)
);

CREATE TABLE IF NOT EXISTS State_to_auth (
    user_state        UUID NOT NULL,
    auth_token        VARCHAR(128) NOT NULL,
    refresh_token     VARCHAR(128) NOT NULL,
    expiration_time   VARCHAR(128) NOT NULL,
    PRIMARY KEY (`user_state`)
);

CREATE TABLE IF NOT EXISTS Rooms (
  room_id UUID NOT NULL,
  room_name VARCHAR(128) NOT NULL,
  owner_id UUID NOT NULL,

  PRIMARY KEY (`room_id`)
  FOREIGN KEY (owner_id) REFERENCES Users(id)
);

CREATE TABLE IF NOT EXISTS Song_Suggestion (
  suggestion_id UUID,
  song_name VARCHAR(128) NOT NULL,
  song_id VARCHAR(128) NOT NULL,
  room_id UUID NOT NULL,

  PRIMARY KEY (`suggestion_id`)
  FOREIGN KEY (room_id) REFERENCES Rooms(room_id)
);

CREATE TABLE IF NOT EXISTS Votes (
  vote_id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  vote INTEGER NOT NULL,
  vote_timestamp VARCHAR(128) NOT NULL,
  suggestion_id UUID NOT NULL,

  FOREIGN KEY (suggestion_id) REFERENCES Song_Suggestion(suggestion_id)
);

INSERT INTO Users VALUES(
  "123456789",
  "test",
  "$2a$04$TOJGmLmeq8/y9cyV5XHtnOEn307hqVx8xyNUXpBC3lCo0sZcePMWK",
  "user-state",
  1
);