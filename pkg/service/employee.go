package service

import (
	"crypto/sha1"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"
	"warehouse_backend/pkg/model"
	"warehouse_backend/pkg/repository"
)

type EmployeeService struct {
	repositoryEmployee  repository.Employee
	repositoryEquipment repository.Equipment
}

const (
	salt       = "12345678"
	signingKey = "12345678"
	tokenTTL   = 10 * time.Minute
	dictionary = "abcdefghijklmnopqrstuvwxyz0123456789"
	length     = 10
)

func NewEmployeeService(
	repositoryEmployee repository.Employee,
	repositoryEquipment repository.Equipment) *EmployeeService {
	return &EmployeeService{
		repositoryEmployee:  repositoryEmployee,
		repositoryEquipment: repositoryEquipment,
	}
}

func (s *EmployeeService) Create(name, phone, email string) error {
	if _, err := s.repositoryEmployee.FindByPhone(phone); err == nil {
		return errors.New("phone already exists")
	}
	return s.repositoryEmployee.Create(name, phone, email)
}

func (s *EmployeeService) GetById(id int) (model.Employee, error) {
	return s.repositoryEmployee.GetById(id)
}

func (s *EmployeeService) GetByDepartment(ids []int, id int) ([]model.Employee, error) {
	return s.repositoryEmployee.GetByDepartment(ids, id)
}

func (s *EmployeeService) GetAll() ([]model.Employee, error) {
	return s.repositoryEmployee.GetAll()
}

func (s *EmployeeService) GetFree() ([]model.Employee, error) {
	return s.repositoryEmployee.GetFree()
}

func (s *EmployeeService) GetAllButOne(id int) ([]model.Employee, error) {
	return s.repositoryEmployee.GetAllButOne(id)
}

func (s *EmployeeService) FindUser(login, password string) (int, error) {
	return s.repositoryEmployee.FindUser(login, generatePasswordHash(password))
}

func (s *EmployeeService) FindByHash(hash string) (int, error) {
	return s.repositoryEmployee.FindByHash(hash)
}

func (s *EmployeeService) AddToDepartment(id, department int) error {
	return s.repositoryEmployee.AddToDepartment(id, department)
}

func (s *EmployeeService) RemoveFromDepartment(idDepartment, idEmployee int) error {
	equipments, err := s.repositoryEquipment.GetByLocationDepartmentEmployee(idDepartment, idEmployee)
	if err != nil {
		return err
	}
	if len(equipments) > 0 {
		return errors.New("have equipment")
	}
	return s.repositoryEmployee.RemoveFromDepartment(idEmployee)
}

func (s *EmployeeService) Update(id int, name, phone, email string) error {
	findId, err := s.repositoryEmployee.FindByPhone(phone)
	if findId != id && err == nil {
		return errors.New("phone already exists")
	}
	return s.repositoryEmployee.Update(id, name, phone, email)
}

func (s *EmployeeService) Delete(id int) error {
	equipments, err := s.repositoryEquipment.GetByLocationEmployee(id)
	if err != nil {
		return err
	}
	if len(equipments) > 0 {
		return errors.New("have equipment")
	}
	return s.repositoryEmployee.Delete(id)
}

func (s *EmployeeService) GenerateToken(id int) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(id),
		ExpiresAt: jwt.NewNumericDate(time.Unix(time.Now().Add(tokenTTL).Unix(), 0)),
		IssuedAt:  jwt.NewNumericDate(time.Unix(time.Now().Unix(), 0)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func (s *EmployeeService) ParseToken(accessToken string) (interface{}, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing")
		}
		return []byte(signingKey), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"], nil
	} else {
		return "", err
	}
}

func (s *EmployeeService) GenerateHash(id int) (string, error) {
	str := generateString()
	err := s.repositoryEmployee.SetHash(id, str)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (s *EmployeeService) Activate(id int) error {
	employee, _ := s.repositoryEmployee.GetById(id)
	password := generateString()
	if err := s.repositoryEmployee.Activate(id, generatePasswordHash(password)); err != nil {
		return err
	}
	if err := sendMail(employee.Email, employee.Phone, password); err != nil {
		return err
	}
	return nil
}

func (s *EmployeeService) Deactivate(id int) error {
	return s.repositoryEmployee.Deactivate(id)
}

func (s *EmployeeService) ResetPassword(id int) error {
	employee, _ := s.repositoryEmployee.GetById(id)
	password := generateString()
	if err := s.repositoryEmployee.ResetPassword(id, generatePasswordHash(password)); err != nil {
		return err
	}
	if err := sendMail(employee.Email, employee.Phone, password); err != nil {
		return err
	}
	return nil
}

func generateString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune(dictionary)
	var str strings.Builder
	for i := 0; i < length; i++ {
		str.WriteRune(chars[rand.Intn(len(chars))])
	}
	return str.String()
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func sendMail(recipient, phone, password string) error {
	authEmail := ""
	authPassword := ""
	host := "smtp.yandex.ru"
	port := "465"
	subject := "Authorization data"
	body := fmt.Sprintf("Login: %s\nPassword: %s", phone, password)
	sendString := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", authEmail, recipient, subject, body)
	msg := []byte(sendString)
	auth := smtp.PlainAuth("", authEmail, authPassword, host)
	conf := &tls.Config{ServerName: host}
	conn, err := tls.Dial("tcp", host+":"+port, conf)
	if err != nil {
		return err
	}
	cl, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	if err := cl.Auth(auth); err != nil {
		return err
	}
	if err := cl.Mail(authEmail); err != nil {
		return err
	}
	if err := cl.Rcpt(recipient); err != nil {
		return err
	}
	w, err := cl.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	if err := cl.Quit(); err != nil {
		return err
	}
	return nil
}
