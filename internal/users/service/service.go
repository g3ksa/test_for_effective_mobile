package service

import (
	"UserService/internal/models"
	"UserService/internal/users"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
)

type AgeResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type GenderResponse struct {
	Count       int    `json:"count"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Probability int    `json:"probability"`
}

type CountryResponse struct {
	Country_id  string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type NationalityResponse struct {
	Count   int               `json:"count"`
	Name    string            `json:"name"`
	Country []CountryResponse `json:"country"`
}

type UserService struct {
	repo users.Repository
	log  *slog.Logger
}

func New(repo users.Repository, log *slog.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
	}
}

func (u *UserService) GetWithFilter(filters map[string]string, page, pageSize int) (*[]models.User, error) {
	users, err := u.repo.GetWithFilter(filters, page, pageSize)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserService) Delete(userId int) error {
	err := u.repo.Delete(userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) Edit(req users.EditUserRequest) (*models.User, error) {
	if req.Field == "" {
		return nil, fmt.Errorf("поле не может быть пустым")
	}
	if req.Id == 0 {
		return nil, fmt.Errorf("id не может быть пустым")
	}
	newUser, err := u.repo.Edit(req)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (u *UserService) Create(user models.User) (*models.User, error) {
	if user.Name == "" {
		return nil, fmt.Errorf("имя пользователя не может быть пустым")
	}
	if user.Surname == "" {
		return nil, fmt.Errorf("фамилия пользователя не может быть пустой")
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		u.setAgeByName(&user)
		wg.Done()
	}()
	go func() {
		u.setGenderByName(&user)
		wg.Done()
	}()
	go func() {
		u.setNationalityByName(&user)
		wg.Done()
	}()
	wg.Wait()
	// u.getAgeByName(&user)
	// u.getGenderByName(&user)
	// u.getNationalityByName(&user)

	if user.Age == 0 {
		return nil, fmt.Errorf("ошибка при получении возвраста пользователя")
	}
	if user.Gender == "" {
		return nil, fmt.Errorf("ошибка при получении гендера пользователя")
	}
	if user.Nationality == "" {
		return nil, fmt.Errorf("ошибка при получении национальности пользователя")
	}

	newUser, err := u.repo.Create(user)
	if err != nil {
		return nil, err
	}
	return newUser, nil
	// return &user, nil
}

func (u *UserService) setAgeByName(user *models.User) {
	op := "UserService.getAgeByName"
	result := new(AgeResponse)
	response, err := http.Get(
		fmt.Sprintf("https://api.agify.io/?name=%s", user.Name),
	)
	u.log.With(slog.String("operation", op)).Info("response", slog.Any("response", response))
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		u.log.Info("ошибка при получении возраста")
	}
	json.Unmarshal(body, &result)
	user.Age = result.Age
}

func (u *UserService) setGenderByName(user *models.User) {
	op := "UserService.getGenderByName"
	result := new(GenderResponse)
	response, err := http.Get(
		fmt.Sprintf("https://api.genderize.io/?name=%s", user.Name),
	)
	u.log.With(slog.String("operation", op)).Info("response", slog.Any("response", response))
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		u.log.Info("ошибка при получении пола")
	}
	json.Unmarshal(body, &result)
	user.Gender = result.Gender
}

func (u *UserService) setNationalityByName(user *models.User) {
	op := "UserService.getNationalityByName"
	result := new(NationalityResponse)
	response, err := http.Get(
		fmt.Sprintf("https://api.nationalize.io/?name=%s", user.Name),
	)
	u.log.With(slog.String("operation", op)).Info("response", slog.Any("response", response))
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		u.log.Info("ошибка при получении национальности")
	}
	json.Unmarshal(body, &result)
	if len(result.Country) > 0 {
		user.Nationality = result.Country[0].Country_id
	}
}
