package handlers

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/saml"
	"github.com/getfider/fider/app/pkg/web"
	webutil "github.com/getfider/fider/app/pkg/web/util"
)

// SAMLLogin initiates SAML auth by redirecting to the IdP (e.g. UAB Shibboleth)
func SAMLLogin() web.HandlerFunc {
	return func(c *web.Context) error {
		c.Response.Header().Add("X-Robots-Tag", "noindex")

		if !saml.IsConfigured() {
			return c.NotFound()
		}

		redirect := c.QueryParam("redirect")
		if redirect == "" {
			redirect = c.BaseURL()
		} else if redirect != c.BaseURL() && !strings.HasPrefix(redirect, c.BaseURL()+"/") {
			return c.Forbidden()
		}

		redirectURL, _ := url.ParseRequestURI(redirect)
		redirectURL.ResolveReference(c.Request.URL)

		if c.IsAuthenticated() {
			return c.Redirect(redirect)
		}

		baseURL := web.OAuthBaseURL(c)
		metadataURL := baseURL + "/saml/metadata"
		acsURL := baseURL + "/saml/acs"

		sp, err := saml.NewServiceProvider(metadataURL, acsURL)
		if err != nil {
			log.Error(c, err)
			return c.Failure(err)
		}

		loginURL, err := saml.LoginURL(sp, redirect, c.SessionID())
		if err != nil {
			log.Error(c, err)
			return c.Failure(err)
		}

		return c.Redirect(loginURL)
	}
}

// SAMLACS handles the SAML Assertion Consumer Service (POST from IdP)
func SAMLACS() web.HandlerFunc {
	return func(c *web.Context) error {
		c.Response.Header().Add("X-Robots-Tag", "noindex")

		if !saml.IsConfigured() {
			return c.NotFound()
		}

		baseURL := web.OAuthBaseURL(c)
		metadataURL := baseURL + "/saml/metadata"
		acsURL := baseURL + "/saml/acs"

		sp, err := saml.NewServiceProvider(metadataURL, acsURL)
		if err != nil {
			log.Error(c, err)
			return c.Failure(err)
		}

		// Restore body so ParseResponse can read PostForm (WrapRequest consumed it)
		rawReq := c.Request.Unwrap()
		rawReq.Body = io.NopCloser(strings.NewReader(c.Request.Body))
		if err := rawReq.ParseForm(); err != nil {
			return c.Failure(err)
		}
		relayState := rawReq.PostForm.Get("RelayState")
		if relayState == "" {
			return c.Forbidden()
		}

		claims, err := jwt.DecodeSAMLStateClaims(relayState)
		if err != nil {
			return c.Forbidden()
		}

		if claims.Identifier != c.SessionID() {
			log.Warn(c, "SAML RelayState identifier doesn't match session. Aborting.")
			return c.Forbidden()
		}

		possibleRequestIDs := []string{claims.RequestID}
		assertion, err := sp.ParseResponse(rawReq, possibleRequestIDs)
		if err != nil {
			log.Error(c, err)
			return c.Failure(err)
		}

		profile := saml.ProfileFromAssertion(assertion)
		if profile.ID == "" {
			return c.Failure(errors.New("SAML assertion missing NameID"))
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
					return c.Failure(err)
				}
			} else {
				return c.Failure(err)
			}
		} else if !user.HasProvider(provider) {
			if err = bus.Dispatch(c, &cmd.RegisterUserProvider{
				UserID:       user.ID,
				ProviderName: provider,
				ProviderUID:  profile.ID,
			}); err != nil {
				return c.Failure(err)
			}
		}

		webutil.AddAuthUserCookie(c, user)

		redirectURL := claims.Redirect
		if redirectURL == "" {
			redirectURL = c.BaseURL()
		}
		return c.Redirect(redirectURL)
	}
}

// SAMLMetadata serves the SP metadata XML for IdP configuration
func SAMLMetadata() web.HandlerFunc {
	return func(c *web.Context) error {
		if !saml.IsConfigured() {
			return c.NotFound()
		}

		baseURL := web.OAuthBaseURL(c)
		metadataURL := baseURL + "/saml/metadata"
		acsURL := baseURL + "/saml/acs"

		sp, err := saml.NewServiceProvider(metadataURL, acsURL)
		if err != nil {
			log.Error(c, err)
			return c.Failure(err)
		}

		md := sp.Metadata()
		xmlBytes, err := xml.MarshalIndent(md, "", "  ")
		if err != nil {
			return c.Failure(err)
		}

		c.Response.Header().Set("Content-Type", "application/samlmetadata+xml; charset=utf-8")
		return c.Blob(http.StatusOK, "application/samlmetadata+xml", xmlBytes)
	}
}
