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

func TestPostgreSQLFriendshipRepository_CreateFriendship(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewFriendshipRepository(gormDB)

	t.Run("successful creation", func(t *testing.T) {
		userId1 := int64(1)
		userId2 := int64(2)
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "friendships"`).WithArgs(userId1, userId2, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.CreateFriendship(userId1, userId2)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		userId1 := int64(3)
		userId2 := int64(4)

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "friendships"`).WithArgs(userId1, userId2, sqlmock.AnyArg()).WillReturnError(assert.AnError)
		mock.ExpectRollback()

		err := repo.CreateFriendship(userId1, userId2)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("duplicate friendship constraint violation", func(t *testing.T) {
		userId1 := int64(1)
		userId2 := int64(2)

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "friendships"`).
			WithArgs(userId1, userId2, sqlmock.AnyArg()).
			WillReturnError(&pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint"})
		mock.ExpectRollback()

		err := repo.CreateFriendship(userId1, userId2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("invalid user ids", func(t *testing.T) {
		userId1 := int64(-1)
		userId2 := int64(0)

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "friendships"`).
			WithArgs(userId1, userId2, sqlmock.AnyArg()).
			WillReturnError(errors.New("invalid input"))
		mock.ExpectRollback()

		err := repo.CreateFriendship(userId1, userId2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid input")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestPostgreSQLFriendshipRepository_RetrieveFriendIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewFriendshipRepository(gormDB)

	t.Run("successful retrieval with friends", func(t *testing.T) {
		userId := int64(1)
		expectedFriends := []int64{2, 3, 4}

		rows := sqlmock.NewRows([]string{"case"}).
			AddRow(2).
			AddRow(3).
			AddRow(4)

		mock.ExpectQuery(`SELECT CASE WHEN user_id1 = \$1 THEN user_id2 ELSE user_id1 END FROM "friendships" WHERE user_id1 = \$2 OR user_id2 = \$3`).
			WithArgs(userId, userId, userId).
			WillReturnRows(rows)

		friends, err := repo.RetrieveFriendIds(userId)

		assert.NoError(t, err)
		assert.Equal(t, expectedFriends, friends)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successful retrieval with no friends", func(t *testing.T) {
		userId := int64(5)

		rows := sqlmock.NewRows([]string{"case"})

		mock.ExpectQuery(`SELECT CASE WHEN user_id1 = \$1 THEN user_id2 ELSE user_id1 END FROM "friendships" WHERE user_id1 = \$2 OR user_id2 = \$3`).
			WithArgs(userId, userId, userId).
			WillReturnRows(rows)

		friends, err := repo.RetrieveFriendIds(userId)

		assert.NoError(t, err)
		assert.Empty(t, friends)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		userId := int64(1)

		mock.ExpectQuery(`SELECT CASE WHEN user_id1 = \$1 THEN user_id2 ELSE user_id1 END FROM "friendships" WHERE user_id1 = \$2 OR user_id2 = \$3`).
			WithArgs(userId, userId, userId).
			WillReturnError(gorm.ErrInvalidDB)

		friends, err := repo.RetrieveFriendIds(userId)

		assert.Error(t, err)
		assert.Nil(t, friends)
		assert.Equal(t, gorm.ErrInvalidDB, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgreSQLFriendshipRepository_GetFriendship(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewFriendshipRepository(gormDB)

	t.Run("successful getting friendship", func(t *testing.T) {
		userId1 := int64(1)
		userId2 := int64(2)
		rows := sqlmock.NewRows([]string{"user_id1", "user_id2"}).AddRow(userId1, userId2)
		mock.ExpectQuery(`SELECT \* FROM "friendships"`).
			WithArgs(userId1, userId2, userId1).
			WillReturnRows(rows)

		friendship, err := repo.GetFriendship(userId1, userId2)

		assert.NoError(t, err)
		assert.NotNil(t, friendship)
		assert.Equal(t, userId1, friendship.UserId1)
		assert.Equal(t, userId2, friendship.UserId2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("no friendship found", func(t *testing.T) {
		userId1 := int64(1)
		userId2 := int64(2)
		mock.ExpectQuery(`SELECT \* FROM "friendships"`).
			WithArgs(userId1, userId2, userId1).
			WillReturnError(gorm.ErrRecordNotFound)

		friendship, err := repo.GetFriendship(userId1, userId2)

		assert.Error(t, err)
		assert.Nil(t, friendship)
		assert.Equal(t, err, gorm.ErrRecordNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("database error", func(t *testing.T) {
		userId1 := int64(1)
		userId2 := int64(2)
		mock.ExpectQuery(`SELECT \* FROM "friendships"`).
			WithArgs(userId1, userId2, userId1).
			WillReturnError(gorm.ErrInvalidDB)

		friendship, err := repo.GetFriendship(userId1, userId2)

		assert.Error(t, err)
		assert.Nil(t, friendship)
		assert.Equal(t, err, gorm.ErrInvalidDB)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
