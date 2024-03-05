package repository

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/models"
	"context"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
)

type CMSRepository struct {
	db         *gorm.DB
	locTime    *time.Location
	timeFormat string
}

func NewCMSRepository(db *gorm.DB, locTime *time.Location, timeFormat string) *CMSRepository {
	return &CMSRepository{
		db:         db,
		locTime:    locTime,
		timeFormat: timeFormat,
	}
}

func (r *CMSRepository) RequestCall(ctx context.Context, requestDTO *dto.RequestCall, userId int) error {

	wg := &sync.WaitGroup{}

	var err1 error
	wg.Add(1)
	go func() {
		defer wg.Done()
		requestDB := &models.RequestCall{
			Name:      requestDTO.Name,
			Phone:     requestDTO.Phone,
			UserID:    userId,
			Message:   requestDTO.Message,
			CreatedAt: time.Now().In(r.locTime),
		}

		err1 = r.db.WithContext(ctx).Create(&requestDB).Error
	}()

	if requestDTO.Message != "" {

		if requestDTO.Email != "" {
			if err := r.db.WithContext(ctx).Exec(fmt.Sprintf(`insert into email (email, user_id) values ('%s', %d) on conflict (user_id) do update set email=excluded.email;`, requestDTO.Email, userId)).Error; err != nil {
				return err
			}

		}
		r.db.WithContext(ctx).Raw(fmt.Sprintf(`
			select email, u.name, u.surname, u.company_name, u.inn, u.phone from "user" u 
			left join email on u.id = email.user_id where u.id = %d`, userId)).Scan(&requestDTO)
	}

	wg.Wait()

	if err1 != nil {
		return err1
	}

	return nil
}

func (r *CMSRepository) GetCourierDeliveryInfo(ctx context.Context) *dto.CourierDeliveryInfo {
	rows, err := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select courier_delivery_info.description, cdti.mon, cdti.tue, cdti.wen, 
	cdti.thu, cdti.fri, cdti.sat, cdti.sun from courier_delivery_info
    inner join courier_delivery_time_info cdti on cdti.id = courier_delivery_info.courier_delivery_time_info_id limit 1`)).Rows()
	if err != nil {
		return nil
	}
	defer func() error {
		if err := rows.Close(); err != nil {
			return err
		}
		return nil
	}()
	var courierInfoDTO *dto.CourierDeliveryInfo
	var courierTimeInfoDTO *dto.CourierDeliveryTimeInfo
	for rows.Next() {
		r.db.WithContext(ctx).ScanRows(rows, &courierInfoDTO)
		r.db.WithContext(ctx).ScanRows(rows, &courierTimeInfoDTO)
		courierInfoDTO.CourierDeliveryTimeInfo = *courierTimeInfoDTO
	}

	return courierInfoDTO
}

func (r *CMSRepository) GetCDEKDeliveryInfo(ctx context.Context) *dto.CDEKDeliveryInfo {
	var cdekInfoDTO *dto.CDEKDeliveryInfo
	if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select description from cdek_delivery_info limit 1`)).Scan(&cdekInfoDTO); result.RowsAffected == 0 {
		return nil
	}
	return cdekInfoDTO
}

func (r *CMSRepository) GetAllVacancies(ctx context.Context) []*dto.Vacancy {
	var allVacancies []*dto.Vacancy
	if err := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select * from vacancy`)).Scan(&allVacancies); err.Error != nil {
		return nil
	}
	return allVacancies
}

func (r *CMSRepository) RequestVacancy(ctx context.Context, requestInfo *dto.RequestVacancy) error {
	if err := r.db.WithContext(ctx).Raw(fmt.Sprintf(`insert into request_vacancy 
		(phone, name, lastname, surname, email, vacancy_id, created_at, comment) values 
		('%s', '%s', '%s', '%s', '%s', %s, ?, '%s') returning (select title from vacancy where id = %s)`, requestInfo.Phone, requestInfo.Name, requestInfo.Lastname, requestInfo.Surname,
		requestInfo.Email, requestInfo.VacancyID, requestInfo.Comment, requestInfo.VacancyID), time.Now().Format(r.timeFormat)).Scan(&requestInfo.VacancyTitle); err != nil {
		return err.Error
	}
	return nil
}

func (r *CMSRepository) GetRequisites(ctx context.Context) *dto.Requisites {
	var requisitesDTO *dto.Requisites
	r.db.WithContext(ctx).Raw(fmt.Sprintf(`select text from requisites;`)).Scan(&requisitesDTO)
	return requisitesDTO
}

func (r *CMSRepository) GetPrivacyPolicy(ctx context.Context) *dto.PrivacyPolicy {
	var privacyPolicyDTO *dto.PrivacyPolicy
	r.db.WithContext(ctx).Raw(fmt.Sprintf(`select text from privacy_policy;`)).Scan(&privacyPolicyDTO)
	return privacyPolicyDTO
}
