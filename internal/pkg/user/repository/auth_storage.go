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

type UserStorage struct {
	pool        *pgxpool.Pool
	logger      *zap.SugaredLogger
	sessionList *models.SessionList
}

func NewUserStorage(pool *pgxpool.Pool, logger *zap.SugaredLogger) *UserStorage {
	return &UserStorage{
		pool:        pool,
		logger:      logger,
		sessionList: NewSessionList(),
	}
}

func (active *UserStorage) userExists(ctx context.Context, tx pgx.Tx, email string) (bool, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
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

func (active *UserStorage) UserExists(ctx context.Context, email string) bool {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var exists bool

	err := pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
		userExists, err := active.userExists(ctx, tx, email)
		exists = userExists

		return err
	})

	if err != nil {
		childLogger.Errorf("Error while executing user exists query, err=%v", err)
	}

	return exists
}

func (active *UserStorage) GetLastID(ctx context.Context) uint {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var lastID uint
	_ = pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
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

func (active *UserStorage) createUser(ctx context.Context, tx pgx.Tx, user *models.User) error {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
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

func (active *UserStorage) CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	if active.UserExists(ctx, email) {
		return nil, errUserExists
	}

	user := models.User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	err := pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
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

func (active *UserStorage) editUserEmail(ctx context.Context, tx pgx.Tx, id uint, email string) (*models.User, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
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

func (active *UserStorage) EditUserEmail(ctx context.Context, id uint, email string) (*models.User, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var user *models.User

	err := pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
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

func (active *UserStorage) getUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.User, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
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
func (active *UserStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user *models.User

	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	err := pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
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

func (active *UserStorage) getUserByID(ctx context.Context, tx pgx.Tx, id uint) (*models.User, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
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
func (active *UserStorage) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user *models.User

	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	err := pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
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
func (active *UserStorage) GetUserBySession(ctx context.Context, sessionID string) (*models.User, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	active.sessionList.Mux.Lock()
	defer active.sessionList.Mux.Unlock()

	userID := active.sessionList.Sessions[sessionID]

	var user *models.User

	err := pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
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

func (active *UserStorage) SessionExists(sessionID string) bool {
	_, exists := active.sessionList.Sessions[sessionID]

	return exists
}

func (active *UserStorage) AddSession(id uint) string {
	sessionID := utils.RandString(models.SessionIDLen)

	active.sessionList.Mux.Lock()
	defer active.sessionList.Mux.Unlock()

	active.sessionList.Sessions[sessionID] = id

	return sessionID
}

func (active *UserStorage) RemoveSession(ctx context.Context, sessionID string) error {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := active.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	if !active.SessionExists(sessionID) {
		childLogger.Errorf("Something went wrong while removing session STILL MAP, err=%v", errSessionNotExists)

		return errSessionNotExists
	}

	active.sessionList.Mux.Lock()
	defer active.sessionList.Mux.Unlock()

	delete(active.sessionList.Sessions, sessionID)

	return nil
}

func NewSessionList() *models.SessionList {
	return &models.SessionList{
		Sessions: make(map[string]uint, 1),
		Mux:      sync.RWMutex{},
	}
}
