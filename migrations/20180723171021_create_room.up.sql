CREATE TABLE rooms (
  id bigserial primary key,
  name varchar(64) NOT NULL UNIQUE
);
