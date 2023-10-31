package utils

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"os"
	"path/filepath"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/initializers"
	"github.com/k3a/html2text"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(userEmail string, data *EmailData, emailTemp string) error {
	config, err := initializers.LoadConfig(os.Getenv("API_ENV_CONFIG_PATH"))

	if err != nil {
		zap.L().Error("could not load config ", zap.Error(err))
		return err
	}

	// Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := userEmail
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer

	template, err := ParseTemplateDir(os.Getenv("GIST_EMAIL_TEMPLATE_DIR"))
	if err != nil {
		zap.L().Error("could not parse template directory ", zap.Error(err))
		return err
	}

	err = template.ExecuteTemplate(&body, emailTemp, &data)
	if err != nil {
		zap.L().Error("Could not execute template ", zap.Error(err))
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err = d.DialAndSend(m); err != nil {
		zap.L().Error("Could not send email: ", zap.Error(err))
		return err
	}

	return nil
}
