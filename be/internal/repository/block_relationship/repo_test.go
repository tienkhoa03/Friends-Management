package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestPostgreSQLBlockRelationshipRepository_CreateBlockRelationship(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewBlockRelationshipRepository(gormDB)

	t.Run("successful creation", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "block_relationships"`).WithArgs(requestorId, targetId, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := gormDB.Begin()
		err := repo.CreateBlockRelationship(tx, requestorId, targetId)
		tx.Commit()
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		requestorId := int64(3)
		targetId := int64(4)

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "block_relationships"`).WithArgs(requestorId, targetId, sqlmock.AnyArg()).WillReturnError(assert.AnError)
		mock.ExpectRollback()

		tx := gormDB.Begin()
		err := repo.CreateBlockRelationship(tx, requestorId, targetId)
		tx.Rollback()
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("duplicate block relationship constraint violation", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "block_relationships"`).
			WithArgs(requestorId, targetId, sqlmock.AnyArg()).
			WillReturnError(&pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint"})
		mock.ExpectRollback()

		tx := gormDB.Begin()
		err := repo.CreateBlockRelationship(tx, requestorId, targetId)
		tx.Rollback()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("invalid user ids", func(t *testing.T) {
		requestorId := int64(-1)
		targetId := int64(0)

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "block_relationships"`).
			WithArgs(requestorId, targetId, sqlmock.AnyArg()).
			WillReturnError(errors.New("invalid input"))
		mock.ExpectRollback()

		tx := gormDB.Begin()
		err := repo.CreateBlockRelationship(tx, requestorId, targetId)
		tx.Rollback()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid input")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestPostgreSQLBlockRelationshipRepository_GetBlockRelationship(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewBlockRelationshipRepository(gormDB)

	t.Run("successful getting block relationship", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)
		rows := sqlmock.NewRows([]string{"user_id1", "user_id2"}).AddRow(requestorId, targetId)
		mock.ExpectQuery(`SELECT \* FROM "block_relationships"`).
			WithArgs(requestorId, targetId, requestorId).
			WillReturnRows(rows)

		blockRelationship, err := repo.GetBlockRelationship(requestorId, targetId)

		assert.NoError(t, err)
		assert.NotNil(t, blockRelationship)
		assert.Equal(t, requestorId, blockRelationship.RequestorId)
		assert.Equal(t, targetId, blockRelationship.TargetId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("no block relationship found", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)
		mock.ExpectQuery(`SELECT \* FROM "block_relationships"`).
			WithArgs(requestorId, targetId, requestorId).
			WillReturnError(gorm.ErrRecordNotFound)

		blockRelationship, err := repo.GetBlockRelationship(requestorId, targetId)

		assert.Error(t, err)
		assert.Nil(t, blockRelationship)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("database error", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)
		mock.ExpectQuery(`SELECT \* FROM "block_relationships"`).
			WithArgs(requestorId, targetId, requestorId).
			WillReturnError(gorm.ErrInvalidDB)

		blockRelationship, err := repo.GetBlockRelationship(requestorId, targetId)

		assert.Error(t, err)
		assert.Nil(t, blockRelationship)
		assert.Equal(t, gorm.ErrInvalidDB, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgreSQLBlockRelationshipRepository_GetBlockRequestorIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewBlockRelationshipRepository(gormDB)

	t.Run("successful retrieval with block requestors", func(t *testing.T) {
		targetId := int64(1)
		expectedBlockRequestors := []int64{2, 3, 4}

		rows := sqlmock.NewRows([]string{"case"}).
			AddRow(2).
			AddRow(3).
			AddRow(4)

		mock.ExpectQuery(`SELECT "requestor_id" FROM "block_relationships" WHERE target_id = \$1`).
			WithArgs(targetId).
			WillReturnRows(rows)

		blockRequestors, err := repo.GetBlockRequestorIds(targetId)

		assert.NoError(t, err)
		assert.Equal(t, expectedBlockRequestors, blockRequestors)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successful retrieval with no block requestors", func(t *testing.T) {
		targetId := int64(5)

		rows := sqlmock.NewRows([]string{"case"})

		mock.ExpectQuery(`SELECT "requestor_id" FROM "block_relationships" WHERE target_id = \$1`).
			WithArgs(targetId).
			WillReturnRows(rows)

		blockRequestors, err := repo.GetBlockRequestorIds(targetId)

		assert.NoError(t, err)
		assert.Empty(t, blockRequestors)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		targetId := int64(1)

		mock.ExpectQuery(`SELECT "requestor_id" FROM "block_relationships" WHERE target_id = \$1`).
			WithArgs(targetId).
			WillReturnError(gorm.ErrInvalidDB)

		blockRequestors, err := repo.GetBlockRequestorIds(targetId)

		assert.Error(t, err)
		assert.Nil(t, blockRequestors)
		assert.Equal(t, gorm.ErrInvalidDB, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgreSQLBlockRelationshipRepository_DeleteBlockRelationship(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewBlockRelationshipRepository(gormDB)

	t.Run("successful deletion", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "block_relationships"`).WithArgs(requestorId, targetId).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.DeleteBlockRelationship(requestorId, targetId)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("deletion error", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "block_relationships"`).
			WithArgs(requestorId, targetId).
			WillReturnError(gorm.ErrInvalidTransaction)
		mock.ExpectRollback()

		err := repo.DeleteBlockRelationship(requestorId, targetId)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrInvalidTransaction, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
