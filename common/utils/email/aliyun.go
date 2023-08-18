package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

func Send(port int, host, sendMail, password, sendName, toEmail, title, body string) {

	header := make(map[string]string)
	if sendName != "" {
		header["From"] = sendName + "<" + sendMail + ">"
	} else {
		header["From"] = "<" + sendMail + ">"
	}
	header["To"] = toEmail
	header["Subject"] = title
	header["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth(
		"",
		sendMail,
		password,
		host,
	)

	err := SendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		sendMail,
		[]string{toEmail},
		[]byte(message),
	)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Send mail success!")
	}
}

// return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// 参考net/smtp的func SendMail()
// 使用net.Dial连接tls（SSL）端口时，smtp.NewClient()会卡住且不提示err
// len(to)>1时，to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

func main() {
	//Send("", "", "test mail", "测试邮件，你的验证码123123哈哈哈哈")
}
