package mail

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"skillForce/internal/models"
	"skillForce/pkg/logs"
)

type EmailData struct {
	UserName string
	Url      string
}

type Mail struct {
	from     string
	password string
	host     string
	port     string
}

func NewMail(from string, password string, host string, port string) *Mail {
	return &Mail{
		from:     from,
		password: password,
		host:     host,
		port:     port,
	}
}

func (m *Mail) SendRegMail(ctx context.Context, user *models.User, token string) error {
	auth := smtp.PlainAuth("", m.from, m.password, m.host)

	subject := "Регистрация на платформе SkillForce"

	templatePath := "./../internal/repository/infrastructure/mail/layouts/confirm_mail.html"
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		logs.PrintLog(ctx, "ReadTemplate", err.Error())
		return err
	}

	tmpl, err := template.New("email").Parse(string(tmplBytes))
	if err != nil {
		logs.PrintLog(ctx, "ParseTemplate", err.Error())
		return err
	}

	url := fmt.Sprintf("http://217.16.21.64:8080/api/validEmail?token=%s", token)
	var body bytes.Buffer
	err = tmpl.Execute(&body, EmailData{UserName: user.Name, Url: url})
	if err != nil {
		logs.PrintLog(ctx, "ExecuteTemplate", err.Error())
		return err
	}

	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n", user.Email, m.from, subject)
	msg += "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
	msg += body.String()

	err = smtp.SendMail(fmt.Sprintf("%s:%s", m.host, m.port), auth, m.from, []string{user.Email}, []byte(msg))
	if err != nil {
		logs.PrintLog(ctx, "SendMail", err.Error())
		return err
	}

	logs.PrintLog(ctx, "SendMail", fmt.Sprintf("mail sent to %s", user.Email))
	return nil
}
