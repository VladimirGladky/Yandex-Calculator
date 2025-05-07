package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/domain/models"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) User(ctx context.Context, email string) (*models.User, error) {
	var q = `SELECT email, pass_hash FROM users WHERE email = $1`
	var user models.User
	err := s.db.QueryRowContext(ctx, q, email).Scan(&user.Email, &user.PassHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Storage) SaveUser(ctx context.Context, user *models.User) error {
	var q = `INSERT INTO users (email, pass_hash) VALUES ($1, $2)`
	_, err := s.db.ExecContext(ctx, q, user.Email, user.PassHash)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddExpression(ctx context.Context, expression *models.Expression, email string) error {
	var q = `SELECT id FROM users WHERE email = $1`
	var id int
	err := s.db.QueryRowContext(ctx, q, email).Scan(&id)
	if err != nil {
		return err
	}
	q = `INSERT INTO user_expressions (id, user_id,res,status) VALUES ($1, $2, $3, $4)`
	_, err = s.db.ExecContext(ctx, q, expression.Id, id, fmt.Sprintf("%f", expression.Result), expression.Status)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetExpressions(ctx context.Context, email string) ([]models.Expression, error) {
	var q = `SELECT id FROM users WHERE email = $1`
	var id int
	err := s.db.QueryRowContext(ctx, q, email).Scan(&id)
	if err != nil {
		return nil, err
	}
	q = `SELECT id, res, status FROM user_expressions WHERE user_id = $1`
	rows, err := s.db.QueryContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	var expressions []models.Expression
	for rows.Next() {
		var expression models.Expression
		var res string
		err := rows.Scan(&expression.Id, &res, &expression.Status)
		if err != nil {
			return nil, err
		}
		expression.Result, err = strconv.ParseFloat(res, 64)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expression)
	}
	return nil, nil
}

func (s *Storage) GetExpression(ctx context.Context, id, email string) (*models.Expression, error) {
	var q = `SELECT id FROM users WHERE email = $1`
	var userId int
	err := s.db.QueryRowContext(ctx, q, email).Scan(&userId)
	if err != nil {
		return nil, err
	}
	q = `SELECT id, res, status FROM user_expressions WHERE id = $1 AND user_id = $2`
	var expression models.Expression
	var res string
	err = s.db.QueryRowContext(ctx, q, id, userId).Scan(&expression.Id, &res, &expression.Status)
	if err != nil {
		return nil, err
	}
	expression.Result, err = strconv.ParseFloat(res, 64)
	if err != nil {
		return nil, err
	}
	return &expression, nil
}
