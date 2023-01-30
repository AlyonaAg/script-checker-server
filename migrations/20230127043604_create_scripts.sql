-- +goose Up
-- +goose StatementBegin
CREATE TABLE scripts (
    id SERIAL PRIMARY KEY,
    url TEXT,
    original_script TEXT NOT NULL,
    result BOOLEAN,
    danger_percent FLOAT,
    virus_total TEXT,

    UNIQUE (url, original_script)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE scripts;
-- +goose StatementEnd
