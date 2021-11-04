package middlewares

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

var htmlBody = `
<html>
<head>
   <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
   <p>You service crashed due to unexpected reasons please check the logs</p>
</body>
`

func PanicHandler(c *gin.Context, err interface{}) {
	go SendMail(err)
	c.JSON(http.StatusInternalServerError,
		gin.H{"error": "Service crashed due to unexpected reason please try again later"})
}

func SendMail(err interface{}) {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("MAIL_FROM"))
	m.SetHeader("To", os.Getenv("MAIL_FROM"))
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Crash Alert-Go Service")
	m.SetBody("text/html", htmlBody)
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASS"))

	// Send the email.
	if err := d.DialAndSend(m); err != nil {
		log.Println(err, "Service crashed and crash alert failed")
	}
}
