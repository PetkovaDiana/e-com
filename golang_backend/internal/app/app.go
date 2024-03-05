package app

import (
	"clean_arch/internal/repository"
	"clean_arch/internal/server"
	"clean_arch/internal/service"
	"clean_arch/internal/transport/http/handler"
	"clean_arch/pkg/bitrix_client"
	"clean_arch/pkg/cache"
	"clean_arch/pkg/client_1c"
	"clean_arch/pkg/email"
	"clean_arch/pkg/order_service_client"
	"clean_arch/pkg/pay_keeper"
	postgresdb "clean_arch/pkg/postgresql"
	"clean_arch/pkg/rec_service_client"
	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	var log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	if err := godotenv.Load(".env"); err != nil {
		log.Errorf("error loading env variables: %s", err.Error())
	}

	if err := initConfig(); err != nil {
		log.Errorf("error initializing configs: %s", err.Error())
	}

	location, err := time.LoadLocation(viper.GetString("time.location"))

	log.Info(location)

	db, err := postgresdb.NewPostgresDB(postgresdb.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
	}, location)

	cookie := http.Cookie{
		Name:     viper.GetString("cookie.name"),
		Path:     viper.GetString("cookie.path"),
		Domain:   viper.GetString("cookie.domain"),
		MaxAge:   viper.GetInt("cookie.max_age"),
		Secure:   viper.GetBool("cookie.secure"),
		HttpOnly: viper.GetBool("cookie.httponly"),
		SameSite: http.SameSite(viper.GetInt("cookie.same_site")),
	}

	if err != nil {
		log.Errorf("error initializing db: %s", err.Error())
	}

	cacheClient := cache.NewRedisClient(&cache.Config{
		Address: viper.GetString("redis.host"),
		DB:      0,
	})

	//TODO move to .env file
	bitrixClient := bitrix_client.NewClient(
		&http.Client{},
		log,
		"zoggmbrz1acc67bi",
		"x8yr5jckqv6n29dp",
		"rqbk7i4wfx0mgayc",
		"13bio84em8batd1n")

	recClient := rec_service_client.NewRecommendationServiceClient(
		&http.Client{},
		log)

	client1c := client_1c.NewClient1C(
		&http.Client{},
		log,
		viper.GetString("cl_1c.username"),
		os.Getenv("1C_PASSWORD"))

	emailClient := email.NewEmail(&email.Config{
		SenderEmail:       viper.GetString("smtp.sender"),
		SenderAppPassword: os.Getenv("SMTP_PASSWORD"),
		SmtpHost:          viper.GetString("smtp.host"),
		SmtpPort:          viper.GetString("smtp.port"),
		DirectorEmail:     viper.GetString("smtp.director_email"),
		CourierEmail:      viper.GetString("smtp.courier_email"),
		MimeHeaders:       viper.GetString("smtp.mime_headers"),
	})

	payKeeperClient := pay_keeper.NewPayKeeperClient(
		pay_keeper.BasicAuth{
			Username: "admin",
			Password: "5858f6d5f2d3",
		}, &http.Client{},
		log,
		os.Getenv("PAY_KEEPER_SK"))

	orderServiceClient := order_service_client.NewOrderServiceClient(&http.Client{}, log)

	repos := repository.NewRepository(db,
		cacheClient,
		log,
		location,
		viper.GetString("time.format"),
		viper.GetString("session_ttl"),
		viper.GetString("media_root"))

	service := service.NewService(
		repos,
		log,
		viper.GetString("token_ttl"),
		os.Getenv("API_KEY"),
		emailClient,
		bitrixClient,
		recClient,
		client1c,
		payKeeperClient,
		orderServiceClient)

	httpHandler := handler.NewHttpHandler(service, &cookie)

	srv := server.Server{}
	go func() {
		if err = srv.Run(viper.GetString("port"), httpHandler.InitRoutes()); err != nil {
			log.Errorf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Info("BESM app started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("BESM app shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Errorf("error occured on server shutting down: %s", err.Error())

	}
	sqlDB, err := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
