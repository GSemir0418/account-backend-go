package email

import (
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func Send() {
	m := gomail.NewMessage()
	m.SetHeader("From", "845217811@qq.com")
	m.SetHeader("To", "845217811@qq.com")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>zs</b>!")

	d := gomail.NewDialer(
		viper.GetString("email.smtp.host"),
		viper.GetInt("email.smtp.port"),
		viper.GetString("email.smtp.user"),
		viper.GetString("email.smtp.password"),
	)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
