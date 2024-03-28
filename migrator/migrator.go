package migrator

import (
	"gorm.io/gorm"
	"time"
)

func New(dbConnection *gorm.DB, schemas []Schema) migrator {
	return migrator{
		dbConnection: dbConnection,
		schemas:      schemas,
	}
}

type Migration struct {
	ApplyTime time.Time
	Name      string
}

type migrator struct {
	dbConnection *gorm.DB
	schemas      []Schema
}

func (m *migrator) Migrate() []error {
	var err error
	var errors []error

	for _, schema := range m.schemas {
		if err != nil {
			errors = append(errors, err)
		}
		err = nil

		tx := m.dbConnection.Begin()

		err = schema.createSchema(tx)
		if err != nil {
			tx.Rollback()
			continue
		}

		err = schema.createTables(tx)
		if err != nil {
			tx.Rollback()
			continue
		}

		err = schema.createMigrationTable(tx)
		if err != nil {
			tx.Rollback()
			continue
		}

		err = schema.createData(tx)
		if err != nil {
			tx.Rollback()
			continue
		}

		tx.Commit()
	}

	return errors
}
