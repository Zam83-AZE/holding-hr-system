package repository

import (
        "database/sql"
        "fmt"
        "time"

        "holding-hr-system/config"

        _ "github.com/go-sql-driver/mysql"
)

func NewDB(cfg *config.Config) (*sql.DB, error) {
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8mb4&multiStatements=true",
                cfg.DBUser,
                cfg.DBPassword,
                cfg.DBHost,
                cfg.DBPort,
                cfg.DBName,
        )

        db, err := sql.Open("mysql", dsn)
        if err != nil {
                return nil, fmt.Errorf("failed to open database: %w", err)
        }

        // Connection pool settings
        db.SetMaxOpenConns(25)
        db.SetMaxIdleConns(5)
        db.SetConnMaxLifetime(5 * time.Minute)

        // Test connection
        if err := db.Ping(); err != nil {
                return nil, fmt.Errorf("failed to ping database: %w", err)
        }

        return db, nil
}
