package mail

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

type KafkaMessage struct {
	Method    string
	Token     string
	UserEmail string
	UserName  string
	CourseId  int
	Url       string
}

type EmailData struct {
	UserName string
	CourseId int
	Url      string
}

type Mail struct {
	from     string
	password string
	host     string
	port     string
	auth     smtp.Auth
}

func NewMail(from string, password string, host string, port string) *Mail {
	return &Mail{
		from:     from,
		password: password,
		host:     host,
		port:     port,
		auth:     smtp.PlainAuth("", from, password, host),
	}
}

func (m *Mail) SendRegMail(ctx context.Context, kafkaMsg KafkaMessage) error {
	subject := "Регистрация на платформе SkillForce"

	templatePath := "./mail/layouts/confirm_mail.html"
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println("SendRegMail", err.Error())
		return err
	}

	tmpl, err := template.New("email").Parse(string(tmplBytes))
	if err != nil {
		fmt.Println("SendRegMail", err.Error())
		return err
	}

	url := fmt.Sprintf("http://217.16.21.64/validate/%s", kafkaMsg.Token)
	var body bytes.Buffer
	err = tmpl.Execute(&body, EmailData{UserName: kafkaMsg.UserName, Url: url})
	if err != nil {
		fmt.Println("SendRegMail", err.Error())
		return err
	}

	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n", kafkaMsg.UserEmail, m.from, subject)
	msg += "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
	msg += body.String()

	err = smtp.SendMail(fmt.Sprintf("%s:%s", m.host, m.port), m.auth, m.from, []string{kafkaMsg.UserEmail}, []byte(msg))
	if err != nil {
		fmt.Println("SendRegMail", err.Error())
		return err
	}

	fmt.Println("SendRegMail", fmt.Sprintf("mail sent to %s", kafkaMsg.UserEmail))
	return nil
}

func (m *Mail) SendWelcomeMail(ctx context.Context, kafkaMsg KafkaMessage) error {
	subject := "Добро пожаловать на платформу SkillForce"

	templatePath := "./mail/layouts/welcome_mail.html"
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println("SendWelcomeMail", err.Error())
		return err
	}

	tmpl, err := template.New("email").Parse(string(tmplBytes))
	if err != nil {
		fmt.Println("SendWelcomeMail", err.Error())
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, EmailData{UserName: kafkaMsg.UserName})
	if err != nil {
		fmt.Println("SendWelcomeMail", err.Error())
		return err
	}

	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n", kafkaMsg.UserEmail, m.from, subject)
	msg += "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
	msg += body.String()

	err = smtp.SendMail(fmt.Sprintf("%s:%s", m.host, m.port), m.auth, m.from, []string{kafkaMsg.UserEmail}, []byte(msg))
	if err != nil {
		fmt.Println("SendWelcomeMail", err.Error())
		return err
	}

	fmt.Println("SendWelcomeMail", fmt.Sprintf("mail sent to %s", kafkaMsg.UserEmail))
	return nil
}

func (m *Mail) SendWelcomeCourseMail(ctx context.Context, kafkaMsg KafkaMessage) error {
	subject := "Продолжайте своё обучение!"

	templatePath := "./mail/layouts/welcome_course_lesson.html"
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println("SendWelcomeCourseMail", err.Error())
		return err
	}

	tmpl, err := template.New("email").Parse(string(tmplBytes))
	if err != nil {
		fmt.Println("SendWelcomeCourseMail", err.Error())
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, EmailData{CourseId: kafkaMsg.CourseId, UserName: kafkaMsg.UserName})
	if err != nil {
		fmt.Println("SendWelcomeCourseMail", err.Error())
		return err
	}

	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n", kafkaMsg.UserEmail, m.from, subject)
	msg += "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
	msg += body.String()

	err = smtp.SendMail(fmt.Sprintf("%s:%s", m.host, m.port), m.auth, m.from, []string{kafkaMsg.UserEmail}, []byte(msg))
	if err != nil {
		fmt.Println("SendWelcomeCourseMail", err.Error())
		return err
	}

	fmt.Println("SendWelcomeCourseMail", fmt.Sprintf("mail sent to %s", kafkaMsg.UserEmail))
	return nil
}

func (m *Mail) SendMiddleCourseMail(ctx context.Context, kafkaMsg KafkaMessage) error {
	subject := "Ты прям огонь, продолжай!"

	templatePath := "./mail/layouts/middle_course_lesson.html"
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println("SenMiddleCourseMail", err.Error())
		return err
	}

	tmpl, err := template.New("email").Parse(string(tmplBytes))
	if err != nil {
		fmt.Println("SendMiddleCourseMail", err.Error())
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, EmailData{CourseId: kafkaMsg.CourseId, UserName: kafkaMsg.UserName})
	if err != nil {
		fmt.Println("SendMiddleCourseMail", err.Error())
		return err
	}

	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n", kafkaMsg.UserEmail, m.from, subject)
	msg += "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
	msg += body.String()

	err = smtp.SendMail(fmt.Sprintf("%s:%s", m.host, m.port), m.auth, m.from, []string{kafkaMsg.UserEmail}, []byte(msg))
	if err != nil {
		fmt.Println("SendMiddleCourseMail", err.Error())
		return err
	}

	fmt.Println("SendMiddleCourseMail", fmt.Sprintf("mail sent to %s", kafkaMsg.UserEmail))
	return nil
}
