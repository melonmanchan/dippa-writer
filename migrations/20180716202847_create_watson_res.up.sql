CREATE TABLE watson_results (
    id bigserial primary key,
    created_at TIMESTAMPTZ,
    contents text,
    user_id integer REFERENCES users(id) NOT NULL
);

CREATE TABLE keywords (
  id bigserial primary key,
  contents text,
  sentiment DECIMAL(6) NOT NULL,
  relevance DECIMAL(6) NOT NULL,
  sadness DECIMAL(4) NOT NULL,
  joy DECIMAL(4) NOT NULL,
  fear DECIMAL(4) NOT NULL,
  disgust DECIMAL(4) NOT NULL,
  anger DECIMAL(4) NOT NULL,
  watson_id integer  REFERENCES watson_results(id) NOT NULL
);
