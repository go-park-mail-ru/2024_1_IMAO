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
	"go.uber.org/zap"
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
	Logger    *zap.SugaredLogger
}

func (active *UsersListWrapper) userExists(ctx context.Context, tx pgx.Tx, email string) (bool, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLUserExists := `SELECT EXISTS(SELECT 1 FROM public."user" WHERE email=$1 );`
	childLogger.Infof(`SELECT EXISTS(SELECT 1 FROM public."user" WHERE email=%s`, email)
	userLine := tx.QueryRow(ctx, SQLUserExists, email)

	var exists bool

	if err := userLine.Scan(&exists); err != nil {
		childLogger.Errorf("Error while scanning user exists, err=%v", err)
		return false, err
	}

	return exists, nil
}

func (active *UsersListWrapper) UserExists(ctx context.Context, email string) bool {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var exists bool

	err := pgx.BeginFunc(ctx, active.Pool, func(tx pgx.Tx) error {
		userExists, err := active.userExists(ctx, tx, email)
		exists = userExists

		return err
	})

	if err != nil {
		childLogger.Errorf("Error while executing user exists query, err=%v", err)
	}

	return exists
}

func (active *UsersListWrapper) GetLastID(ctx context.Context) uint {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var lastID uint
	_ = pgx.BeginFunc(ctx, active.Pool, func(tx pgx.Tx) error {
		id, err := repository.GetLastValSeq(ctx, tx, NameSeqUser)
		if err != nil {
			childLogger.Errorf("Something went wrong while getting user id from seq, err=%v", err)
			return fmt.Errorf("Something went wrong while getting user id from seq in func GetLastID", err)
		}
		lastID = uint(id)

		return nil
	})

	return lastID
}

// ВОЗМОЖНО НУЖНО СРОЧНО ПЕРЕПИСАТЬ, ПОТОМУ ЧТО ИЗ-ЗА ЭТОГО ЛОЖИТСЯ ПОХОД В БАЗУ
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
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLCreateUser := `INSERT INTO public."user"(email, password_hash) VALUES ($1, $2);`
	childLogger.Infof(`INSERT INTO public."user"(email, password_hash) VALUES (%s, %s)`, user.Email, user.PasswordHash)
	var err error

	_, err = tx.Exec(ctx, SQLCreateUser, user.Email, user.PasswordHash)

	if err != nil {
		childLogger.Errorf("Something went wrong while executing create user query, err=%v", err)
		return fmt.Errorf("Something went wrong while executing create user query in func createUser", err)
	}

	return nil
}

func (active *UsersListWrapper) CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

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
			childLogger.Errorf("Something went wrong while creating user, err=%v", err)
			return fmt.Errorf("Something went wrong while creating user", err)
		}
		id, err := repository.GetLastValSeq(ctx, tx, NameSeqUser)
		if err != nil {
			childLogger.Errorf("Something went wrong getting user id from seq, err=%v", err)
			return fmt.Errorf("Something went wrong getting user id from seq", err)
		}
		user.ID = uint(id)

		return nil
	})

	if err != nil {
		childLogger.Errorf("Error while creating user, err=%v", err)

		return nil, err
	}

	fmt.Println("user", user)
	return &user, nil

}

func (active *UsersListWrapper) editUserEmail(ctx context.Context, tx pgx.Tx, id uint, email string) (*models.User, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLUserExists := `UPDATE public."user"	SET email=$1 WHERE id=$2 RETURNING id, email;`
	childLogger.Infof(`UPDATE public."user"	SET email=%s WHERE id=%s RETURNING id, email;`, email, id)
	userLine := tx.QueryRow(ctx, SQLUserExists, email, id)

	user := models.User{}

	if err := userLine.Scan(&user.ID, &user.Email); err != nil {
		childLogger.Errorf("Error while scanning edit user email, err=%v", err)

		return nil, err
	}

	return &user, nil
}

func (active *UsersListWrapper) EditUserEmail(ctx context.Context, id uint, email string) (*models.User, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var user *models.User

	err := pgx.BeginFunc(ctx, active.Pool, func(tx pgx.Tx) error {
		userInner, err := active.editUserEmail(ctx, tx, id, email)
		user = userInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while editing user profile , err=%v", errUserNotExists)

		return nil, errUserNotExists
	}

	return user, nil
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

func (active *UsersListWrapper) getUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.User, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLUserByEmail := `SELECT id, email, password_hash	FROM public."user" where email = $1 `
	childLogger.Infof(`SELECT id, email, password_hash	FROM public."user" where email = %s`, email)
	userLine := tx.QueryRow(ctx, SQLUserByEmail, email)

	user := models.User{}

	if err := userLine.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		childLogger.Errorf("Something went wrong while getting user by email from seq, err=%v", err)

		return nil, err
	}

	return &user, nil
}

// НЕ ПРОТЕСТИРОВАНО
func (active *UsersListWrapper) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user *models.User

	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

	err := pgx.BeginFunc(ctx, active.Pool, func(tx pgx.Tx) error {
		userInner, err := active.getUserByEmail(ctx, tx, email)
		user = userInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while getting user by email from seq, err=%v", err)
		return nil, errUserNotExists
	}

	return user, nil
}

func (active *UsersListWrapper) getUserByID(ctx context.Context, tx pgx.Tx, id uint) (*models.User, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLUserById := `SELECT id, email, password_hash	FROM public."user" where id = $1 `
	childLogger.Infof(`SELECT id, email, password_hash	FROM public."user" where id = %s`, id)
	userLine := tx.QueryRow(ctx, SQLUserById, id)

	user := models.User{}

	if err := userLine.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		childLogger.Errorf("Something went wrong while getting user by id from seq, err=%v", err)
		return nil, err
	}

	return &user, nil
}

// НЕ ПРОТЕСТИРОВАНО
func (active *UsersListWrapper) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user *models.User

	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

	err := pgx.BeginFunc(ctx, active.Pool, func(tx pgx.Tx) error {
		userInner, err := active.getUserByID(ctx, tx, userID)
		user = userInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while getting user by id from seq, err=%v", err)
		return nil, errUserNotExists
	}

	return user, nil
}

// НЕ ПРОТЕСТИРОВАНО
func (active *UsersListWrapper) GetUserBySession(ctx context.Context, sessionID string) (*models.User, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	userID := active.UsersList.Sessions[sessionID]

	var user *models.User

	err := pgx.BeginFunc(ctx, active.Pool, func(tx pgx.Tx) error {
		userInner, err := active.getUserByID(ctx, tx, userID)
		user = userInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while getting user by session from seq, err=%v", err)
		return nil, errUserNotExists
	}

	return user, nil
}

func (active *UsersListWrapper) SessionExists(sessionID string) bool {
	_, exists := active.UsersList.Sessions[sessionID]

	return exists
}

func (active *UsersListWrapper) AddSession(id uint) string {
	sessionID := utils.RandString(models.SessionIDLen)

	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	active.UsersList.Users[id] = &models.User{
		ID: id,
	}

	user := active.UsersList.Users[id]

	active.UsersList.Sessions[sessionID] = user.ID

	return sessionID
}

func (active *UsersListWrapper) RemoveSession(ctx context.Context, sessionID string) error {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.Logger.With(
		zap.String("requestUUID", requestUUID),
	)

	if !active.SessionExists(sessionID) {
		childLogger.Errorf("Something went wrong while removing session STILL MAP, err=%v", errSessionNotExists)

		return errSessionNotExists
	}

	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	delete(active.UsersList.Sessions, sessionID)

	return nil
}

func NewActiveUser(pool *pgxpool.Pool, logger *zap.SugaredLogger) *UsersListWrapper {
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
		Pool:   pool,
		Logger: logger,
	}
}
