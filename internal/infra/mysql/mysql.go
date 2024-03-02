package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"snappfood/internal/config"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	migrateMySQL "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewClient(ctx context.Context, cfg *config.Config) (*sql.DB, error) {
	loc, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		return nil, err
	}

	c := mysql.Config{
		User:                    cfg.Mysql.UserName,
		Passwd:                  cfg.Mysql.Password,
		DBName:                  cfg.Mysql.Database,
		Net:                     "tcp",
		Addr:                    fmt.Sprintf("%s:%s", cfg.Mysql.Host, cfg.Mysql.Port),
		AllowNativePasswords:    true,
		AllowCleartextPasswords: true,
		ParseTime:               true,
		MultiStatements:         true,
		Loc:                     loc,
	}

	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func NewGormWithInstance(db *sql.DB, debug bool) (*gorm.DB, error) {
	cfg := gorm.Config{}
	if debug {
		cfg.Logger = logger.Default.LogMode(logger.Info)
	} else {
		cfg.Logger = logger.Default.LogMode(logger.Silent)
	}
	gormDB, err := gorm.Open(gormMySQL.New(gormMySQL.Config{
		Conn: db,
	}), &cfg)
	return gormDB, err
}

func Migrate(db *sql.DB) error {
	driver, err := migrateMySQL.WithInstance(db, &migrateMySQL.Config{})
	if err != nil {
		return fmt.Errorf("migration driver failed: %w", err)
	}

	f := "file://migrations"
	m, err := migrate.NewWithDatabaseInstance(f, "mysql", driver)
	if err != nil {
		return fmt.Errorf("failed to init migration: %w", err)
	}
	m.Log = MigrateLogger{}
	log.Print("migration starting")
	start := time.Now()
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return fmt.Errorf("migration failed: %w", err)
		}
		log.Printf("migration no change. duration %v", time.Since(start).Seconds())
		return nil
	}
	log.Printf("migration successful. duration : %v", time.Since(start).Seconds())
	return nil
}

type MigrateLogger struct {
}

func (m MigrateLogger) Printf(format string, v ...interface{}) {
	log.Printf("migration running "+format, v...)
}

func (m MigrateLogger) Verbose() bool {
	return true
}
