package postgresql

import (
	"fmt"
	"os"
	"testing"

	"github.com/tecsisa/authorizr/database"
	"time"
)

var repoDB PostgresRepo

func TestMain(m *testing.M) {
	// Wait for DB
	time.Sleep(3 * time.Second)
	// Retrieve db connector to run test
	dbmap, err := InitDb("postgres://postgres:password@localhost:54320/postgres?sslmode=disable")
	if err != nil {
		fmt.Fprintln(os.Stderr, "There was an error starting connector", err)
		os.Exit(1)
	}
	repoDB = PostgresRepo{
		Dbmap: dbmap,
	}

	result := m.Run()

	os.Exit(result)
}

// Aux methods

func insertUser(id string, externalID string, path string, createAt int64, urn string) error {
	err := repoDB.Dbmap.Exec("INSERT INTO public.users (id, external_id, path, create_at, urn) VALUES (?, ?, ?, ?, ?)",
		id, externalID, path, createAt, urn).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}
	return nil
}

func getUsersCountFiltered(id string, externalID string, path string, createAt int64, urn string) (int, error) {
	query := repoDB.Dbmap.Table(User{}.TableName())
	if id != "" {
		query = query.Where("id = ?", id)
	}
	if externalID != "" {
		query = query.Where("external_id = ?", externalID)
	}
	if path != "" {
		query = query.Where("path = ?", path)
	}
	if createAt != 0 {
		query = query.Where("create_at = ?", createAt)
	}
	if urn != "" {
		query = query.Where("urn = ?", urn)
	}
	var number int
	if err := query.Count(&number).Error; err != nil {
		return 0, err
	}

	return number, nil
}

func cleanUserTable() error {
	if err := repoDB.Dbmap.Delete(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// GROUP

func insertGroup(id string, name string, path string, createAt int64, urn string, org string) error {
	err := repoDB.Dbmap.Exec("INSERT INTO public.groups (id, name, path, create_at, urn, org) VALUES (?, ?, ?, ?, ?, ?)",
		id, name, path, createAt, urn, org).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}
	return nil
}

func getGroupsCountFiltered(id string, name string, path string, createAt int64, urn string, org string) (int, error) {
	query := repoDB.Dbmap.Table(Group{}.TableName())
	if id != "" {
		query = query.Where("id = ?", id)
	}
	if name != "" {
		query = query.Where("name = ?", name)
	}
	if path != "" {
		query = query.Where("path = ?", path)
	}
	if createAt != 0 {
		query = query.Where("create_at = ?", createAt)
	}
	if urn != "" {
		query = query.Where("urn = ?", urn)
	}
	if org != "" {
		query = query.Where("org = ?", org)
	}
	var number int
	if err := query.Count(&number).Error; err != nil {
		return 0, err
	}

	return number, nil
}

func cleanGroupTable() error {
	if err := repoDB.Dbmap.Delete(&Group{}).Error; err != nil {
		return err
	}
	return nil
}