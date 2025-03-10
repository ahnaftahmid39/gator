-- +goose Up
create table users (
  id UUID,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  name TEXT NOT NULL UNIQUE
);

-- +goose Down
drop table users;
