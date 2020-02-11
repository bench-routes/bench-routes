package alertways

import (
	// "net/http"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

// Mail ....
type Mail struct {
	Sender  string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
}

// SMTPServer ...
type SMTPServer struct {
	Host      string
	Port      string
	TLSConfig *tls.Config
}

// ServerName ...
func (s *SMTPServer) ServerName() string {
	return s.Host + ":" + s.Port
}

// EmailController controls email operations
func EmailController(to, cc, bcc, sender, subject string) {

	toArr := strings.Split(to, ";")
	ccArr := strings.Split(cc, ";")
	bccArr := strings.Split(bcc, ";")
	body := "I am here"
	request := Mail{
		Sender: sender,
		To: toArr,
		Cc: ccArr,
		Bcc: bccArr,
		Subject: subject,
		Body: body,
	}
	fmt.Println("request: ", request)

	send(request)

}

// BuildMessage ...
func (mail *Mail) BuildMessage() string {
	header := ""
	header += fmt.Sprintf("From: %s\r\n", mail.Sender)
	if len(mail.To) > 0 {
		header += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	}
	if len(mail.Cc) > 0 {
		header += fmt.Sprintf("Cc: %s\r\n", strings.Join(mail.Cc, ";"))
	}

	header += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	header += "\r\n" + mail.Body

	return header
}

func send(mailObject Mail) {

	mail := Mail{}
	mail.Sender = mailObject.Sender
	mail.To = mailObject.To
	mail.Subject = mailObject.Subject
	mail.Body = mailObject.Body

	messageBody := mail.BuildMessage()

	smtpServer := SMTPServer{Host: "smtp.gmail.com", Port: "465"}
	smtpServer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.Host,
	}

	fmt.Println("Sender:: ", mail.Sender)
	auth := smtp.PlainAuth("", mail.Sender, "password", smtpServer.Host)

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), smtpServer.TLSConfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, smtpServer.Host)
	if err != nil {
		log.Panic(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// step 2: add all from and to
	if err = client.Mail(mail.Sender); err != nil {
		log.Panic(err)
	}
	receivers := append(mail.To, mail.Cc...)
	receivers = append(receivers, mail.Bcc...)
	for _, k := range receivers {
		log.Println("sending to: ", k)
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	// Data
	wr, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = wr.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = wr.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")

}
