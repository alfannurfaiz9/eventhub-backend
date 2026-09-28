package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PsqlDb struct {
	User   string
	Pass   string
	Host   string
	Port   string
	DBName string
}

func NewPsqlDb(user, pass, host, port, dbname string) *PsqlDb {
	return &PsqlDb{
		User:   user,
		Pass:   pass,
		Host:   host,
		Port:   port,
		DBName: dbname,
	}
}

func (p *PsqlDb) Connect() (*pgxpool.Pool, error) {
	// postgres://username:password@localhost:5432/database_name
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", p.User, p.Pass, p.Host, p.Port, p.DBName)

	return pgxpool.New(context.Background(), connStr)
}
