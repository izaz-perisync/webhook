package service

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	s3 "webhook/utils/s3"
)

var (
	ErrNotInserted = Generic{Code: 1, Msg: "err not inserted"}
)

type IService interface {
	SubmitForm(ctx context.Context, f SubmitForm, file multipart.File, fileName, formType string) error
}

type Service struct {
	jwtKey, webhooks string

	db  *sql.DB
	loc *time.Location
	l   *log.Logger
	s3  s3.S3Service
}

func New(jwtKey, webhooks string, db *sql.DB, loc *time.Location, l *log.Logger, s3 s3.S3Service) IService {

	return &Service{
		jwtKey:   jwtKey,
		webhooks: webhooks,
		db:       db,
		loc:      loc,
		l:        l,
		s3:       s3,
	}
}

func (s *Service) SubmitForm(ctx context.Context, f SubmitForm, file multipart.File, fileName, formType string) error {
	var (
		r   sql.Result
		err error
	)

	buf, err := io.ReadAll(file)
	if err != nil {
		return Generic{Msg: "invalid file"}
	}

	url, err := s.s3.Upload(ctx, buf, uuid.New().String()+fileName)
	if err != nil {
		return Generic{Msg: "err upload issue"}
	}

	switch formType {

	case "careers":
		r, err = s.db.Exec(`INSERT INTO careers.applications ("name", email, phone_no, created_at, resume, linkedin_url, portfolio_url) VALUES($1, $2, $3, $4, $5, $6, $7);`, f.Name, f.Email, f.MobileNo, time.Now().In(s.loc), url, f.LinkedInUrl, f.PortfolioUrl)

	case "contactUs", "enquiry", "vulnerability":

		r, err = s.db.Exec(`INSERT INTO public.contact_vulnerability_enquiry_form ("name", email, mobile, created_at,company,url,message,form_type) VALUES($1, $2, $3, $4,$5,$6,$7,$8);`, f.Name, f.Email, f.MobileNo, time.Now().In(s.loc), f.Company, url, f.Message, formType)
	default:
		return Generic{Msg: "err invalid request"}
	}

	if err != nil {
		s.l.Print("err in insert", err)
		return err
	}

	c, _ := r.RowsAffected()
	if c == 0 {
		return ErrNotInserted
	}

	return s.WrapDiscord(ctx, f.Name, f.Email+f.MobileNo, s.webhooks)

}

func (s *Service) WrapDiscord(ctx context.Context, header, body, webhookUrl string) error {

	jsonStr := fmt.Sprintf(`{"content":"%s\n%s"}`, strings.ToUpper(header), body)

	_, err := httpCall(jsonStr, webhookUrl)
	if err != nil {

		fmt.Println("call err", err)
		return err
	}

	return nil
}

func httpCall(jsonStr, webhookUrl string) ([]byte, error) {

	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return nil, nil

}
