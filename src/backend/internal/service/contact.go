package service

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/config"
	"main/internal/model"
	"net/http"
	"net/url"
	"strings"
)

type ContactService struct {
	config *config.Config
}

func NewContactService(cfg *config.Config) *ContactService {
	return &ContactService{config: cfg}
}

func (cs *ContactService) verifyRecaptcha(token string) (bool, error) {
	const verifyURL = "https://www.google.com/recaptcha/api/siteverify"
	data := url.Values{}
	data.Set("secret", cs.config.RecaptchaSecretKey)
	data.Set("response", token)

	res, err := http.Post(verifyURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return false, fmt.Errorf("failed reCAPTCHA API request: %v", err)
	}
	defer res.Body.Close()

	var recaptchaRes model.RecaptchaResponse
	if err = json.NewDecoder(res.Body).Decode(&recaptchaRes); err != nil {
		return false, err
	}

	const minScore = 0.5
	if !recaptchaRes.Success || recaptchaRes.Score < minScore || recaptchaRes.Action != "contact" {
		return false, nil
	}
	return true, nil
}

func (cs *ContactService) ProcessContactRequest(ctx context.Context, contactReq *model.ContactRequest) error {
	good, err := cs.verifyRecaptcha(contactReq.Recaptcha)
	if err != nil {
		return err
	}
	if !good {
		return fmt.Errorf("failed reCAPTCHA: %v", err)
	}
	return nil
}
