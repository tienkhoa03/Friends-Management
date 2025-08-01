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

func TestPostgreSQLSubscriptionRepository_CreateSubscription(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewSubscriptionRepository(gormDB)

	t.Run("successful creation", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "subscriptions"`).WithArgs(requestorId, targetId, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.CreateSubscription(requestorId, targetId)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		requestorId := int64(3)
		targetId := int64(4)

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "subscriptions"`).WithArgs(requestorId, targetId, sqlmock.AnyArg()).WillReturnError(assert.AnError)
		mock.ExpectRollback()

		err := repo.CreateSubscription(requestorId, targetId)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("duplicate subscription constraint violation", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "subscriptions"`).
			WithArgs(requestorId, targetId, sqlmock.AnyArg()).
			WillReturnError(&pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint"})
		mock.ExpectRollback()

		err := repo.CreateSubscription(requestorId, targetId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("invalid user ids", func(t *testing.T) {
		requestorId := int64(-1)
		targetId := int64(0)

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "subscriptions"`).
			WithArgs(requestorId, targetId, sqlmock.AnyArg()).
			WillReturnError(errors.New("invalid input"))
		mock.ExpectRollback()

		err := repo.CreateSubscription(requestorId, targetId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid input")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
func TestPostgreSQLSubscriptionRepository_DeleteSubscription(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewSubscriptionRepository(gormDB)

	t.Run("successful deletion", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "subscriptions"`).WithArgs(requestorId, targetId).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		tx := gormDB.Begin()
		err := repo.DeleteSubscription(tx, requestorId, targetId)
		tx.Commit()
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("deletion error", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "subscriptions"`).
			WithArgs(requestorId, targetId).
			WillReturnError(gorm.ErrInvalidTransaction)
		mock.ExpectRollback()

		tx := gormDB.Begin()
		err := repo.DeleteSubscription(tx, requestorId, targetId)
		tx.Rollback()

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrInvalidTransaction, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
func TestPostgreSQLSubscriptionRepository_GetSubscription(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewSubscriptionRepository(gormDB)

	t.Run("successful getting subscription", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)
		rows := sqlmock.NewRows([]string{"user_id1", "user_id2"}).AddRow(requestorId, targetId)
		mock.ExpectQuery(`SELECT \* FROM "subscriptions"`).
			WithArgs(requestorId, targetId, requestorId).
			WillReturnRows(rows)

		subscription, err := repo.GetSubscription(requestorId, targetId)

		assert.NoError(t, err)
		assert.NotNil(t, subscription)
		assert.Equal(t, requestorId, subscription.RequestorId)
		assert.Equal(t, targetId, subscription.TargetId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("no subscription found", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)
		mock.ExpectQuery(`SELECT \* FROM "subscriptions"`).
			WithArgs(requestorId, targetId, requestorId).
			WillReturnError(gorm.ErrRecordNotFound)

		subscription, err := repo.GetSubscription(requestorId, targetId)

		assert.Error(t, err)
		assert.Nil(t, subscription)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("database error", func(t *testing.T) {
		requestorId := int64(1)
		targetId := int64(2)
		mock.ExpectQuery(`SELECT \* FROM "subscriptions"`).
			WithArgs(requestorId, targetId, requestorId).
			WillReturnError(gorm.ErrInvalidDB)

		subscription, err := repo.GetSubscription(requestorId, targetId)

		assert.Error(t, err)
		assert.Nil(t, subscription)
		assert.Equal(t, gorm.ErrInvalidDB, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgreSQLSubscriptionRepository_GetAllSubscriberIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewSubscriptionRepository(gormDB)

	t.Run("successful retrieval with subscribers", func(t *testing.T) {
		targetId := int64(1)
		expectedSubscribers := []int64{2, 3, 4}

		rows := sqlmock.NewRows([]string{"case"}).
			AddRow(2).
			AddRow(3).
			AddRow(4)

		mock.ExpectQuery(`SELECT "requestor_id" FROM "subscriptions" WHERE target_id = \$1`).
			WithArgs(targetId).
			WillReturnRows(rows)

		subscribers, err := repo.GetAllSubscriberIds(targetId)

		assert.NoError(t, err)
		assert.Equal(t, expectedSubscribers, subscribers)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successful retrieval with no subscribers", func(t *testing.T) {
		targetId := int64(5)

		rows := sqlmock.NewRows([]string{"case"})

		mock.ExpectQuery(`SELECT "requestor_id" FROM "subscriptions" WHERE target_id = \$1`).
			WithArgs(targetId).
			WillReturnRows(rows)

		subscribers, err := repo.GetAllSubscriberIds(targetId)

		assert.NoError(t, err)
		assert.Empty(t, subscribers)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		targetId := int64(1)

		mock.ExpectQuery(`SELECT "requestor_id" FROM "subscriptions" WHERE target_id = \$1`).
			WithArgs(targetId).
			WillReturnError(gorm.ErrInvalidDB)

		subscribers, err := repo.GetAllSubscriberIds(targetId)

		assert.Error(t, err)
		assert.Nil(t, subscribers)
		assert.Equal(t, gorm.ErrInvalidDB, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
