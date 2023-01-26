package scripts

import (
	"database/sql"

	"github.com/AlyonaAg/script-checker-server/internal/app/config"
	"github.com/AlyonaAg/script-checker-server/internal/app/model"
	_ "github.com/lib/pq"
)

type Repository interface {
	CreateScript(s model.Scripts) error
}

type repo struct{
	db  *sql.DB
}

func (r *repo) CreateScript(script model.Scripts) error {
	for _, s := range script {
		if err := r.db.QueryRow(
			`INSERT INTO "scripts" (url, script) VALUES ($1, $2) RETURNING id`,
			s.URL, s.Script).Scan(&s.ID); err != nil {
			return err
		}
	}
	return nil
}

func (r *repo) ListScripts(filter ListScriptsFilter) (model.Scripts, error) {
	rows, err := r.db.Query(`SELECT id, url, script FROM "scripts" LIMIT $1 OFFSET $2`,
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