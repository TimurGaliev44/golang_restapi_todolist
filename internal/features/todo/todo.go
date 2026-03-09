package todo

import (
	"context"
	"errors"
	appErrors "main/internal/core/errors"
	"main/internal/core/jwt"
	"main/internal/features/dto"
	"main/internal/features/models"
	"main/internal/features/postgres"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TODO struct {
	pg *postgres.DBstruct
}

func New(db *postgres.DBstruct) *TODO {
	return &TODO{pg: db}
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

func (td *TODO) CreateTask(ctx context.Context, userID int, task dto.TaskDTO) (*models.TaskModel, error) {
	result, err := td.pg.CreateTask(ctx, userID, task.Title, task.Description, time.Now())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (td *TODO) SelectTask(ctx context.Context, userID, id int) (*models.TaskModel, error) {
	result, err := td.pg.SelectID(ctx, userID, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, appErrors.ErrTaskNotFound
		}
		return nil, err
	}
	return result, nil
}

func (td *TODO) SelectAllTask(ctx context.Context, userID int) ([]models.TaskModel, error) {
	tasks, err := td.pg.SelectAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (td *TODO) SelectUncompletedTasks(ctx context.Context, userID int) ([]models.TaskModel, error) {
	tasks, err := td.pg.SelectUncompleted(ctx, userID)
	if err != nil {
		return nil, err
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
	return nil
}
