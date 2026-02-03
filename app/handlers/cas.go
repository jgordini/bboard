package handlers

import (
	"net/url"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/cas"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
	webutil "github.com/getfider/fider/app/pkg/web/util"
)

// CASLogin handles the initiation of the CAS login flow
func CASLogin() web.HandlerFunc {
	return func(c *web.Context) error {
		if !cas.IsConfigured() {
			return c.NotFound()
		}

		redirectURL := c.QueryParam("redirect")
		if redirectURL == "" {
			redirectURL = "/"
		}

		cacheKey := "cas_redirect:" + c.SessionID()
		c.Engine().Cache().Set(cacheKey, redirectURL, 10*time.Minute)

		loginURL, err := cas.LoginURL(redirectURL)
		if err != nil {
			log.Error(c, err)
			return c.Redirect("/signin?error=" + url.QueryEscape("CAS login failed"))
		}

		return c.Redirect(loginURL)
	}
}

// CASCallback handles the callback from the CAS server
func CASCallback() web.HandlerFunc {
	return func(c *web.Context) error {
		if !cas.IsConfigured() {
			return c.NotFound()
		}

		ticket := c.QueryParam("ticket")
		if ticket == "" {
			log.Warn(c, "CAS callback received without a ticket")
			return c.Redirect("/signin?error=" + url.QueryEscape("CAS authentication failed: no ticket received"))
		}

		profile, err := cas.ValidateTicket(ticket)
		if err != nil {
			log.Error(c, err)
			return c.Redirect("/signin?error=" + url.QueryEscape("CAS authentication failed: invalid ticket"))
		}

		redirectURL := "/"
		cacheKey := "cas_redirect:" + c.SessionID()
		if val, found := c.Engine().Cache().Get(cacheKey); found {
			if s, ok := val.(string); ok {
				redirectURL = s
			}
			c.Engine().Cache().Delete(cacheKey)
		}

		provider := app.UABProvider

		var user *entity.User
		userByProvider := &query.GetUserByProvider{Provider: provider, UID: profile.ID}
		err = bus.Dispatch(c, userByProvider)
		user = userByProvider.Result

		if errors.Cause(err) == app.ErrNotFound && profile.Email != "" {
			userByEmail := &query.GetUserByEmail{Email: profile.Email}
			err = bus.Dispatch(c, userByEmail)
			user = userByEmail.Result
		}

		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				if c.Tenant().IsPrivate {
					return c.Redirect("/not-invited")
				}
				user = &entity.User{
					Name:   profile.Name,
					Tenant: c.Tenant(),
					Email:  profile.Email,
					Role:   enum.RoleVisitor,
					Providers: []*entity.UserProvider{
						{UID: profile.ID, Name: provider},
					},
				}
				if err = bus.Dispatch(c, &cmd.RegisterUser{User: user}); err != nil {
					log.Error(c, err)
					return c.Redirect("/signin?error=" + url.QueryEscape("Failed to sign in"))
				}
			} else {
				log.Error(c, err)
				return c.Redirect("/signin?error=" + url.QueryEscape("Failed to sign in"))
			}
		} else if !user.HasProvider(provider) {
			if err = bus.Dispatch(c, &cmd.RegisterUserProvider{
				UserID:       user.ID,
				ProviderName: provider,
				ProviderUID:  profile.ID,
			}); err != nil {
				log.Error(c, err)
				return c.Redirect("/signin?error=" + url.QueryEscape("Failed to sign in"))
			}
		}

		webutil.AddAuthUserCookie(c, user)

		return c.Redirect(redirectURL)
	}
}
