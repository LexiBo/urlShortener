package sql

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type StorageInDB struct {
	db *sqlx.DB
}

func New(dbURL string) (*StorageInDB, error) {
	db, err := sqlx.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	s := &StorageInDB{
		db: db,
	}
	return s, nil
}

func (s *StorageInDB) Close() {
	if err := s.db.Close(); err != nil {
		zap.L().Error("failed stop DB: ", zap.Error(err))
	} else {
		zap.L().Info("connection to DB is closed gracefully...")
	}
}
