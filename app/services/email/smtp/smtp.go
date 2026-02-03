package smtp

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"
	gosmtp "net/smtp"
	"net/url"
	"strconv"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/services/email"
)

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "SMTP"
}

func (s Service) Category() string {
	return "email"
}

func (s Service) Enabled() bool {
	return env.Config.Email.Type == "smtp"
}

func (s Service) Init() {
	bus.AddListener(sendMail)
	bus.AddHandler(fetchRecentSupressions)
}

func fetchRecentSupressions(ctx context.Context, c *query.FetchRecentSupressions) error {
	//not implemented for SMTP
	return nil
}

func sendMail(ctx context.Context, c *cmd.SendMail) {
	if c.Props == nil {
		c.Props = dto.Props{}
	}

	if c.From.Address == "" {
		c.From.Address = email.NoReply
	}

	for _, to := range c.To {
		if to.Address == "" {
			return
		}

		u, err := url.Parse(web.BaseURL(ctx))
		localname := "localhost"
		if err == nil {
			localname = u.Hostname()
		}

		if !email.CanSendTo(to.Address) {
			log.Warnf(ctx, "Skipping email to '@{Name} <@{Address}>'.", dto.Props{
				"Name":    to.Name,
				"Address": to.Address,
			})
			return
		}

		log.Debugf(ctx, "Sending email to @{Address} with template @{TemplateName} and params @{Props}.", dto.Props{
			"Address":      to.Address,
			"TemplateName": c.TemplateName,
			"Props":        to.Props,
		})

		message := email.RenderMessage(ctx, c.TemplateName, c.From.Address, c.Props.Merge(to.Props))
		b := builder{}
		b.Set("From", c.From.String())
		b.Set("Reply-To", c.From.Address)
		b.Set("To", to.String())
		b.Set("Subject", email.EncodeSubject(message.Subject))
		b.Set("MIME-version", "1.0")
		b.Set("Content-Type", "text/html; charset=\"UTF-8\"")
		b.Set("Date", time.Now().Format(time.RFC1123Z))
		b.Set("Message-ID", generateMessageID(localname))
		b.Body(message.Body)

		smtpConfig := env.Config.Email.SMTP
		err = sendWithConfig(localname, smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.Password, smtpConfig.EnableStartTLS, to.Address, b.Bytes())
		if err != nil && smtpConfig.BackupHost != "" {
			log.Warnf(ctx, "Primary SMTP failed. Falling back to backup host @{Host}. Error: @{Error}", dto.Props{
				"Host":  smtpConfig.BackupHost,
				"Error": err.Error(),
			})
			err = sendWithConfig(localname, smtpConfig.BackupHost, smtpConfig.BackupPort, smtpConfig.BackupUsername, smtpConfig.BackupPassword, smtpConfig.BackupEnableStartTLS, to.Address, b.Bytes())
		}
		if err != nil {
			panic(errors.Wrap(err, "failed to send email with template %s", c.TemplateName))
		}
		log.Debug(ctx, "Email sent.")
	}
}

var Send = func(localName, serverAddress string, enableStartTLS bool, a gosmtp.Auth, from string, to []string, msg []byte) error {
	host, _, _ := net.SplitHostPort(serverAddress)
	c, err := gosmtp.Dial(serverAddress)
	if err != nil {
		return err
	}
	defer func() { _ = c.Close() }()
	if err = c.Hello(localName); err != nil {
		return err
	}
	if enableStartTLS {
		if ok, _ := c.Extension("STARTTLS"); ok {
			config := &tls.Config{ServerName: host}
			if err = c.StartTLS(config); err != nil {
				return err
			}
		}
	}
	if a != nil {
		if ok, _ := c.Extension("AUTH"); !ok {
			return errors.New("smtp: server doesn't support AUTH")
		}
		if err = c.Auth(a); err != nil {
			return err
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

func sendWithConfig(localname, host, port, username, password string, enableStartTLS bool, to string, msg []byte) error {
	if host == "" {
		return errors.New("smtp: host is empty")
	}
	servername := fmt.Sprintf("%s:%s", host, port)
	auth := authenticate(username, password, host)
	return Send(localname, servername, enableStartTLS, auth, email.NoReply, []string{to}, msg)
}

func generateMessageID(localName string) string {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	randStr := hex.EncodeToString(buf)
	messageID := fmt.Sprintf("<%s.%s@%s>", randStr, timestamp, localName)
	return messageID
}

func authenticate(username string, password string, host string) gosmtp.Auth {
	if username == "" && password == "" {
		return nil
	}
	return AgnosticAuth("", username, password, host)
}

type builder struct {
	content string
}

func (b *builder) Set(key, value string) {
	b.content += fmt.Sprintf("%s: %s\r\n", key, value)
}

func (b *builder) Body(body string) {
	b.content += "\r\n" + body
}

func (b *builder) Bytes() []byte {
	return []byte(b.content)
}
