-- +migrate Up
CREATE TABLE domains (
  name text PRIMARY KEY,
  ip text NOT NULL
);
-- +migrate Down
DROP TABLE domains;