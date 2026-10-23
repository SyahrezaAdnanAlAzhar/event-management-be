package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool initializes a pgx connection pool
func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Verify connection
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}

// RunMigrations applies SQL migration files found in migrationsDir
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		// If directory does not exist, return nil or error
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			filePath := filepath.Join(migrationsDir, entry.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
			}

			log.Printf("Executing migration: %s", entry.Name())
			if _, err := pool.Exec(ctx, string(content)); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}
