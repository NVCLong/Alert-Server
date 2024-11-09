package abstractrepo

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

const (
	USER_TABLE = "student_users"
)

type UserAbstractRepository interface {
	FindUserHaveDeadline(startTime time.Time, endTime time.Time) []UserAlertResponse
}

type UserRepository struct {
	db *gorm.DB
}

type UserAlertResponse struct {
	Id    uint
	Name  string
	Email string
}

func (u UserRepository) FindUserHaveDeadline(startDate time.Time, endDate time.Time) []UserAlertResponse {
	query := fmt.Sprintf(`
		SELECT distinct student_users.name AS name, student_users.id AS id, student_users.email AS email
		FROM student_users
		JOIN deadline ON deadline."userId" = student_users.id
		WHERE deadline.deadline > '%s'
		AND deadline.deadline < '%s';
	`, strings.Split(startDate.String(), " ")[0], strings.Split(endDate.String(), " ")[0])

	fmt.Println("Executing query:", query)
	var results []UserAlertResponse
	result := u.db.Raw(query).Scan(&results)
	if result.Error != nil {
		log.Println("Execute query error")
		return []UserAlertResponse{}
	}

	return results
}

func NewUserRepository(db *gorm.DB) UserAbstractRepository {
	return &UserRepository{
		db: db,
	}
}
