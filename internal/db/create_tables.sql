DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS State_to_auth;

CREATE TABLE IF NOT EXISTS Users (
  id                TEXT NOT NULL,
  username          VARCHAR(128) NOT NULL,
  pass              VARCHAR(255) NOT NULL,
  user_state        VARCHAR(128) NOT NULL,
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

