package service

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/repository"
	"clean_arch/pkg/bitrix_client"
	"context"
	"github.com/sirupsen/logrus"
)

type CMSService struct {
	repo         repository.CMS
	log          *logrus.Logger
	bitrixClient bitrix_client.BitrixClient
}

func NewCMSService(repo repository.CMS, log *logrus.Logger, bitrixClient bitrix_client.BitrixClient) *CMSService {
	return &CMSService{
		repo:         repo,
		log:          log,
		bitrixClient: bitrixClient,
	}
}

func (s *CMSService) RequestCall(ctx context.Context, requestDTO *dto.RequestCall, userId int) error {

	if err := s.repo.RequestCall(ctx, requestDTO, userId); err != nil {
		return err
	}

	switch {
	case requestDTO.Message != "":
		if err := s.bitrixClient.Request(ctx, s.bitrixClient.CallBackLK(requestDTO)); err != nil {
			return err
		}
	case requestDTO.Name != "":
		if err := s.bitrixClient.Request(ctx, s.bitrixClient.CallBackHeader(requestDTO)); err != nil {
			return err
		}
	default:
		if err := s.bitrixClient.Request(ctx, s.bitrixClient.CallBackFooter(requestDTO)); err != nil {
			return err
		}
	}

	return nil
}

func (s *CMSService) GetCourierDeliveryInfo(ctx context.Context) *dto.CourierDeliveryInfo {
	return s.repo.GetCourierDeliveryInfo(ctx)
}

func (s *CMSService) GetCDEKDeliveryInfo(ctx context.Context) *dto.CDEKDeliveryInfo {
	return s.repo.GetCDEKDeliveryInfo(ctx)
}

func (s *CMSService) GetAllVacancies(ctx context.Context) []*dto.Vacancy {
	return s.repo.GetAllVacancies(ctx)
}

func (s *CMSService) RequestVacancy(ctx context.Context, requestInfo *dto.RequestVacancy) error {
	if err := s.repo.RequestVacancy(ctx, requestInfo); err != nil {
		return err
	}

	if err := s.bitrixClient.Request(ctx, s.bitrixClient.CallBackVacancy(requestInfo)); err != nil {
		return err
	}

	return nil
}

// TODO реализовать потоковое чтение данных
func (s *CMSService) GetRequisites(ctx context.Context) *dto.Requisites {
	return s.repo.GetRequisites(ctx)
}

func (s *CMSService) GetPrivacyPolicy(ctx context.Context) *dto.PrivacyPolicy {
	return s.repo.GetPrivacyPolicy(ctx)
}
