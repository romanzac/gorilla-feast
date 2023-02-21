package dbhandler

import (
	"errors"
	"github.com/romanzac/gorilla-feast/domain/model"
	"github.com/romanzac/gorilla-feast/infra/database"
	"github.com/romanzac/gorilla-feast/infra/ssha"
	"github.com/romanzac/gorilla-feast/middleware"
	"gorm.io/gorm"
	"time"
)

// DbUserRepo represents access to user data
type DbUserRepo struct {
	DB *gorm.DB
}

// NewDbUserRepo creates new database repository for Users
func NewDbUserRepo() *DbUserRepo {
	dbUserRepo := new(DbUserRepo)
	dbUserRepo.DB = database.DB

	return dbUserRepo
}

func (r *DbUserRepo) Find(acct, fullname, sortQuery string, limit, offset int, noDetail bool) ([]model.User, error) {
	var users []model.User

	if acct == "" && fullname == "" && sortQuery != "" && limit != 0 && offset != 0 {
		if err := r.DB.Select("acct", "fullname").
			Limit(limit).Offset(offset).Order(sortQuery).Find(&users).Error; err != nil {
			return []model.User{}, err
		}
		return users, nil
	}

	if acct != "" && noDetail == false {
		if err := database.DB.Select("acct", "fullname", "created_at", "updated_at").
			Where("acct = ?", acct).Find(&users).Error; err != nil {
			return []model.User{}, err
		}
		return users, nil
	}

	if acct != "" && noDetail == true {
		if err := database.DB.Select("acct", "fullname").
			Where("acct = ?", acct).Find(&users).Error; err != nil {
			return []model.User{}, err
		}
		return users, nil
	}

	if acct == "" && fullname != "" {
		if err := database.DB.Select("acct", "fullname").Where("fullname LIKE ?", fullname).
			Find(&users).Error; err != nil {
			return []model.User{}, err
		}
		return users, nil
	}

	return []model.User{}, errors.New("invalid query")
}

func (r *DbUserRepo) Create(acct, fullname, pwd string) error {
	var u model.User
	u.Acct = acct
	u.Fullname = fullname

	u.Pwd, _ = ssha.GeneratePassword(pwd, 32)
	if err := r.DB.Create(&u).Error; err != nil {
		return err
	}

	return nil
}

func (r *DbUserRepo) Update(acct, fullname, pwd string) error {
	var u model.User

	pwd, _ = ssha.GeneratePassword(pwd, 32)

	t := time.Now()
	result := database.DB.Model(&u).Where("acct = ?", acct).
		Updates(model.User{Pwd: pwd, Fullname: fullname, UpdatedAt: &t})

	if result.RowsAffected == 0 {
		return errors.New("no rows were affected")
	}
	if result.Error != nil {
		return result.Error

	}

	return nil
}

func (r *DbUserRepo) Delete(acct string) error {
	var u model.User

	result := database.DB.Where("acct = ?", acct).Delete(&u)

	if result.RowsAffected == 0 {
		return errors.New("nothing was deleted")
	}
	if result.Error != nil {
		return result.Error

	}

	return nil
}

// Validate user for login purposes, return JWT token if passed
func (r *DbUserRepo) Validate(acct, pwd string) (middleware.JWTToken, error) {
	var u model.User

	if err := r.DB.Select("acct", "pwd", "fullname").
		Where("acct = ?", acct).First(&u).Error; err != nil {
		return middleware.JWTToken{}, errors.New("DB query error to find user \"" + acct + "\"")
	}

	pwdOK, _ := ssha.ValidatePassword(pwd, u.Pwd)
	if !pwdOK {
		return middleware.JWTToken{}, errors.New("Password incorrect for user \"" + acct + "\"")
	}

	token, err := middleware.GenerateJWT(u.Acct, u.Fullname)
	if err != nil {
		return middleware.JWTToken{}, errors.New("Error: " + err.Error())
	}

	return token, nil
}
