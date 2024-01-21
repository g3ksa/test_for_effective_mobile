package psql

import (
	"UserService/internal/models"
	"UserService/internal/users"
	"database/sql"
	"fmt"
	"log/slog"
)

type UserRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func New(db *sql.DB, log *slog.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

func (r *UserRepository) Create(user models.User) (*models.User, error) {
	op := "psql.Create"
	r.log.With(slog.String("operation", op)).Info("inserting user", slog.Any("user", user))
	query := r.db.QueryRow(
		"INSERT INTO users(name, surname, patronymic, age, gender, nationality) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, name, surname, patronymic, age, gender, nationality",
		user.Name,
		user.Surname,
		user.Patronymic,
		user.Age,
		user.Gender,
		user.Nationality,
	)
	err := query.Scan(
		&user.Id,
		&user.Name,
		&user.Surname,
		&user.Patronymic,
		&user.Age,
		&user.Gender,
		&user.Nationality,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании пользователя")
	}

	return &user, nil
}

func (r *UserRepository) Edit(req users.EditUserRequest) (*models.User, error) {
	op := "psql.Edit"

	user := models.User{}

	r.log.
		With(slog.String("operation", op)).
		Info("editing user", slog.Any("request", req))
	_, err := r.db.Exec(
		fmt.Sprintf("UPDATE users SET %s = $1 WHERE id = $2", req.Field),
		req.NewValue,
		req.Id,
	)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow("SELECT * FROM users WHERE id = $1", req.Id).Scan(
		&user.Id,
		&user.Name,
		&user.Surname,
		&user.Patronymic,
		&user.Age,
		&user.Gender,
		&user.Nationality,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Delete(userId int) error {
	op := "psql.Delete"
	r.log.With(slog.String("operation", op)).Info("deleting user", slog.Any("user_id", userId))
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetWithFilter(filters map[string]string, page, pageSize int) (*[]models.User, error) {
	op := "psql.GetWithFilter"
	users := make([]models.User, 0)
	r.log.With(slog.String("operation", op)).Info("getting users with filter", slog.Any("filters", filters), slog.Any("page", page), slog.Any("pageSize", pageSize))
	queryString := "SELECT * FROM users WHERE 1=1"
	for k, v := range filters {
		if v == "" {
			continue
		}
		queryString += fmt.Sprintf(" AND %s = '%s'", k, v)
	}
	queryString += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(
		queryString,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Surname,
			&user.Patronymic,
			&user.Age,
			&user.Gender,
			&user.Nationality,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &users, nil
}
