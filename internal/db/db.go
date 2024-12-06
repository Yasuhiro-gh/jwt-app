package db

import (
	"database/sql"
	"github.com/Yasuhiro-gh/jwt-app/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB() *PostgresDB {
	return &PostgresDB{}
}

func CreateTable(pdb *PostgresDB) error {
	_, err := pdb.DB.Exec(`CREATE TABLE IF NOT EXISTS tokens ("user_id" UUID UNIQUE, "refresh_token" VARCHAR(255) NOT NULL UNIQUE)`)
	return err
}

func (pdb *PostgresDB) SetNewToken(userID, refreshToken string) error {
	_, err := pdb.DB.Exec("INSERT INTO tokens (user_id, refresh_token) VALUES ($1, $2)", userID, refreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (pdb *PostgresDB) RefreshToken(refreshToken string) error {
	return nil
}

func (pdb *PostgresDB) GetTokenByUserID(userID string) (string, error) {
	qr := pdb.DB.QueryRow("SELECT refresh_token FROM tokens WHERE user_id=$1", userID)
	if qr != nil && qr.Err() != nil {
		return "", qr.Err()
	}
	var UID string
	err := qr.Scan(&UID)
	if err != nil {
		return "", err
	}
	return UID, nil
}

func (pdb *PostgresDB) OpenConnection() error {
	db, err := sql.Open("pgx", config.Options.DatabaseDSN)
	if err != nil {
		return err
	}
	pdb.DB = db
	return nil
}

func (pdb *PostgresDB) CloseConnection() error {
	return pdb.DB.Close()
}
