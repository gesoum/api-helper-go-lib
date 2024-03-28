package migrator

import (
	"api-helper/utils"
	"fmt"
	"gorm.io/gorm"
	"log"
	"path/filepath"
	"sort"
	"time"
)

var (
	createSchemaTemplate  = "CREATE SCHEMA IF NOT EXISTS %s;"
	setSearchPathTemplate = "SET search_path TO %s;"
	migrationTableName    = "migration"
)

type Schema struct {
	Name               string
	DataAccessObjects  []any
	PredefinedDataPath string
}

func (s *Schema) createSchema(tx *gorm.DB) error {
	r := tx.Exec(fmt.Sprintf(createSchemaTemplate, s.Name))
	if r.Error != nil {
		return r.Error
	}

	r = tx.Exec(fmt.Sprintf(setSearchPathTemplate, s.Name))
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Schema) createMigrationTable(tx *gorm.DB) error {
	return tx.Exec(
		fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s.%s (name varchar NOT NULL PRIMARY KEY, apply_time timestamp DEFAULT NOW())",
			s.Name, migrationTableName,
		)).Error
}

func (s *Schema) getAppliedMigrations(tx *gorm.DB) ([]*Migration, error) {
	var appliedMigrations []*Migration
	r := tx.Table(fmt.Sprintf("%s.%s", s.Name, migrationTableName)).Select("*").Scan(&appliedMigrations)
	return appliedMigrations, r.Error
}

func (s *Schema) createTables(tx *gorm.DB) error {
	return tx.AutoMigrate(s.DataAccessObjects...)
}

func (s *Schema) createData(tx *gorm.DB) error {
	migrationFilesPaths, err := getMigrationFiles(s.PredefinedDataPath)
	if err != nil {
		return err
	}

	sort.Strings(migrationFilesPaths)

	appliedMigrations, err := s.getAppliedMigrations(tx)
	if err != nil {
		return err
	}

	var appliedMigrationMap = make(map[string]Migration, len(appliedMigrations))
	for _, item := range appliedMigrations {
		appliedMigrationMap[item.Name] = *item
	}

	for _, migrationFilePath := range migrationFilesPaths {
		migrationName := filepath.Base(migrationFilePath)
		if _, ok := appliedMigrationMap[migrationName]; ok {
			log.Printf("migration %s aplready applied", migrationName)
			continue
		}

		sqlScript, err := utils.ReadFileContent(migrationFilePath)
		if err != nil {
			return err
		}

		r := tx.Exec(sqlScript)
		if r.Error != nil {
			log.Printf("cannot apply script[%s]: %s", migrationName, r.Error.Error())
			return r.Error
		}

		r = tx.
			Table(fmt.Sprintf("%s.%s", s.Name, migrationTableName)).
			Create(&Migration{ApplyTime: time.Now(), Name: migrationName})
		if r.Error != nil {
			log.Printf("cannot create migration record: %s", err.Error())
			return r.Error
		}

		log.Printf("migration %s applied", migrationName)
	}

	return nil
}
