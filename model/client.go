package model

import (
	"strings"

	"github.com/ory/fosite"
)

// Client provides the underlying structured make up of an OAuth2.0 Client
type Client struct {
	// ID is the id for this client.
	ID string `bson:"_id" json:"id"`

	// Name is the human-readable string name of the client to be presented to the
	// end-user during authorization.
	Name string `bson:"client_name" json:"client_name"`

	// Secret is the client's secret. The secret will be included in the create request as cleartext, and then
	// never again. The secret is stored using BCrypt so it is impossible to recover it. Tell your users
	// that they need to write the secret down as it will not be made available again.
	Secret string `bson:"client_secret" json:"client_secret,omitempty"`

	// RedirectURIs is an array of allowed redirect urls for the client, for example: http://mydomain/oauth/callback .
	RedirectURIs []string `bson:"redirect_uris" json:"redirect_uris"`

	// GrantTypes is an array of grant types the client is allowed to use.
	//
	// Pattern: client_credentials|authorize_code|implicit|refresh_token
	GrantTypes []string `bson:"grant_types" json:"grant_types"`

	// ResponseTypes is an array of the OAuth 2.0 response type strings that the client can
	// use at the authorization endpoint.
	//
	// Pattern: id_token|code|token
	ResponseTypes []string `bson:"response_types" json:"response_types"`

	// Scope is a string containing a space-separated list of scope values (as
	// described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client
	// can use when requesting access tokens.
	//
	// Pattern: ([a-zA-Z0-9\.]+\s)+
	Scope string `bson:"scope" json:"scope"`

	// Owner is a string identifying the owner of the OAuth 2.0 Client.
	Owner string `bson:"owner" json:"owner"`

	// PolicyURI is a URL string that points to a human-readable privacy policy document
	// that describes how the deployment organization collects, uses,
	// retains, and discloses personal data.
	PolicyURI string `bson:"policy_uri" json:"policy_uri"`

	// TermsOfServiceURI is a URL string that points to a human-readable terms of service
	// document for the client that describes a contractual relationship
	// between the end-user and the client that the end-user accepts when
	// authorizing the client.
	TermsOfServiceURI string `bson:"tos_uri" json:"tos_uri"`

	// ClientURI is an URL string of a web page providing information about the client.
	// If present, the server SHOULD display this URL to the end-user in
	// a clickable fashion.
	ClientURI string `bson:"client_uri" json:"client_uri"`

	// LogoURI is an URL string that references a logo for the client.
	LogoURI string `bson:"logo_uri" json:"logo_uri"`

	// Contacts is a array of strings representing ways to contact people responsible
	// for this client, typically email addresses.
	Contacts []string `bson:"contacts" json:"contacts"`

	// Public is a boolean that identifies this client as public, meaning that it
	// does not have a secret. It will disable the client_credentials grant type for this client if set.
	Public bool `bson:"public" json:"public"`
}

// GetID returns the client's Client ID.
func (c *Client) GetID() string {
	return c.ID
}

// GetRedirectURIs returns the OAuth2.0 authorized Client redirect URIs.
func (c *Client) GetRedirectURIs() []string {
	return c.RedirectURIs
}

// GetHashedSecret returns the Client's Hashed Secret for authenticating with the Identity Provider.
func (c *Client) GetHashedSecret() []byte {
	return []byte(c.Secret)
}

// GetScopes returns an array of strings, wrapped as `fosite.Arguments` to provide functions that allow verifying the
// Client's scopes against incoming requests.
func (c *Client) GetScopes() fosite.Arguments {
	return fosite.Arguments(strings.Split(c.Scope, " "))
}

// GetGrantTypes returns an array of strings, wrapped as `fosite.Arguments` to provide functions that allow verifying
// the Client's Grant Types against incoming requests.
func (c *Client) GetGrantTypes() fosite.Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	//
	// JSON array containing a list of the OAuth 2.0 Grant Types that the Client is declaring
	// that it will restrict itself to using.
	// If omitted, the default is that the Client will use only the authorization_code Grant Type.
	if len(c.GrantTypes) == 0 {
		return fosite.Arguments{"authorization_code"}
	}
	return fosite.Arguments(c.GrantTypes)
}

// GetResponseTypes returns an array of strings, wrapped as `fosite.Arguments` to provide functions that allow verifying
// the Client's Response Types against incoming requests.
func (c *Client) GetResponseTypes() fosite.Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	//
	// <JSON array containing a list of the OAuth 2.0 response_type values that the Client is declaring
	// that it will restrict itself to using. If omitted, the default is that the Client will use
	// only the code Response Type.
	if len(c.ResponseTypes) == 0 {
		return fosite.Arguments{"code"}
	}
	return fosite.Arguments(c.ResponseTypes)
}

// GetOwner returns a string which contains the OAuth Client owner's name.Generally speaking, this will be a developer
// or an organisation.
func (c *Client) GetOwner() string {
	return c.Owner
}

// IsPublic returns a boolean as to whether the Client itself is either private or public. If public, only trusted
// OAuth grant types should be used as client secrets shouldn't be exposed to a public client.
func (c *Client) IsPublic() bool {
	return c.Public
}