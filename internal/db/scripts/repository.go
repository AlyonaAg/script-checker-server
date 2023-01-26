package scriptsdb

import (
	"database/sql"

	"github.com/AlyonaAg/script-checker-server/internal/config"
	"github.com/AlyonaAg/script-checker-server/internal/model"
	_ "github.com/lib/pq"
)

type Repository interface {
	CreateScript(script model.Script) (int64, error)
}

type repo struct {
	db *sql.DB
}

func (r *repo) CreateScript(script model.Script) (int64, error) {
	if err := r.db.QueryRow(
		`INSERT INTO "scripts" (url, original_script) VALUES ($1, $2) RETURNING id`,
		script.URL, script.Script).Scan(&script.ID); err != nil {
		return 0, err
	}
	return script.ID, nil
}

func (r *repo) UpdateScript(scriptID int64, deobfScript string) error {
	if _, err := r.db.Exec(
		`UPDATE "scripts" SET deobf_script = $1 WHERE id = $2`, deobfScript, scriptID); err != nil {
		return err
	}
	return nil
}

func (r *repo) ListScripts(filter ListScriptsFilter) (model.Scripts, error) {
	rows, err := r.db.Query(`SELECT id, url, original_script FROM "scripts" LIMIT $1 OFFSET $2`,
		filter.Limit, filter.Page)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scripts model.Scripts
	for rows.Next() {
		s := model.Script{}
		if err := rows.Scan(
			s.ID,
			s.URL,
			s.Script,
		); err != nil {
			continue
		}

		scripts = append(scripts, &s)
	}
	return scripts, nil
}

func NewRepository() (Repository, error) {
	databaseURL, err := config.GetValue(config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", databaseURL.(string))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &repo{
		db: db,
	}, nil
}
