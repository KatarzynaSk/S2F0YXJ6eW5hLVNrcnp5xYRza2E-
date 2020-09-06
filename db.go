package main

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type request struct {
	ID       int    `json:"id" db:"id"`
	URL      string `json:"url" db:"url"`
	Interval int    `json:"interval" db:"interval"`
	Active   bool   `json:"-" db:"active"`
}

type requestResult struct {
	ID        int       `json:"id" db:"id"`
	Response  *string   `json:"response" db:"response"`
	Duration  float64   `json:"duration" db:"duration"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	RequestID int       `json:"-" db:"request_id"`
}

func insertRequest(db *sqlx.DB, e request) (int, error) {
	var id int
	query := `INSERT INTO request (url, interval) 
		VALUES (:url, :interval) RETURNING id`

	rows, err := db.NamedQuery(query, e)
	if err != nil {
		return 0, err
	}

	if rows.Next() {
		err = rows.Scan(&id)
	}

	if err != nil {
		return 0, err
	}

	return id, nil

}

func deactivateRequest(db *sqlx.DB, id int) error {
	query := `UPDATE request 
	SET active=false
	WHERE id=$1`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func selectAllActiveRequest(db *sqlx.DB) ([]request, error) {
	query := `SELECT * FROM request 
	WHERE active=true`

	all := make([]request, 0, 0)
	err := db.Select(&all, query)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func selectActiveRequest(db *sqlx.DB, requestID int) (*request, error) {
	query := `SELECT * FROM request 
	WHERE request.id=$1 AND active=true`

	var r request
	err := db.Get(&r, query, requestID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, nil
		}

		return nil, err
	}

	return &r, nil
}

func insertRequestResult(db *sqlx.DB, r requestResult) error {
	query := `INSERT INTO request_result (response, duration, request_id) 
		VALUES (:response, :duration, :request_id)`

	_, err := db.NamedExec(query, r)
	if err != nil {
		return err
	}

	return nil
}

func selectAllRequestResult(db *sqlx.DB, requestID int) ([]requestResult, error) {
	query := `SELECT * FROM request_result WHERE request_id=$1`

	all := make([]requestResult, 0, 0)
	err := db.Select(&all, query, requestID)
	if err != nil {
		return nil, err
	}

	return all, nil
}
