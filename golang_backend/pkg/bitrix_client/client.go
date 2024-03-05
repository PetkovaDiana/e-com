package bitrix_client

import (
	"clean_arch/internal/dto"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

type Method string

const (
	CallBackHeader  Method = "call_back_header"
	CallBackFooter  Method = "call_back_footer"
	CallBackLK      Method = "call_back_lk"
	CallBackVacancy Method = "call_back_vacancy"
)

type bitrixClient struct {
	client  *http.Client
	log     *logrus.Logger
	methods map[Method]string
}

func (b *bitrixClient) Request(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)

	if err != nil {
		b.log.Errorf("error occured bitrix initializng request: %v", err)
		return err
	}

	res, err := b.client.Do(req)

	defer res.Body.Close()

	if err != nil {
		b.log.Errorf("error occured bitrix client do request: %v", err)
		return err
	}

	b.log.Printf("bitrix request status code: %s", res.Status)

	return nil
}

func (b *bitrixClient) CallBackHeader(callBackInfo *dto.RequestCall) string {
	params := url.Values{}
	params.Set("FIELDS[TITLE]", "Заявка на обратный звонок")
	params.Set("FIELDS[NAME]", callBackInfo.Name)
	params.Set("FIELDS[PHONE][0][VALUE]", callBackInfo.Phone)
	url := fmt.Sprintf("https://ufaelectro.bitrix24.ru/rest/1/%s/crm.lead.add.json?%s", b.methods["call_back_header"], params.Encode())
	return url
}

func (b *bitrixClient) CallBackFooter(callBackInfo *dto.RequestCall) string {
	params := url.Values{}
	params.Set("FIELDS[TITLE]", "Заявка на обратный звонок")
	params.Set("FIELDS[PHONE][0][VALUE]", callBackInfo.Phone)
	url := fmt.Sprintf("https://ufaelectro.bitrix24.ru/rest/1/%s/crm.lead.add.json?%s", b.methods["call_back_footer"], params.Encode())
	return url
}

func (b *bitrixClient) CallBackLK(callBackInfo *dto.RequestCall) string {
	params := url.Values{}
	params.Set("FIELDS[TITLE]", "Заявка с личного кабинета")

	if callBackInfo.Inn == "" {
		params.Set("FIELDS[NAME]", fmt.Sprintf("%s %s", callBackInfo.Name, callBackInfo.Surname))
	} else {
		params.Set("FIELDS[NAME]", callBackInfo.ManagerName)
		params.Set("FIELDS[COMPANY_TITLE]", fmt.Sprintf("%s (ИНН: %s)", callBackInfo.CompanyName, callBackInfo.Inn))
	}
	params.Set("FIELDS[PHONE][0][VALUE]", callBackInfo.Phone)
	params.Set("FIELDS[EMAIL][0][VALUE]", callBackInfo.Email)
	params.Set("FIELDS[COMMENTS]", callBackInfo.Message)
	url := fmt.Sprintf("https://ufaelectro.bitrix24.ru/rest/1/%s/crm.lead.add.json?%s", b.methods["call_back_lk"], params.Encode())
	return url
}

func (b *bitrixClient) CallBackVacancy(requestVacancyInfo *dto.RequestVacancy) string {
	params := url.Values{}
	params.Set("FIELDS[TITLE]", fmt.Sprintf("Отклик на вакансию (%s)", requestVacancyInfo.VacancyTitle))
	params.Set("FIELDS[NAME]", requestVacancyInfo.Name)
	params.Set("FIELDS[LAST_NAME]", requestVacancyInfo.Lastname)
	params.Set("FIELDS[SECOND_NAME]", requestVacancyInfo.Surname)
	params.Set("FIELDS[PHONE][0][VALUE]", requestVacancyInfo.Phone)
	params.Set("FIELDS[EMAIL][0][VALUE]", requestVacancyInfo.Email)
	params.Set("FIELDS[COMMENTS]", requestVacancyInfo.Comment)
	url := fmt.Sprintf("https://ufaelectro.bitrix24.ru/rest/1/%s/crm.lead.add.json?%s", b.methods["call_back_vacancy"], params.Encode())
	return url
}
