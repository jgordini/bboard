package cas

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
)

// Profile holds the user information returned by CAS
type Profile struct {
	ID    string
	Email string
	Name  string
}

// IsConfigured returns true if CAS is configured
func IsConfigured() bool {
	return env.Config.CAS.ServerURL != ""
}

// LoginURL builds the CAS login redirect URL
func LoginURL(redirectURL string) (string, error) {
	if !IsConfigured() {
		return "", errors.New("CAS is not configured")
	}

	serviceURL := env.Config.CAS.ServiceURL
	if serviceURL == "" {
		serviceURL = env.Config.BaseURL
	}

	loginURL, err := url.Parse(env.Config.CAS.ServerURL + "/login")
	if err != nil {
		return "", errors.Wrap(err, "failed to parse CAS Server URL")
	}

	q := loginURL.Query()
	q.Set("service", serviceURL+"/cas/callback?redirect="+url.QueryEscape(redirectURL))
	loginURL.RawQuery = q.Encode()

	return loginURL.String(), nil
}

// ValidateTicket validates the CAS ticket and returns the user profile
func ValidateTicket(ticket string) (*Profile, error) {
	if !IsConfigured() {
		return nil, errors.New("CAS is not configured")
	}

	serviceURL := env.Config.CAS.ServiceURL
	if serviceURL == "" {
		serviceURL = env.Config.BaseURL
	}

	validateURL, err := url.Parse(env.Config.CAS.ServerURL + "/serviceValidate")
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse CAS Server URL")
	}

	q := validateURL.Query()
	q.Set("ticket", ticket)
	q.Set("service", serviceURL+"/cas/callback")
	validateURL.RawQuery = q.Encode()

	resp, err := http.Get(validateURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate CAS ticket")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("CAS ticket validation failed with status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CAS validation response body")
	}

	var casResponse casServiceResponse
	err = xml.Unmarshal(body, &casResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal CAS validation response: %s", string(body))
	}

	if casResponse.ServiceResponse.AuthenticationSuccess.User == "" {
		return nil, errors.New("CAS authentication failed: %s", string(body))
	}

	username := strings.ToLower(casResponse.ServiceResponse.AuthenticationSuccess.User)

	profile := &Profile{
		ID:    username,
		Email: username + "@uab.edu",
		Name:  username,
	}

	return profile, nil
}

// casServiceResponse represents the XML structure of the CAS service validation response (CAS 2.0)
type casServiceResponse struct {
	XMLName       xml.Name `xml:"serviceResponse"`
	ServiceResponse struct {
		AuthenticationSuccess struct {
			User string `xml:"user"`
		} `xml:"authenticationSuccess"`
	}
}
