CREATE TABLE watson_results (
    id bigserial primary key,
    created_at TIMESTAMPTZ,
    contents text,
    user_id integer REFERENCES users(id) NOT NULL
);

CREATE TABLE keywords (
  id bigserial primary key,
  contents text,
  sentiment float8 NOT NULL,
  relevance float8 NOT NULL,
  sadness float8 NOT NULL,
  joy float8 NOT NULL,
  fear float8 NOT NULL,
  disgust float8 NOT NULL,
  anger float8 NOT NULL,
  watson_id integer  REFERENCES watson_results(id) NOT NULL
);
