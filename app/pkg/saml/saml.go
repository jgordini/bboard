package saml

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/url"
	"strings"
	"time"

	"github.com/crewjam/saml"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
)

// Profile holds user attributes extracted from a SAML assertion
type Profile struct {
	ID    string // NameID or persistent name identifier
	Name  string
	Email string
}

// IsConfigured returns true if SAML is enabled via environment
func IsConfigured() bool {
	c := &env.Config.SAML
	return c.EntityID != "" && c.IdPSSOURL != "" && c.IdPCert != "" && c.SPCertPath != "" && c.SPKeyPath != ""
}

// NewServiceProvider builds a SAML SP from env config and base URL
func NewServiceProvider(metadataURL, acsURL string) (*saml.ServiceProvider, error) {
	c := &env.Config.SAML
	if !IsConfigured() {
		return nil, errors.New("SAML is not configured")
	}

	keyPair, err := tls.LoadX509KeyPair(c.SPCertPath, c.SPKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load SP key pair")
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse SP certificate")
	}

	metaURL, err := url.Parse(metadataURL)
	if err != nil {
		return nil, errors.Wrap(err, "invalid metadata URL")
	}
	acsURLParsed, err := url.Parse(acsURL)
	if err != nil {
		return nil, errors.Wrap(err, "invalid ACS URL")
	}

	idpEntityID := c.IdPEntityID
	if idpEntityID == "" {
		idpEntityID = c.IdPSSOURL
	}

	idpMeta := buildIdPMetadata(idpEntityID, c.IdPSSOURL, c.IdPCert)
	idpCertPEM := strings.TrimSpace(c.IdPCert)

	sp := &saml.ServiceProvider{
		EntityID:    c.EntityID,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		MetadataURL: *metaURL,
		AcsURL:      *acsURLParsed,
		IDPMetadata: idpMeta,
		IDPCertificate: &idpCertPEM,
	}

	return sp, nil
}

// buildIdPMetadata builds a minimal IdP EntityDescriptor from SSO URL and cert
func buildIdPMetadata(entityID, ssoURL, idpCertPEM string) *saml.EntityDescriptor {
	trimmed := strings.TrimSpace(idpCertPEM)
	var certData string
	if strings.HasPrefix(trimmed, "-----") {
		block, _ := pem.Decode([]byte(trimmed))
		if block != nil {
			certData = base64.StdEncoding.EncodeToString(block.Bytes)
		}
	} else {
		certData = strings.ReplaceAll(strings.ReplaceAll(trimmed, "\n", ""), "\r", "")
	}

	keyDescriptors := []saml.KeyDescriptor{}
	if certData != "" {
		keyDescriptors = append(keyDescriptors, saml.KeyDescriptor{
			Use: "signing",
			KeyInfo: saml.KeyInfo{
				X509Data: saml.X509Data{
					X509Certificates: []saml.X509Certificate{{Data: certData}},
				},
			},
		})
	}

	return &saml.EntityDescriptor{
		EntityID: entityID,
		IDPSSODescriptors: []saml.IDPSSODescriptor{
			{
				SSODescriptor: saml.SSODescriptor{
					RoleDescriptor: saml.RoleDescriptor{
						ProtocolSupportEnumeration: "urn:oasis:names:tc:SAML:2.0:protocol",
						KeyDescriptors:            keyDescriptors,
					},
				},
				SingleSignOnServices: []saml.Endpoint{
					{
						Binding:  saml.HTTPRedirectBinding,
						Location: ssoURL,
					},
					{
						Binding:  saml.HTTPPostBinding,
						Location: ssoURL,
					},
				},
			},
		},
	}
}

// LoginURL creates an AuthnRequest, encodes redirect/identifier/requestID in RelayState JWT, and returns the IdP redirect URL
func LoginURL(sp *saml.ServiceProvider, redirect, identifier string) (redirectURL string, err error) {
	idpURL := sp.GetSSOBindingLocation(saml.HTTPRedirectBinding)
	req, err := sp.MakeAuthenticationRequest(idpURL, saml.HTTPRedirectBinding, saml.HTTPPostBinding)
	if err != nil {
		return "", errors.Wrap(err, "failed to create SAML AuthnRequest")
	}
	relayState, err := jwt.Encode(jwt.SAMLStateClaims{
		Redirect:   redirect,
		Identifier: identifier,
		RequestID:  req.ID,
		Metadata:   jwt.Metadata{ExpiresAt: jwt.Time(time.Now().Add(10 * time.Minute))},
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to encode SAML RelayState")
	}
	redirectTo, err := req.Redirect(relayState, sp)
	if err != nil {
		return "", errors.Wrap(err, "failed to build SAML redirect URL")
	}
	return redirectTo.String(), nil
}

// ProfileFromAssertion extracts user profile from a SAML assertion (NameID + attributes)
func ProfileFromAssertion(assertion *saml.Assertion) *Profile {
	p := &Profile{}
	if assertion.Subject != nil && assertion.Subject.NameID != nil {
		p.ID = assertion.Subject.NameID.Value
	}
	for _, stmt := range assertion.AttributeStatements {
		for _, attr := range stmt.Attributes {
			switch {
			case attr.Name == "email" || attr.FriendlyName == "email" || attr.Name == "mail":
				if len(attr.Values) > 0 {
					p.Email = attr.Values[0].Value
				}
			case attr.Name == "displayName" || attr.FriendlyName == "displayName" || attr.Name == "cn" || attr.Name == "urn:oid:2.16.840.1.113730.3.1.241":
				if len(attr.Values) > 0 && p.Name == "" {
					p.Name = attr.Values[0].Value
				}
			case attr.Name == "givenName" || attr.FriendlyName == "givenName":
				if len(attr.Values) > 0 {
					p.Name = strings.TrimSpace(attr.Values[0].Value + " " + p.Name)
				}
			case attr.Name == "sn" || attr.FriendlyName == "sn":
				if len(attr.Values) > 0 {
					p.Name = strings.TrimSpace(p.Name + " " + attr.Values[0].Value)
				}
			}
		}
	}
	if p.Name == "" && p.Email != "" {
		p.Name = strings.Split(p.Email, "@")[0]
	}
	if p.Name == "" {
		p.Name = "SAML User"
	}
	return p
}
