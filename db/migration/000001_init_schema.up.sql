CREATE TABLE users (
  id   BIGSERIAL PRIMARY KEY,
  name      text      NOT NULL,
  email     text   NOT NULL,
  phone     text  NOT NULL
);