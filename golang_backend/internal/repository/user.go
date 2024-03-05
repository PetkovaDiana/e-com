package repository

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/models"
	"clean_arch/pkg"
	"clean_arch/pkg/cache"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

var (
	ErrDB                = fmt.Errorf("error uccured in db")
	ErrUserRegistered    = fmt.Errorf("user already registered")
	ErrUserNotRegistered = fmt.Errorf("user not registered")
	ErrWrongCode         = fmt.Errorf("wrong code")
	ErrUnregistered      = fmt.Errorf("unregistered")
	ErrUnauthorized      = fmt.Errorf("unauthorized")
	ErrWrongPassword     = fmt.Errorf("invalid password")
	ErrUserNotFound      = fmt.Errorf("user not found")

	ErrCreatingReview = fmt.Errorf("error occured creating reveiw")
)

type UserRepository struct {
	db         *gorm.DB
	cache      cache.Cache
	log        *logrus.Logger
	sessionTTL string
}

func NewUserRepository(db *gorm.DB, log *logrus.Logger, cache cache.Cache, sessionTTL string) *UserRepository {
	return &UserRepository{
		db:         db,
		cache:      cache,
		log:        log,
		sessionTTL: sessionTTL,
	}
}

func (r *UserRepository) RegisterUser(ctx context.Context, userDTO *dto.RegisterUser, session string, userId int) (int, error) {
	var newUser *models.User

	code, err := r.cache.Get(ctx, userDTO.Phone)

	if err != nil {
		r.log.Errorf("error occured register: %s", ErrWrongCode)
		return 0, ErrWrongCode
	}

	if userDTO.Code != code {
		r.log.Errorf("error occurred register: %s", ErrWrongCode)
		return 0, ErrWrongCode
	} else {
		tx := r.db.Begin()
		r.cache.Del(ctx, userDTO.Phone)

		// TODO bad code
		if len(userDTO.Inn) != 0 {
			switch r.IdentificationUserType(ctx, userDTO) {
			case nil:
				if session != "" {
					var updateUserDB *models.User
					r.db.Where("id = ?", userId).First(&updateUserDB)
					if err := tx.Model(&models.Session{}).Where("session = ?", session).Delete(&updateUserDB).Error; err != nil {
						tx.Rollback()
						return 0, err
					}
					updateUserDB.Phone = userDTO.Phone
					updateUserDB.Name = userDTO.Name
					updateUserDB.Surname = userDTO.Surname
					updateUserDB.Inn = userDTO.Inn
					updateUserDB.Kpp = userDTO.Kpp
					updateUserDB.CompanyName = userDTO.CompanyName
					updateUserDB.CompanyAddress = userDTO.CompanyAddress
					updateUserDB.ManagerName = userDTO.ManagerName

					if err := tx.Save(&updateUserDB).Error; err != nil {
						tx.Rollback()
						return 0, err
					}
					tx.Commit()
					return updateUserDB.Id, nil
				}
				newUser = &models.User{
					Phone:          userDTO.Phone,
					Name:           userDTO.Name,
					Surname:        userDTO.Surname,
					Inn:            userDTO.Inn,
					Kpp:            userDTO.Kpp,
					CompanyName:    userDTO.CompanyName,
					CompanyAddress: userDTO.CompanyAddress,
					ManagerName:    userDTO.ManagerName,
					Cart:           []models.Cart{models.Cart{}},
					Favourites:     []models.Favourite{models.Favourite{}},
					Comparison:     []models.Comparison{models.Comparison{}},
				}
				if err := tx.Create(&newUser).Error; err != nil {
					tx.Rollback()
					return 0, err
				}
				tx.Commit()
				return int(newUser.Id), nil
			default:
				return 0, fmt.Errorf("user already registered")
			}
		} else {
			if session != "" {
				var updateUserDB *models.User
				r.db.Where("id = ?", userId).First(&updateUserDB)
				if err := tx.Model(&models.Session{}).Where("session = ?", session).Delete(&updateUserDB).Error; err != nil {
					tx.Rollback()
					return 0, err
				}

				updateUserDB.Phone = userDTO.Phone
				updateUserDB.Name = userDTO.Name
				updateUserDB.Surname = userDTO.Surname
				updateUserDB.Inn = userDTO.Inn
				updateUserDB.Kpp = userDTO.Kpp
				updateUserDB.CompanyName = userDTO.CompanyName
				updateUserDB.CompanyAddress = userDTO.CompanyAddress
				updateUserDB.ManagerName = userDTO.ManagerName

				if err := tx.Save(&updateUserDB).Error; err != nil {
					tx.Rollback()
					return 0, err
				}
				tx.Commit()
				return updateUserDB.Id, nil
			}
			newUser = &models.User{
				Phone:          userDTO.Phone,
				Name:           userDTO.Name,
				Surname:        userDTO.Surname,
				Inn:            userDTO.Inn,
				Kpp:            userDTO.Kpp,
				CompanyName:    userDTO.CompanyName,
				CompanyAddress: userDTO.CompanyAddress,
				ManagerName:    userDTO.ManagerName,
				Cart:           []models.Cart{models.Cart{}},
				Favourites:     []models.Favourite{models.Favourite{}},
				Comparison:     []models.Comparison{models.Comparison{}},
			}
			if err := tx.Create(&newUser).Error; err != nil {
				tx.Rollback()
				return 0, err
			}
			tx.Commit()
			return int(newUser.Id), nil
		}
	}
}

func (r *UserRepository) AuthenticateUser(ctx context.Context, userDTO *dto.UserAuth, session string, userId int) (int, error) {
	var user *models.User

	code, err := r.cache.Get(ctx, userDTO.Phone)

	if err != nil {
		r.log.Errorf("error occured authenticate user: %s", ErrUnauthorized)
		return 0, ErrUnauthorized
	}
	if userDTO.Code != code {
		return 0, ErrWrongCode
	} else {
		r.cache.Del(ctx, userDTO.Phone)
		r.db.Where("phone = ?", userDTO.Phone).Preload(clause.Associations).First(&user)

		if session != "" {
			tx := r.db.Begin()
			var userSessionCart *models.Cart
			var userCart *models.Cart
			r.db.Where("in_order = false").Where("user_id = ?", userId).Preload("CartProducts." + clause.Associations).First(&userSessionCart)
			r.db.Where("in_order = false").Where("user_id = ?", user.Id).Preload("CartProducts." + clause.Associations).First(&userCart)
			for _, x := range userSessionCart.CartProducts {
				var newCartProducts *models.CartProduct
				if err := r.db.Preload(clause.Associations).
					Joins("inner join cartm2ms ug on ug.cart_products_id = cart_products.id").
					Joins("inner join carts g on g.id= ug.cart_id ").
					Where("g.in_order = false AND g.user_id = ?", user.Id).
					Where("product_uuid = ?", x.ProductUUID).
					First(&newCartProducts).
					Error; err != nil {
					newUserCartProduct := models.CartProduct{
						Product:    x.Product,
						Count:      x.Count,
						TotalPrice: x.TotalPrice,
						Carts:      []models.Cart{*userCart},
					}
					if err := tx.Model(&userCart.CartProducts).Create(&newUserCartProduct).Error; err != nil {
						tx.Rollback()
						return 0, err
					}
				} else {
					if err := tx.Model(&newCartProducts).Updates(&models.CartProduct{
						Count:      newCartProducts.Count + x.Count,
						TotalPrice: x.TotalPrice,
					}).Error; err != nil {
						tx.Rollback()
						return 0, err
					}
				}
			}

			var userSessionFavourites *models.Favourite
			var userFavourites *models.Favourite

			r.db.Where("user_id = ?", userId).Preload("FavouritesProducts." + clause.Associations).First(&userSessionFavourites)
			r.db.Where("user_id = ?", user.Id).Preload("FavouritesProducts." + clause.Associations).First(&userFavourites)

			for _, x := range userSessionFavourites.FavouriteProducts {
				var userFavouiriteProduct *models.FavouriteProduct
				if err := r.db.Preload(clause.Associations).
					Joins("inner join favouritesm2ms ug on ug.favourites_products_id = favourites_products.id ").
					Joins("inner join favourites g on g.id= ug.favourites_id ").
					Where("g.user_id = ?", user.Id).
					Where("product_uuid = ?", x.ProductUUID).
					First(&userFavouiriteProduct).Error; err != nil {
					newUserFavouriteProduct := models.FavouriteProduct{
						Product:    x.Product,
						Favourites: []models.Favourite{*userFavourites},
					}
					if err := tx.Model(&userFavourites.FavouriteProducts).Create(&newUserFavouriteProduct).Error; err != nil {
						tx.Rollback()
						return 0, err
					}
				}
			}

			var sessionDB *models.Session
			var userSession *models.User

			r.db.Where("id = ?", userId).First(&userSession)

			var dbError error

			dbError = tx.Select(clause.Associations).Delete(&userSessionCart.CartProducts).Error
			dbError = tx.Select(clause.Associations).Delete(&userSessionCart).Error
			dbError = tx.Select(clause.Associations).Delete(&userSessionFavourites.FavouriteProducts).Error
			dbError = tx.Select(clause.Associations).Delete(&userSessionFavourites).Error

			dbError = tx.Model(&models.Session{}).Where("session = ?", session).Delete(&userSession).Error

			dbError = tx.Where("session = ?", session).Delete(&sessionDB).Error

			dbError = tx.Unscoped().Delete(&userSession).Error

			if dbError != nil {
				tx.Rollback()
				return 0, dbError
			}
			tx.Commit()
		}
		return int(user.Id), nil
	}
}

func (r *UserRepository) RegCodeGenerator(ctx context.Context, code *dto.CodeGenerate, newCode int) error {
	var userId int

	if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select id from "user" where phone = '%s';`, code.Phone)).Scan(&userId); result.RowsAffected != 0 {
		return fmt.Errorf("user already registered")
	}

	if err := r.cache.Set(ctx, code.Phone, strconv.Itoa(newCode), 10*time.Minute); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) AuthCodeGenerator(ctx context.Context, code *dto.CodeGenerate, newCode int) error {
	var userId int

	if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select id from "user" where phone = '%s';`, code.Phone)).Scan(&userId); result.RowsAffected == 0 {
		r.log.Errorf("error occurred authenticate code sender: %s", ErrUserNotRegistered)
		return ErrUserNotRegistered
	}

	if err := r.cache.Set(ctx, code.Phone, strconv.Itoa(newCode), 10*time.Minute); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) SessionValidator() {
	r.db.Exec(fmt.Sprintf(`WITH user_id AS (
	DELETE FROM session WHERE created_at < NOW() - INTERVAL '%s minute' RETURNING user_id
    )
	DELETE FROM "user" WHERE id IN (select user_id from user_id);`, r.sessionTTL))
}

func (r *UserRepository) CreateSession(sessionKey string) (int, error) {
	newSessionDB := &models.Session{
		Session: sessionKey,
		User: models.User{
			Cart:       []models.Cart{models.Cart{}},
			Favourites: []models.Favourite{models.Favourite{}},
			Comparison: []models.Comparison{models.Comparison{}},
		},
	}
	if err := r.db.Create(&newSessionDB).Error; err != nil {
		return 0, err
	}
	return newSessionDB.UserID, nil
}

func (r *UserRepository) CheckSessionInDb(sessionKey string) (int, string, error) {
	var sessionDB *models.Session
	if err := r.db.Where("session = ?", sessionKey).First(&sessionDB).Error; err == nil {
		return sessionDB.UserID, sessionKey, nil
	} else {
		var userId int
		newSessionKey := pkg.GenerateSession()
		userId, err = r.CreateSession(newSessionKey)
		if err != nil {
			return 0, "", err
		}
		return userId, newSessionKey, nil
	}
}

func (r *UserRepository) GetUserData(ctx context.Context, userId int) (*dto.UserData, error) {
	var userDataDTO *dto.UserData
	if err := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select "user".*, e.email, e.can_to_send_personal_offers, e.can_to_send_news from "user" left join email e on "user".id = e.user_id where "user".id = %d;`, userId)).Scan(&userDataDTO).Error; err != nil {
		return nil, err
	}
	return userDataDTO, nil
}

func (r *UserRepository) UpdateEmailUser(ctx context.Context, emailInfo *dto.UpdateEmail, userId int) error {
	if err := r.db.WithContext(ctx).Exec(fmt.Sprintf(`insert into email (email, user_id) values ('%s', %d) on conflict (user_id) do update set email=excluded.email;`, emailInfo.Email, userId)).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CanToSendEmail(ctx context.Context, emailInfo *dto.CanToSendEmail, userId int) error {
	if err := r.db.WithContext(ctx).Exec(fmt.Sprintf(`update email 
	set can_to_send_news = %v, 
	can_to_send_personal_offers = %v 
	where user_id = %d`, emailInfo.CanToSendNews, emailInfo.CanToSendPersonalOffers, userId)).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateManagerName(ctx context.Context, userInfo *dto.UpdateManagerName, id int) error {
	if result := r.db.WithContext(ctx).Exec(fmt.Sprintf(`update "user" set manager_name = '%s' where id = %d`, userInfo.ManagerName, id)); result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) IdentificationUserType(ctx context.Context, userDTO *dto.RegisterUser) error {
	var userDB *models.User
	var err error

	switch {
	case len(userDTO.Inn) == 12:
		// ИП
		err = r.db.WithContext(ctx).Where("inn = ?", userDTO.Inn).Select("id").First(&userDB).Error
	default:
		// Компания
		err = r.db.WithContext(ctx).Where("kpp = ?", userDTO.Kpp).Select("id").First(&userDB).Error
	}
	// Физик никогда не попадет на этот этап, тк ему не отправится код на регистрацию
	if err != nil {
		// Такого пользователя нету
		return nil
	}
	// Такой пользователь есть
	return fmt.Errorf("user already registered")
}

func (r *UserRepository) CreateSiteReview(ctx context.Context, siteReview *dto.SiteReview) error {
	if result := r.db.WithContext(ctx).Exec(fmt.Sprintf(`insert into site_review (rating, comment) values (%d, '%s')`, siteReview.Rating, siteReview.Comment)); result.Error != nil {
		return ErrCreatingReview
	}
	return nil
}
