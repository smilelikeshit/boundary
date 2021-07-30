// +build linux darwin windows

package dbtemplate

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	// postgres dialect
	_ "github.com/lib/pq"
)

func init() {
	StartDbFromTemplate = startDbFromTemplateSupported
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func startDbFromTemplateSupported(dialect string, opt ...Option) (cleanup func() error, retURL, container string, err error) {
	db, err := sql.Open("postgres", "postgres://boundary:boundary@127.0.0.1/boundary?sslmode=disable")
	if err != nil {
		return func() error { return nil }, "", "", fmt.Errorf("could not connect to source database: %w", err)
	}
	defer db.Close()

	dbname := fmt.Sprintf("boundary_test_%s", randStr(16))

	_, err = db.Exec(fmt.Sprintf("create database %s template boundary_template", dbname))
	if err != nil {
		return func() error { return nil }, "", "", fmt.Errorf("could not create test database: %w", err)
	}

	url := fmt.Sprintf("postgres://boundary:boundary@127.0.0.1/%s?sslmode=disable", dbname)

	cleanup = func() error {
		return dropDatabase(dbname)
	}

	tdb, err := sql.Open(dialect, url)
	if err != nil {
		return func() error { return nil }, "", "", fmt.Errorf("could not create test database: %w", err)
	}
	defer tdb.Close()

	if err := tdb.Ping(); err != nil {
		return func() error { return nil }, "", "", fmt.Errorf("could not ping test database: %w", err)
	}

	return cleanup, url, "", nil
}

const killconns = `
select
  pg_terminate_backend(pg_stat_activity.pid)
from
  pg_stat_activity
where pg_stat_activity.datname = $1
  and pid <> pg_backend_pid();
`

func dropDatabase(dbname string) error {
	db, err := sql.Open("postgres", "postgres://boundary:boundary@127.0.0.1/boundary?sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(killconns, dbname)
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("drop database %s", dbname))
	if err != nil {
		return err
	}
	return nil
}
