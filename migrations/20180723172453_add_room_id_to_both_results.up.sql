ALTER TABLE goole_results ADD COLUMN room_id integer REFERENCES rooms(id) NOT NULL;
ALTER TABLE watson_results ADD COLUMN room_id integer REFERENCES rooms(id) NOT NULL;
