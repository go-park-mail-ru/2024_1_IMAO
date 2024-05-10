package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	utils "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	errUserNotExists    = errors.New("user does not exist")
	errUserExists       = errors.New("user already exists")
	errSessionNotExists = errors.New("session does not exist")
	errWrongData        = errors.New("wrong registration data")
	NameSeqUser         = pgx.Identifier{"public", "user_id_seq"} //nolint:gochecknoglobals
)

type UserStorage struct {
	pool        *pgxpool.Pool
	sessionList *models.SessionList
	metrics     *mymetrics.DatabaseMetrics
}

func NewUserStorage(pool *pgxpool.Pool, metrics *mymetrics.DatabaseMetrics) *UserStorage {
	return &UserStorage{
		pool:        pool,
		sessionList: NewSessionList(),
		metrics:     metrics,
	}
}

func (active *UserStorage) userExists(ctx context.Context, tx pgx.Tx, email string) (bool, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUserExists := `SELECT EXISTS(SELECT 1 FROM public."user" WHERE email=$1 );`

	logging.LogInfo(logger, "SELECT FROM user")

	start := time.Now()
	userLine := tx.QueryRow(ctx, SQLUserExists, email)
	active.metrics.AddDuration(funcName, time.Since(start))

	var exists bool

	if err := userLine.Scan(&exists); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning user exists, err=%v", err))

		return false, err
	}

	return exists, nil
}

func (active *UserStorage) UserExists(ctx context.Context, email string) bool {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var exists bool

	err := pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
		userExists, err := active.userExists(ctx, tx, email)
		exists = userExists

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing user exists query, err=%v", err))

	}

	return exists
}

func (active *UserStorage) GetLastID(ctx context.Context) uint {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var lastID uint
	_ = pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
		id, err := repository.GetLastValSeq(ctx, tx, NameSeqUser)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting user id from seq, err=%v", err))

			return err
		}
		lastID = uint(id)

		return nil
	})

	return lastID
}

func (active *UserStorage) createUser(ctx context.Context, tx pgx.Tx, user *models.User) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateUser := `INSERT INTO public."user"(email, password_hash) VALUES ($1, $2);`

	logging.LogInfo(logger, "INSERT INTO user")

	var err error

	start := time.Now()
	_, err = tx.Exec(ctx, SQLCreateUser, user.Email, user.PasswordHash)
	active.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing create user query, err=%v", err))
		active.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func (active *UserStorage) CreateUser(ctx context.Context, email,
	password, passwordRepeat string) (*models.User, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if active.UserExists(ctx, email) {
		return nil, errUserExists
	}

	errs := utils.Validate(email, password)
	if password != passwordRepeat {
		errs = append(errs, "Passwords do not match")
	}
	if errs != nil {
		return nil, errWrongData
	}

	user := models.User{
		Email:        email,
		PasswordHash: utils.HashPassword(password),
	}

	err := pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
		err := active.createUser(ctx, tx, &user)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while creating user, err=%v", err))

			return err
		}
		id, err := repository.GetLastValSeq(ctx, tx, NameSeqUser)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong getting user id from seq, err=%v", err))

			return err
		}
		user.ID = uint(id)

		return nil
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while creating user, err=%v", err))

		return nil, err
	}

	fmt.Println("user", user)
	return &user, nil

}

func (active *UserStorage) editUserEmail(ctx context.Context, tx pgx.Tx, id uint, email string) (*models.User, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUserExists := `UPDATE public."user"	SET email=$1 WHERE id=$2 RETURNING id, email;`

	logging.LogInfo(logger, "UPDATE user")

	start := time.Now()
	userLine := tx.QueryRow(ctx, SQLUserExists, email, id)
	active.metrics.AddDuration(funcName, time.Since(start))

	user := models.User{}

	if err := userLine.Scan(&user.ID, &user.Email); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning edit user email, err=%v", err))

		return nil, err
	}

	return &user, nil
}

func (active *UserStorage) EditUserEmail(ctx context.Context, id uint, email string) (*models.User, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))
	err := utils.ValidateEmail(email)
	if err != nil {
		return nil, err
	}

	var user *models.User

	err = pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
		userInner, err := active.editUserEmail(ctx, tx, id, email)
		user = userInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while editing user profile , err=%v", errUserNotExists))

		return nil, errUserNotExists
	}

	return user, nil
}

func (active *UserStorage) getUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.User, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUserByEmail := `SELECT id, email, password_hash	FROM public."user" where email = $1 `

	logging.LogInfo(logger, "SELECT FROM user")

	start := time.Now()
	userLine := tx.QueryRow(ctx, SQLUserByEmail, email)
	active.metrics.AddDuration(funcName, time.Since(start))

	user := models.User{}

	if err := userLine.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting user by email from seq, err=%v", err))

		return nil, err
	}

	return &user, nil
}

// НЕ ПРОТЕСТИРОВАНО
func (active *UserStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var user *models.User

	err := pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
		userInner, err := active.getUserByEmail(ctx, tx, email)
		user = userInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting user by email from seq, err=%v", err))

		return nil, errUserNotExists
	}

	return user, nil
}

func (active *UserStorage) getUserByID(ctx context.Context, tx pgx.Tx, id uint) (*models.User, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUserById := `SELECT id, email, password_hash	FROM public."user" where id = $1 `

	logging.LogInfo(logger, "SELECT FROM user")

	start := time.Now()
	userLine := tx.QueryRow(ctx, SQLUserById, id)
	active.metrics.AddDuration(funcName, time.Since(start))

	user := models.User{}

	if err := userLine.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting user by id from seq, err=%v", err))

		return nil, err
	}

	return &user, nil
}

// НЕ ПРОТЕСТИРОВАНО
func (active *UserStorage) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var user *models.User

	err := pgx.BeginFunc(ctx, active.pool, func(tx pgx.Tx) error {
		userInner, err := active.getUserByID(ctx, tx, userID)
		user = userInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting user by id from seq, err=%v", err))

		return nil, errUserNotExists
	}

	return user, nil
}

// НЕ ПРОТЕСТИРОВАНО
func (active *UserStorage) GetUserBySession(ctx context.Context, sessionID string) (*models.User, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

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
		logging.LogError(logger,
			fmt.Errorf("something went wrong while getting user by session from seq, err=%v", err))

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
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if !active.SessionExists(sessionID) {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while removing session from MAP, err=%v", errSessionNotExists))

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
