package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"employee-management/config"

	_ "github.com/go-sql-driver/mysql"
)

func FindConfigPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := cwd
	for i := 0; i < 6; i++ {
		candidate := filepath.Join(dir, "config.yaml")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("config.yaml not found from %s", cwd)
}

func LoadConfigForTest() error {
	path, err := FindConfigPath()
	if err != nil {
		return err
	}
	return config.InitConfig(path)
}

func EnsureDatabaseExists(cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Charset,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", cfg.DBName))
	return err
}
