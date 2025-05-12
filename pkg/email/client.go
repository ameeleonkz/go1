package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"github.com/go-mail/mail/v2"
)

type Client struct {
	dialer *mail.Dialer
	from   string
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewClient(cfg Config) *Client {
	d := mail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	d.TLSConfig = &tls.Config{
		ServerName:         cfg.Host,
		InsecureSkipVerify: false,
	}

	return &Client{
		dialer: d,
		from:   cfg.From,
	}
}

func (c *Client) SendPaymentNotification(to string, amount float64) error {
	content := fmt.Sprintf(`
		<h1>Спасибо за оплату!</h1>
		<p>Сумма: <strong>%.2f RUB</strong></p>
		<small>Это автоматическое уведомление</small>
	`, amount)

	m := mail.NewMessage()
	m.SetHeader("From", c.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Платеж успешно проведен")
	m.SetBody("text/html", content)

	if err := c.dialer.DialAndSend(m); err != nil {
		log.Printf("SMTP error: %v", err)
		return fmt.Errorf("ошибка отправки email: %v", err)
	}

	log.Printf("Email sent to %s", to)
	return nil
}

func (c *Client) SendCreditNotification(to string, creditID int64, amount float64, dueDate string) error {
	content := fmt.Sprintf(`
		<h1>Напоминание о платеже по кредиту</h1>
		<p>Номер кредита: %d</p>
		<p>Сумма платежа: <strong>%.2f RUB</strong></p>
		<p>Дата платежа: %s</p>
		<small>Это автоматическое уведомление</small>
	`, creditID, amount, dueDate)

	m := mail.NewMessage()
	m.SetHeader("From", c.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Напоминание о платеже по кредиту")
	m.SetBody("text/html", content)

	if err := c.dialer.DialAndSend(m); err != nil {
		log.Printf("SMTP error: %v", err)
		return fmt.Errorf("ошибка отправки email: %v", err)
	}

	log.Printf("Email sent to %s", to)
	return nil
} 