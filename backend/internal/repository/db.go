package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

type Repository struct {
	DB         *sql.DB
	isPostgres bool
}

func rebind(query string) string {
	var b strings.Builder
	n := 1
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			fmt.Fprintf(&b, "$%d", n)
			n++
		} else {
			b.WriteByte(query[i])
		}
	}
	return b.String()
}

func (r *Repository) Exec(query string, args ...interface{}) (sql.Result, error) {
	if r.isPostgres {
		query = rebind(query)
	}
	return r.DB.Exec(query, args...)
}

func (r *Repository) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if r.isPostgres {
		query = rebind(query)
	}
	return r.DB.Query(query, args...)
}

func (r *Repository) QueryRow(query string, args ...interface{}) *sql.Row {
	if r.isPostgres {
		query = rebind(query)
	}
	return r.DB.QueryRow(query, args...)
}

func (r *Repository) txExec(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	if r.isPostgres {
		query = rebind(query)
	}
	return tx.Exec(query, args...)
}

func (r *Repository) txQueryRow(tx *sql.Tx, query string, args ...interface{}) *sql.Row {
	if r.isPostgres {
		query = rebind(query)
	}
	return tx.QueryRow(query, args...)
}

func (r *Repository) InsertReturningID(query string, args ...interface{}) (int64, error) {
	if r.isPostgres {
		query = rebind(query) + " RETURNING id"
		var id int64
		err := r.DB.QueryRow(query, args...).Scan(&id)
		return id, err
	}
	res, err := r.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) txInsertReturningID(tx *sql.Tx, query string, args ...interface{}) (int64, error) {
	if r.isPostgres {
		query = rebind(query) + " RETURNING id"
		var id int64
		err := tx.QueryRow(query, args...).Scan(&id)
		return id, err
	}
	res, err := tx.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func New(dbType, dsn string) (*Repository, error) {
	isPostgres := dbType == "postgres"

	driverName := "sqlite"
	if isPostgres {
		driverName = "pgx"
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if !isPostgres {
		pragmas := []string{
			"PRAGMA journal_mode=WAL",
			"PRAGMA foreign_keys=ON",
			"PRAGMA busy_timeout=5000",
		}
		for _, p := range pragmas {
			if _, err := db.Exec(p); err != nil {
				return nil, fmt.Errorf("executing %s: %w", p, err)
			}
		}
	}

	if isPostgres {
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
	}

	return &Repository{DB: db, isPostgres: isPostgres}, nil
}

func (r *Repository) RunMigrations(migrationsDir string) error {
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("reading migrations: %w", err)
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("reading %s: %w", f, err)
		}
		if _, err := r.DB.Exec(string(data)); err != nil {
			return fmt.Errorf("executing %s: %w", f, err)
		}
		log.Printf("Migration applied: %s", filepath.Base(f))
	}
	return nil
}

func (r *Repository) RunSeed(seedFile string) error {
	data, err := os.ReadFile(seedFile)
	if err != nil {
		return fmt.Errorf("reading seed file: %w", err)
	}
	if _, err := r.DB.Exec(string(data)); err != nil {
		return fmt.Errorf("executing seed: %w", err)
	}
	log.Println("Seed data applied")
	return nil
}

func (r *Repository) IsEmpty() (bool, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM categorias").Scan(&count)
	if err != nil {
		return true, nil
	}
	return count == 0, nil
}

func (r *Repository) Close() error {
	return r.DB.Close()
}
