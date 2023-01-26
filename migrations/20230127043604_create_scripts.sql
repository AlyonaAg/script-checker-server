-- +goose Up
-- +goose StatementBegin
CREATE TABLE scripts (
    id SERIAL PRIMARY KEY,
    url TEXT,
    original_script TEXT UNIQUE NOT NULL,
    deobf_script TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE scripts;
-- +goose StatementEnd
