package postgresdb

import (
	"clean_arch/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func NewPostgresDB(cfg Config, locTime *time.Location) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}, NowFunc: func() time.Time {
			return time.Now().In(locTime)
		},
		DisableAutomaticPing: true,
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()

	// TODO получать данные из конфига
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)

	err = migrate(db)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *gorm.DB) error {

	return db.AutoMigrate(
		&models.Session{},
		&models.Review{},
		&models.ReviewPhotos{},
		&models.Product{},
		&models.Characteristic{},
		&models.ProductCharacteristic{},
		&models.ProductFiles{},
		&models.ProductStatistic{},
		&models.Category{},
		&models.User{},
		&models.Email{},
		&models.Code{},
		&models.CartProduct{},
		&models.FavouriteProduct{},
		&models.ComparisonProduct{},
		&models.Cart{},
		&models.Comparison{},
		&models.PaymentMethod{},
		&models.OrderStatus{},
		&models.PickUpPointTime{},
		&models.PickUpPointStockTitle{},
		&models.PickUpPointStockDescription{},
		&models.PickUpPoint{},
		&models.Order{},
		&models.PayKeeperInfo{},
		&models.DeliveryType{},
		&models.CourierDelivery{},
		&models.SelfDelivery{},
		&models.CDEKDelivery{},
		&models.RequestCall{},
		&models.PromoCode{},
		&models.CourierDeliveryInfo{},
		&models.CourierDeliveryTimeInfo{},
		&models.CDEKDeliveryInfo{},
		&models.DeliveryTypeInfo{},
		&models.EmailStatic{},
		&models.Vacancy{},
		&models.RequestVacancy{},
		&models.Requisites{},
		&models.PrivacyPolicy{},
	)
}
