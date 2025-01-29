package engines

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mburaksoran/insider-case/internal/app/config"
)

type SqlDbEngine struct {
	Client *sql.DB
}

var dbEngine *SqlDbEngine

func GetSqlDbEngine() *SqlDbEngine {
	return dbEngine
}

func SetSqlDBEngine(cfg *config.AppConfig) (*SqlDbEngine, error) {
	if dbEngine == nil {
		dbEngine = new(SqlDbEngine)
		ConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.Postgres.SqlUser,
			cfg.Postgres.SqlPassword,
			cfg.Postgres.SqlHost,
			cfg.Postgres.SqlPort,
			cfg.Postgres.SqlDatabaseName,
			cfg.Postgres.SqlSslMode)

		db, err := sql.Open("postgres", ConnectionString)
		if err != nil {
			return nil, err
		}
		dbEngine.Client = db
	}
	return dbEngine, nil
}
