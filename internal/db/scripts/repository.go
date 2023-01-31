package scriptsdb

import (
	"database/sql"
	"fmt"

	"github.com/AlyonaAg/script-checker-server/internal/config"
	"github.com/AlyonaAg/script-checker-server/internal/model"
	_ "github.com/lib/pq"
)

type Repository interface {
	CreateScript(script model.Script) (int64, error)
	GetScript(id int64) (*model.Script, error)
	UpdateResultByID(scriptID int64, result bool) error
	UpdateDangerByID(scriptID int64, dangerPercent float64, vtDanger string) error
	ListScripts(filter model.ListScriptsFilter) (model.Scripts, error)
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

func (r *repo) UpdateDangerByID(scriptID int64, dangerPercent float64, vtDanger string) error {
	if _, err := r.db.Exec(
		`UPDATE "scripts" SET danger_percent = $1, virus_total = $2 WHERE id = $3`, dangerPercent, vtDanger, scriptID); err != nil {
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

func (r *repo) ListScripts(filter model.ListScriptsFilter) (model.Scripts, error) {
	fmt.Println(filter)

	rows, err := r.db.Query(`SELECT id, url, original_script, result, danger_percent, virus_total FROM "scripts" LIMIT $1 OFFSET $2`,
		filter.Limit, filter.Page)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scripts model.Scripts
	for rows.Next() {
		s := model.Script{}
		if err := rows.Scan(
			&s.ID,
			&s.URL,
			&s.Script,
			&s.Result,
			&s.DangerPercent,
			&s.VirusTotal,
		); err != nil {
			continue
		}
		fmt.Println(s)

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
