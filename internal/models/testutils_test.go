package models

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "test_web:pass@/test_pastepanda?parseTime=true&multiStatements=true")
	require.NoError(t, err)

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer db.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		require.NoError(t, err)

		_, err = db.Exec(string(script))
		require.NoError(t, err)
	})

	return db
}
