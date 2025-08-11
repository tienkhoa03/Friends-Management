package repository

import (
	"BE_Friends_Management/internal/domain/entity"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestPostgreSQLUserRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(gormDB)

	t.Run("successful creation", func(t *testing.T) {
		userId := int64(1)
		userEmail := "example@gmail.com"
		userPassword := "123"
		userRole := "user"
		user := &entity.User{Email: userEmail, Password: userPassword, Role: userRole}
		rows := sqlmock.NewRows([]string{"id", "email", "password", "role"}).AddRow(userId, userEmail, userPassword, userRole)
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "users"`).WithArgs(userEmail, userPassword, userRole, sqlmock.AnyArg()).WillReturnRows(rows)
		mock.ExpectCommit()

		createdUser, err := repo.CreateUser(user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, userId, createdUser.Id)
		assert.Equal(t, userEmail, createdUser.Email)
		assert.Equal(t, userRole, createdUser.Role)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		userEmail := "example@gmail.com"
		userPassword := "123"
		userRole := "user"
		user := &entity.User{Email: userEmail, Password: userPassword, Role: userRole}
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "users"`).WithArgs(userEmail, userPassword, userRole, sqlmock.AnyArg()).WillReturnError(assert.AnError)
		mock.ExpectRollback()

		createdUser, err := repo.CreateUser(user)
		assert.Error(t, err)
		assert.Nil(t, createdUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate subscription constraint violation", func(t *testing.T) {
		userEmail := "example@gmail.com"
		userPassword := "123"
		userRole := "user"
		user := &entity.User{Email: userEmail, Password: userPassword, Role: userRole}
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "users"`).WithArgs(userEmail, userPassword, userRole, sqlmock.AnyArg()).WillReturnError(&pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint"})
		mock.ExpectRollback()

		createdUser, err := repo.CreateUser(user)
		assert.Error(t, err)
		assert.Nil(t, createdUser)
		assert.Contains(t, err.Error(), "duplicate key")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgreSQLUserRepository_GetAllUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(gormDB)

	t.Run("successful retrieval with users", func(t *testing.T) {
		createdAt1 := time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC)
		createdAt2 := time.Date(2025, 8, 2, 10, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "created_at"}).
			AddRow(1, "user1@example.com", "123", "user", createdAt1).
			AddRow(2, "user2@example.com", "123", "user", createdAt2)

		mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(rows)

		users, err := repo.GetAllUser()

		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, int64(1), users[0].Id)
		assert.Equal(t, "user1@example.com", users[0].Email)
		assert.Equal(t, "123", users[0].Password)
		assert.Equal(t, "user", users[0].Role)
		assert.Equal(t, createdAt1, users[0].CreatedAt)
		assert.Equal(t, int64(2), users[1].Id)
		assert.Equal(t, "user2@example.com", users[1].Email)
		assert.Equal(t, "123", users[1].Password)
		assert.Equal(t, "user", users[1].Role)
		assert.Equal(t, createdAt2, users[1].CreatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successful retrieval with no users", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "created_at"})

		mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(rows)

		users, err := repo.GetAllUser()

		assert.NoError(t, err)
		assert.Len(t, users, 0)
		assert.Empty(t, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error retrieval", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnError(assert.AnError)

		users, err := repo.GetAllUser()

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgreSQLUserRepository_GetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(gormDB)

	t.Run("successful getting user", func(t *testing.T) {
		userId := int64(1)
		createdAt := time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "created_at"}).
			AddRow(1, "user1@example.com", "123", "user", createdAt)
		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(userId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		user, err := repo.GetUserById(userId)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userId, user.Id)
		assert.Equal(t, "user1@example.com", user.Email)
		assert.Equal(t, "123", user.Password)
		assert.Equal(t, "user", user.Role)
		assert.Equal(t, createdAt, user.CreatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no user found", func(t *testing.T) {
		userId := int64(1)
		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(userId, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.GetUserById(userId)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, err, gorm.ErrRecordNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("database error", func(t *testing.T) {
		userId := int64(1)
		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(userId, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrInvalidDB)

		user, err := repo.GetUserById(userId)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, err, gorm.ErrInvalidDB)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgreSQLUserRepository_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(gormDB)

	t.Run("successful getting user", func(t *testing.T) {
		email := "user1@example.com"
		createdAt := time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "created_at"}).
			AddRow(1, "user1@example.com", "123", "user", createdAt)
		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(email, sqlmock.AnyArg()).
			WillReturnRows(rows)

		user, err := repo.GetUserByEmail(email)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(1), user.Id)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, "123", user.Password)
		assert.Equal(t, "user", user.Role)
		assert.Equal(t, createdAt, user.CreatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no user found", func(t *testing.T) {
		email := "user1@example.com"
		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(email, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.GetUserByEmail(email)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, err, gorm.ErrRecordNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("database error", func(t *testing.T) {
		email := "user1@example.com"
		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(email, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrInvalidDB)

		user, err := repo.GetUserByEmail(email)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, err, gorm.ErrInvalidDB)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgreSQLUserRepository_GetUsersFromIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(gormDB)

	t.Run("successful getting users", func(t *testing.T) {
		userIds := []int64{1, 2, 3}
		createdAt1 := time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC)
		createdAt2 := time.Date(2025, 8, 2, 10, 0, 0, 0, time.UTC)
		createdAt3 := time.Date(2025, 8, 3, 10, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "created_at"}).
			AddRow(1, "user1@example.com", "123", "user", createdAt1).
			AddRow(2, "user2@example.com", "123", "user", createdAt2).
			AddRow(3, "user3@example.com", "123", "user", createdAt3)
		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(userIds[0], userIds[1], userIds[2]).
			WillReturnRows(rows)

		users, err := repo.GetUsersFromIds(userIds)

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Equal(t, userIds[0], users[0].Id)
		assert.Equal(t, "user1@example.com", users[0].Email)
		assert.Equal(t, "123", users[0].Password)
		assert.Equal(t, "user", users[0].Role)
		assert.Equal(t, createdAt1, users[0].CreatedAt)
		assert.Equal(t, userIds[1], users[1].Id)
		assert.Equal(t, "user2@example.com", users[1].Email)
		assert.Equal(t, "123", users[1].Password)
		assert.Equal(t, "user", users[1].Role)
		assert.Equal(t, createdAt2, users[1].CreatedAt)
		assert.Equal(t, userIds[2], users[2].Id)
		assert.Equal(t, "user3@example.com", users[2].Email)
		assert.Equal(t, "123", users[2].Password)
		assert.Equal(t, "user", users[2].Role)
		assert.Equal(t, createdAt3, users[2].CreatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no users found", func(t *testing.T) {
		userIds := []int64{1, 2, 3}
		rows := sqlmock.NewRows([]string{"id", "email", "created_at"})

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(userIds[0], userIds[1], userIds[2]).
			WillReturnRows(rows)

		users, err := repo.GetUsersFromIds(userIds)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Empty(t, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("database error", func(t *testing.T) {
		userIds := []int64{1, 2, 3}
		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(userIds[0], userIds[1], userIds[2]).
			WillReturnError(gorm.ErrInvalidDB)

		user, err := repo.GetUsersFromIds(userIds)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, err, gorm.ErrInvalidDB)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgreSQLUserRepository_GetUserFromEmails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(gormDB)

	t.Run("successful getting users", func(t *testing.T) {
		userEmails := []string{"user1@example.com", "user2@example.com", "user3@example.com"}
		createdAt1 := time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC)
		createdAt2 := time.Date(2025, 8, 2, 10, 0, 0, 0, time.UTC)
		createdAt3 := time.Date(2025, 8, 3, 10, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "created_at"}).
			AddRow(1, "user1@example.com", "123", "user", createdAt1).
			AddRow(2, "user2@example.com", "123", "user", createdAt2).
			AddRow(3, "user3@example.com", "123", "user", createdAt3)
		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(userEmails[0], userEmails[1], userEmails[2]).
			WillReturnRows(rows)

		users, err := repo.GetUsersFromEmails(userEmails)

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Equal(t, int64(1), users[0].Id)
		assert.Equal(t, "user1@example.com", users[0].Email)
		assert.Equal(t, "123", users[0].Password)
		assert.Equal(t, "user", users[0].Role)
		assert.Equal(t, createdAt1, users[0].CreatedAt)
		assert.Equal(t, int64(2), users[1].Id)
		assert.Equal(t, "user2@example.com", users[1].Email)
		assert.Equal(t, "123", users[1].Password)
		assert.Equal(t, "user", users[1].Role)
		assert.Equal(t, createdAt2, users[1].CreatedAt)
		assert.Equal(t, int64(3), users[2].Id)
		assert.Equal(t, "user3@example.com", users[2].Email)
		assert.Equal(t, "123", users[2].Password)
		assert.Equal(t, "user", users[2].Role)
		assert.Equal(t, createdAt3, users[2].CreatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no users found", func(t *testing.T) {
		userEmails := []string{"user1@example.com", "user2@example.com", "user3@example.com"}
		rows := sqlmock.NewRows([]string{"id", "email", "created_at"})

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(userEmails[0], userEmails[1], userEmails[2]).
			WillReturnRows(rows)

		users, err := repo.GetUsersFromEmails(userEmails)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Empty(t, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("database error", func(t *testing.T) {
		userEmails := []string{"user1@example.com", "user2@example.com", "user3@example.com"}
		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(userEmails[0], userEmails[1], userEmails[2]).
			WillReturnError(gorm.ErrInvalidDB)

		user, err := repo.GetUsersFromEmails(userEmails)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, err, gorm.ErrInvalidDB)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgreSQLUserRepository_UpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(gormDB)

	t.Run("successful updating user", func(t *testing.T) {
		userId := int64(1)
		userEmail := "user1@example.com"
		userPassword := "123"
		userRole := "user"
		createdAt := time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC)
		user := &entity.User{Id: userId, Email: userEmail, Password: userPassword, Role: userRole, CreatedAt: createdAt}
		rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "created_at"}).
			AddRow(userId, userEmail, userPassword, userRole, createdAt)
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users"`).
			WithArgs(userId, userEmail, userPassword, userRole, createdAt, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(userId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		updatedUser, err := repo.UpdateUser(user)

		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		assert.Equal(t, userId, updatedUser.Id)
		assert.Equal(t, userEmail, updatedUser.Email)
		assert.Equal(t, userPassword, updatedUser.Password)
		assert.Equal(t, userRole, updatedUser.Role)
		assert.Equal(t, createdAt, updatedUser.CreatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error user not found", func(t *testing.T) {
		userId := int64(1)
		userEmail := "user1@example.com"
		createdAt := time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC)
		user := &entity.User{Id: userId, Email: userEmail, CreatedAt: createdAt}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users"`).
			WithArgs(userId, userEmail, createdAt, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(userId, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		updatedUser, err := repo.UpdateUser(user)

		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		userId := int64(1)
		userEmail := "user1@example.com"
		createdAt := time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC)
		user := &entity.User{Id: userId, Email: userEmail, CreatedAt: createdAt}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users"`).
			WithArgs(userId, userEmail, createdAt, sqlmock.AnyArg()).
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		updatedUser, err := repo.UpdateUser(user)

		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgreSQLUserRepository_DeleteUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(gormDB)

	t.Run("successful deletion", func(t *testing.T) {
		userId := int64(1)
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "users"`).WithArgs(userId).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.DeleteUserById(userId)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("deletion error", func(t *testing.T) {
		userId := int64(1)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "users"`).
			WithArgs(userId).
			WillReturnError(gorm.ErrInvalidTransaction)
		mock.ExpectRollback()

		err := repo.DeleteUserById(userId)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrInvalidTransaction, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
