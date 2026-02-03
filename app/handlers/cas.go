package handlers

import (
	"net/http"
	"net/url"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/cas"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/web"
)

// CASLogin handles the initiation of the CAS login flow
func CASLogin() web.Handler {
	return func(c *app.Context) error {
		if !cas.IsConfigured() {
			return errors.NotFound(c.Request.URL.Path)
		}

		redirectURL := c.QueryParam("redirect")
		if redirectURL == "" {
			redirectURL = "/"
		}

		// Store the redirect URL in a session or cookie for later use in the callback
		c.Session().Set("cas_redirect_url", redirectURL)

		loginURL, err := cas.LoginURL(redirectURL)
		if err != nil {
			c.Log().Error("Failed to build CAS login URL: %v", err)
			return c.Redirect(http.StatusTemporaryRedirect, "/signin?error="+url.QueryEscape("CAS login failed"))
		}

		return c.Redirect(http.StatusTemporaryRedirect, loginURL)
	}
}

// CASCallback handles the callback from the CAS server
func CASCallback() web.Handler {
	return func(c *app.Context) error {
		if !cas.IsConfigured() {
			return errors.NotFound(c.Request.URL.Path)
		}

		ticket := c.QueryParam("ticket")
		if ticket == "" {
			c.Log().Warn("CAS callback received without a ticket")
			return c.Redirect(http.StatusTemporaryRedirect, "/signin?error="+url.QueryEscape("CAS authentication failed: no ticket received"))
		}

		profile, err := cas.ValidateTicket(ticket)
		if err != nil {
			c.Log().Error("Failed to validate CAS ticket: %v", err)
			return c.Redirect(http.StatusTemporaryRedirect, "/signin?error="+url.QueryEscape("CAS authentication failed: invalid ticket"))
		}

		redirectURL := "/"
		if val := c.Session().GetString("cas_redirect_url"); val != "" {
			redirectURL = val
			c.Session().Remove("cas_redirect_url")
		}

		ctx := c.Request.Context()

		// Use the existing user management logic
		signinCmd := cmd.GetOrCreateUserFromProvider{
			Provider:  app.UABProvider, // Using "uab" as provider for both SAML and CAS
			Reference: profile.ID,
			Email:     profile.Email,
			Name:      profile.Name,
		}

		// Ensure we are in a transaction
		if err := c.Service().Execute(ctx, signinCmd); err != nil {
			c.Log().Error("Failed to get or create user from CAS provider: %v", err)
			return c.Redirect(http.StatusTemporaryRedirect, "/signin?error="+url.QueryEscape("Failed to sign in"))
		}

		uabUser, ok := signinCmd.Result.(*dto.User)
		if !ok {
			c.Log().Error("Failed to cast result to User DTO after CAS authentication")
			return c.Redirect(http.StatusTemporaryRedirect, "/signin?error="+url.QueryEscape("Failed to sign in"))
		}

		if uabUser == nil {
			c.Log().Error("Received nil user after CAS authentication")
			return c.Redirect(http.StatusTemporaryRedirect, "/signin?error="+url.QueryEscape("Failed to sign in"))
		}

		// Log the user in
		err = c.SignIn(uabUser)
		if err != nil {
			c.Log().Error("Failed to sign in user after CAS authentication: %v", err)
			return c.Redirect(http.StatusTemporaryRedirect, "/signin?error="+url.QueryEscape("Failed to sign in"))
		}

		return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	}
}
