package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	utils "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	errUserNotExists    = errors.New("user does not exist")
	errUserExists       = errors.New("user already exists")
	errSessionNotExists = errors.New("session does not exist")
	NameSeqUser         = pgx.Identifier{"public", "user_id_seq"} //nolint:gochecknoglobals
)

type UsersListWrapper struct {
	UsersList *models.UsersList
	Pool      *pgxpool.Pool
}

func (active *UsersListWrapper) userExists(ctx context.Context, tx pgx.Tx, email string) (bool, error) {
	SQLUserExists := `SELECT EXISTS(SELECT 1 FROM public."user" WHERE email=$1 );`
	userLine := tx.QueryRow(ctx, SQLUserExists, email)

	var exists bool

	if err := userLine.Scan(&exists); err != nil {

		return false, err
	}

	return exists, nil
}

func (active *UsersListWrapper) UserExists(ctx context.Context, email string) bool {
	var exists bool

	_ = pgx.BeginFunc(ctx, active.Pool, func(tx pgx.Tx) error {
		userExists, err := active.userExists(ctx, tx, email)
		exists = userExists

		return err
	})

	return exists
}

func (active *UsersListWrapper) GetLastID() uint {
	active.UsersList.UsersCount++

	return active.UsersList.UsersCount
}

func (active *UsersListWrapper) getIDByEmail(email string) (uint, error) {
	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	for _, val := range active.UsersList.Users {
		if val.Email == email {
			return val.ID, nil
		}
	}

	return 0, errUserNotExists
}

func (active *UsersListWrapper) createUser(ctx context.Context, tx pgx.Tx, user *models.User) error {
	SQLCreateUser := `INSERT INTO public."user"(email, password_hash) VALUES ($1, $2);`

	var err error

	_, err = tx.Exec(ctx, SQLCreateUser, user.Email, user.PasswordHash)

	if err != nil {
		return fmt.Errorf("Something went wrong while executing create user query", err)
	}

	return nil
}

func (active *UsersListWrapper) CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	if active.UserExists(ctx, email) {
		return nil, errUserExists
	}

	user := models.User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	err := pgx.BeginFunc(ctx, active.Pool, func(tx pgx.Tx) error {
		err := active.createUser(ctx, tx, &user)
		if err != nil {

			return fmt.Errorf("Something went wrong while creating user", err)
		}
		id, err := repository.GetLastValSeq(ctx, tx, NameSeqUser)
		if err != nil {

			return fmt.Errorf("Something went wrong getting user id from seq", err)
		}
		user.ID = uint(id)

		return nil
	})

	// userInsertQuery := fmt.Sprintf(`INSERT INTO public."user"(email, password_hash) VALUES ('%s', '%s')`, email, passwordHash)
	// execquery.ExecuteInsertQuery(active.Pool, userInsertQuery)

	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	//id := active.GetLastID()

	active.UsersList.Users[user.ID] = &models.User{
		ID:           user.ID,
		PasswordHash: user.PasswordHash,
		Email:        user.Email,
	}

	if err != nil {

		return nil, err
	}
	fmt.Println("user", user)
	return &user, nil

}

func (active *UsersListWrapper) EditUser(id uint, email, passwordHash string) (*models.User, error) {
	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	usr, ok := active.UsersList.Users[id]

	if !ok {
		return nil, errUserNotExists
	}

	usr.PasswordHash = passwordHash
	usr.Email = email

	return usr, nil
}

func (active *UsersListWrapper) GetUserByEmail(email string) (*models.User, error) {
	usr, err := active.getIDByEmail(email)

	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	if err != nil {

		return nil, errUserNotExists
	}

	return active.UsersList.Users[usr], nil
}

func (active *UsersListWrapper) GetUserByID(userID uint) (*models.User, error) {
	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	usr, ok := active.UsersList.Users[userID]

	if ok {
		return usr, nil
	}

	return nil, errUserNotExists
}

func (active *UsersListWrapper) GetUserBySession(sessionID string) (*models.User, error) {
	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	id := active.UsersList.Sessions[sessionID]

	for _, val := range active.UsersList.Users {
		if val.ID == id {
			return val, nil
		}
	}

	return nil, errUserNotExists
}

func (active *UsersListWrapper) SessionExists(sessionID string) bool {
	_, exists := active.UsersList.Sessions[sessionID]

	return exists
}

func (active *UsersListWrapper) AddSession(id uint) string {
	sessionID := utils.RandString(models.SessionIDLen)

	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	user := active.UsersList.Users[id]

	active.UsersList.Sessions[sessionID] = user.ID

	return sessionID
}

func (active *UsersListWrapper) RemoveSession(sessionID string) error {
	if !active.SessionExists(sessionID) {
		return errSessionNotExists
	}

	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	delete(active.UsersList.Sessions, sessionID)

	return nil
}

func NewActiveUser(pool *pgxpool.Pool) *UsersListWrapper {
	return &UsersListWrapper{
		UsersList: &models.UsersList{
			Sessions: make(map[string]uint, 1),
			Users: map[uint]*models.User{
				1: {
					ID:           1,
					Email:        "example@mail.ru",
					PasswordHash: utils.HashPassword("123456"),
				},
			},
			UsersCount: 1,
			Mux:        sync.RWMutex{}},
		Pool: pool,
	}
}
