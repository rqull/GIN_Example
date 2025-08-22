package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "require" // default aman untuk Railway
		}
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", url.PathEscape(user), url.PathEscape(pass), host, port, name, sslmode)
	} else {
		// pastikan sslmode ada
		u, err := url.Parse(dsn)
		if err == nil {
			q := u.Query()
			if q.Get("sslmode") == "" {
				q.Set("sslmode", "require")
				u.RawQuery = q.Encode()
				dsn = u.String()
			}
		}
	}

	db, err := sql.Open("pgx", dsn) // sesuaikan nama driver
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}
