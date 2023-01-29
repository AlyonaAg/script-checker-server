package scriptsdb

import (
	"database/sql"

	"github.com/AlyonaAg/script-checker-server/internal/config"
	"github.com/AlyonaAg/script-checker-server/internal/model"
	_ "github.com/lib/pq"
)

type Repository interface {
	CreateScript(script model.Script) (int64, error)
	GetScript(id int64) (*model.Script, error)
	UpdateResultByID(scriptID int64, result bool) error
	UpdateDangerPercentByID(scriptID int64, dangerPercent float64) error
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

func (r *repo) UpdateResultByID(scriptID int64, result bool) error {
	if _, err := r.db.Exec(
		`UPDATE "scripts" SET result = $1 WHERE id = $2`, result, scriptID); err != nil {
		return err
	}
	return nil
}

func (r *repo) UpdateDangerPercentByID(scriptID int64, dangerPercent float64) error {
	if _, err := r.db.Exec(
		`UPDATE "scripts" SET danger_percent = $1 WHERE id = $2`, dangerPercent, scriptID); err != nil {
		return err
	}
	return nil
}

func (r *repo) GetScript(id int64) (*model.Script, error) {
	var script = &model.Script{}
	if err := r.db.QueryRow(`SELECT id, url, original_script FROM "scripts" WHERE id = $1`, id).Scan(
		&script.ID,
		&script.URL,
		&script.Script,
	); err != nil {
		return nil, err
	}

	return script, nil
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
