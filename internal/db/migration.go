package database

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/tmozzze/SubChecker/internal/config"
)

func RunMigration(ctx context.Context, pool *pgxpool.Pool, cfg *config.Config, log *logrus.Logger) error {
	log.Info("Starting migrations...")
	migrationsDir := cfg.MigrationsDir
	if migrationsDir == "" {
		migrationsDir = "./database/migrations"
	}

	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Игнорируем директории и файлы, которые не заканчиваются на .up.sql
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".up.sql") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", d.Name(), err)
		}

		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", d.Name(), err)
		}

		log.Infof("Migration %s applied", d.Name())
		return nil
	})

	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
