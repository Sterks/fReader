package services

import (
	"database/sql"

	"github.com/Sterks/fReader/config"
	"github.com/Sterks/fReader/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/lib/pq"
)

// Migrations Структура для миграций
type Migrations struct {
	DBMigration *sql.DB
	Config      *config.Config
	Logger      *logger.Logger
}

// NewMigration ... Новая структура для миграций
func NewMigration(config *config.Config) *Migrations {
	return &Migrations{
		DBMigration: &sql.DB{},
		Config:      config,
		Logger:      logger.NewLogger(),
	}
}

// ConfigureMigration Настройка миграций
func (m *Migrations) ConfigureMigration(config *config.Config) {
	// m.Logger.ConfigureLogger(config)
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Postgres.Host, config.Postgres.Port, config.Postgres.User, config.Postgres.Password, config.Postgres.DBName)
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// m.Logger.ErrorLog("Не могу подключится к базе для выполения миграций", err)
	mig, err := migrate.New(
		"github://sterks:0ab1fa3e589b03243854163d4a84692a5dc61f52@Sterks/fReader/migrations",
		"postgres://user_ro:4r2w3e1q@Postgres:5432/freader?sslmode=disable",
		// psqlInfo,
	)
	if err != nil {
		m.Logger.ErrorLog("Не могу подключиться к базе", err)
	}
	if err2 := mig.Up(); err2 != nil {

	}
}
