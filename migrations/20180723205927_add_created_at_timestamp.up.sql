ALTER TABLE ONLY google_results ALTER COLUMN created_at SET DEFAULT now();
ALTER TABLE ONLY watson_results ALTER COLUMN created_at SET DEFAULT now();
