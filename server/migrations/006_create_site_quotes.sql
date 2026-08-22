CREATE TABLE IF NOT EXISTS site_quotes (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  content VARCHAR(500) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT site_quotes_content_not_blank CHECK (length(trim(content)) > 0)
);

INSERT INTO site_quotes (content)
SELECT value FROM site_settings WHERE key = 'home_intro'
  AND NOT EXISTS (SELECT 1 FROM site_quotes);
