package main
package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	_ "github.com/lib/pq"
)

func TestPostgresConnectionWithDocker(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15-alpine",
		Env:        []string{"POSTGRES_USER=test", "POSTGRES_PASSWORD=test", "POSTGRES_DB=testdb"},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}
	defer pool.Purge(resource)

	host := "localhost"
	port := resource.GetPort("5432/tcp")
	dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port)

	var db *sql.DB
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		t.Fatalf("Could not connect to dockerized postgres: %s", err)
	}
	defer db.Close()

	_ = os.Setenv("DB_HOST", host)
	_ = os.Setenv("DB_PORT", port)
	_ = os.Setenv("DB_USER", "test")
	_ = os.Setenv("DB_PASSWORD", "test")
	_ = os.Setenv("DB_NAME", "testdb")

	// Simple sanity query
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS health(id INT PRIMARY KEY);")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	time.Sleep(1 * time.Second)
}
