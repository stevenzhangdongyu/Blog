CREATE TABLE IF NOT EXISTS site_settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO site_settings (key, value)
VALUES ('home_intro', '只有一种成功，那就是按照自己的意愿过完一生。')
ON CONFLICT (key) DO NOTHING;
