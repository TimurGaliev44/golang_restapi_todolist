package todo

import (
	"context"
	"errors"
	"fmt"
	appErrors "main/internal/core/errors"
	"main/internal/core/jwt"
	"main/internal/core/storage"
	"main/internal/features/dto"
	"main/internal/features/models"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

type TODO struct {
	pg     *storage.DBstruct
	rdb    *redis.Client
	logger *zap.Logger
}

func New(db *storage.DBstruct, rdb *redis.Client, logger *zap.Logger) *TODO {
	return &TODO{pg: db, rdb: rdb, logger: logger}
}

func (td *TODO) GetRedis() *redis.Client {
	return td.rdb
}

func (td *TODO) RegisterUser(ctx context.Context, user dto.UserDTO) error {
	hash_pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = td.pg.CreateUser(ctx, user.Username, hash_pass)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return appErrors.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (td *TODO) LoginUser(ctx context.Context, userDTO dto.UserDTO) (string, error) {
	user, err := td.pg.SelectUser(ctx, userDTO.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", appErrors.ErrUserNotFound
		}
		return "", err
	}
	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(userDTO.Password))
	if err != nil {
		return "", err
	}

	return jwt.GenerateToken(user.Id)
}

func (td *TODO) LogoutUser(ctx context.Context, token string, ttl time.Duration) error {
	return storage.BlockToken(ctx, td.rdb, token, ttl)
}

func (td *TODO) CreateTask(ctx context.Context, userID int, task dto.TaskDTO) (*models.TaskModel, error) {
	result, err := td.pg.CreateTask(ctx, userID, task.Title, task.Description, time.Now())
	if err != nil {
		return nil, err
	}
	if err = storage.Delete(ctx, td.rdb, fmt.Sprintf("tasks:%d", userID)); err != nil {
		td.logger.Error("failed to delete from redis", zap.Error(err))
	}
	if err = storage.Delete(ctx, td.rdb, fmt.Sprintf("tasks:uncompleted:%d", userID)); err != nil {
		td.logger.Error("failed to delete from redis", zap.Error(err))
	}
	return result, nil
}

func (td *TODO) SelectTask(ctx context.Context, userID, id int) (*models.TaskModel, error) {
	key := fmt.Sprintf("task:%d:%d", userID, id)
	result, err := storage.Get(ctx, td.rdb, key)
	if err != nil {
		td.logger.Error("redis error", zap.Error(err))
	}
	if result != nil {
		return result, nil
	}

	result, err = td.pg.SelectID(ctx, userID, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, appErrors.ErrTaskNotFound
		}
		return nil, err
	}
	if err = storage.Set(ctx, td.rdb, key, result, time.Hour); err != nil {
		td.logger.Error("failed to set value in redis", zap.Error(err))
	}
	return result, nil
}

func (td *TODO) SelectAllTask(ctx context.Context, userID int) ([]models.TaskModel, error) {
	key := fmt.Sprintf("tasks:%d", userID)
	result, err := storage.GetTasks(ctx, td.rdb, key)
	if err != nil {
		td.logger.Error("redis error", zap.Error(err))
	}
	if result != nil {
		return result, nil
	}
	tasks, err := td.pg.SelectAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err = storage.SetTasks(ctx, td.rdb, key, tasks, time.Hour); err != nil {
		td.logger.Error("failed to set value in redis", zap.Error(err))
	}
	return tasks, nil
}

func (td *TODO) SelectUncompletedTasks(ctx context.Context, userID int) ([]models.TaskModel, error) {
	key := fmt.Sprintf("tasks:uncompleted:%d", userID)
	result, err := storage.GetTasks(ctx, td.rdb, key)
	if err != nil {
		td.logger.Error("redis error", zap.Error(err))
	}
	if result != nil {
		return result, nil
	}
	tasks, err := td.pg.SelectUncompleted(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err = storage.SetTasks(ctx, td.rdb, key, tasks, time.Hour); err != nil {
		td.logger.Error("failed to set value in redis", zap.Error(err))
	}
	return tasks, nil
}

func (td *TODO) CompleteTask(ctx context.Context, userID, id int) (*models.TaskModel, error) {
	err := td.pg.CompleteTask(ctx, userID, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, appErrors.ErrTaskNotFound
		}
		return nil, err
	}
	if err = storage.Delete(ctx, td.rdb, fmt.Sprintf("task:%d:%d", userID, id)); err != nil {
		td.logger.Error("failed to delete from redis", zap.Error(err))
	}
	if err = storage.Delete(ctx, td.rdb, fmt.Sprintf("tasks:%d", userID)); err != nil {
		td.logger.Error("failed to delete from redis", zap.Error(err))
	}
	if err = storage.Delete(ctx, td.rdb, fmt.Sprintf("tasks:uncompleted:%d", userID)); err != nil {
		td.logger.Error("failed to delete from redis", zap.Error(err))
	}
	return td.pg.SelectID(ctx, userID, id)
}

func (td *TODO) DeleteTask(ctx context.Context, userID, id int) error {
	err := td.pg.DeleteTask(ctx, userID, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return appErrors.ErrTaskNotFound
		}
		return err
	}
	if err = storage.Delete(ctx, td.rdb, fmt.Sprintf("task:%d:%d", userID, id)); err != nil {
		td.logger.Error("failed to delete from redis", zap.Error(err))
	}
	if err = storage.Delete(ctx, td.rdb, fmt.Sprintf("tasks:%d", userID)); err != nil {
		td.logger.Error("failed to delete from redis", zap.Error(err))
	}
	if err = storage.Delete(ctx, td.rdb, fmt.Sprintf("tasks:uncompleted:%d", userID)); err != nil {
		td.logger.Error("failed to delete from redis", zap.Error(err))
	}
	return nil
}
