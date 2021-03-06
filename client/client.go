package client

import (
	"github.com/ory/fosite"
)

// Client provides the underlying structured make up of an OAuth2.0 Client. In order to update mongo records efficiently
// omitempty is used for all bson casting, with exception to ID, as this should always be provided in queries and
// updates.
type Client struct {
	// ID is the id for this client.
	ID string `bson:"_id" json:"id" xml:"id"`

	// AllowedTenantAccess contains the Tenant IDs that the client has been given rights to access
	AllowedTenantAccess []string `bson:"allowedTenantAccess" json:"allowedTenantAccess" xml:"allowedTenantAccess"`

	// Name is the human-readable string name of the client to be presented to the
	// end-user during authorization.
	Name string `bson:"clientName" json:"clientName" xml:"clientName"`

	// Secret is the client's secret. The secret will be included in the create request as cleartext, and then
	// never again. The secret is stored using BCrypt so it is impossible to recover it. Tell your users
	// that they need to write the secret down as it will not be made available again.
	Secret []byte `bson:"clientSecret,omitempty" json:"clientSecret,omitempty" xml:"clientSecret,omitempty"`

	// RedirectURIs is an array of allowed redirect urls for the client, for example:
	// http://mydomain/oauth/callback.
	RedirectURIs []string `bson:"redirectUris" json:"redirectUris" xml:"redirectUris"`

	// GrantTypes is an array of grant types the client is allowed to use.
	//
	// Pattern: client_credentials|authorize_code|implicit|refresh_token
	GrantTypes []string `bson:"grantTypes" json:"grantTypes" xml:"grantTypes"`

	// ResponseTypes is an array of the OAuth 2.0 response type strings that the client can
	// use at the authorization endpoint.
	//
	// Pattern: id_token|code|token
	ResponseTypes []string `bson:"responseTypes" json:"responseTypes" xml:"responseTypes"`

	// Scope is a string containing a space-separated list of scope values (as
	// described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client
	// can use when requesting access tokens.
	//
	// Pattern: ([a-zA-Z0-9\.]+\s)+
	Scopes []string `bson:"scopes" json:"scopes" xml:"scopes"`

	// Owner is a string identifying the owner of the OAuth 2.0 Client.
	Owner string `bson:"owner" json:"owner" xml:"owner"`

	// PolicyURI is a URL string that points to a human-readable privacy policy document
	// that describes how the deployment organization collects, uses,
	// retains, and discloses personal data.
	PolicyURI string `bson:"policyUri" json:"policyUri" xml:"policyUri"`

	// TermsOfServiceURI is a URL string that points to a human-readable terms of service
	// document for the client that describes a contractual relationship
	// between the end-user and the client that the end-user accepts when
	// authorizing the client.
	TermsOfServiceURI string `bson:"termsOfServiceUri" json:"termsOfServiceUri" xml:"termsOfServiceUri"`

	// ClientURI is an URL string of a web page providing information about the client.
	// If present, the server SHOULD display this URL to the end-user in
	// a clickable fashion.
	ClientURI string `bson:"clientUri" json:"clientUri" xml:"clientUri"`

	// LogoURI is an URL string that references a logo for the client.
	LogoURI string `bson:"logoUri" json:"logoUri" xml:"logoUri"`

	// Contacts is a array of strings representing ways to contact people responsible
	// for this client, typically email addresses.
	Contacts []string `bson:"contacts" json:"contacts" xml:"contacts"`

	// Public is a boolean that identifies this client as public, meaning that it
	// does not have a secret. It will disable the client_credentials grant type for this client if set.
	Public bool `bson:"public" json:"public" xml:"public"`

	// Disabled is a boolean that identifies whether the client has had it's access disabled.
	Disabled bool `bson:"disabled" json:"disabled" xml:"disabled"`
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
	return c.Secret
}

// GetScopes returns an array of strings, wrapped as `fosite.Arguments` to provide functions that allow verifying the
// Client's scopes against incoming requests.
func (c *Client) GetScopes() fosite.Arguments {
	return fosite.Arguments(c.Scopes)
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

// Disabled returns a boolean as to whether the Client itself has had it's access disabled.
func (c *Client) IsDisabled() bool {
	return c.Disabled
}

// EnableScopeAccess enables client scope access
func (c *Client) EnableScopeAccess(scopes ...string) {
	for i := range scopes {
		found := false
		for j := range c.Scopes {
			if scopes[i] == c.Scopes[j] {
				found = true
				break
			}
		}
		if !found {
			c.Scopes = append(c.Scopes, scopes[i])
		}
	}
}

// DisableScopeAccess disables client scope access.
func (c *Client) DisableScopeAccess(scopes ...string) {
	for i := range scopes {
		for j := range c.Scopes {
			if scopes[i] == c.Scopes[j] {
				copy(c.Scopes[j:], c.Scopes[j+1:])
				c.Scopes[len(c.Scopes)-1] = ""
				c.Scopes = c.Scopes[:len(c.Scopes)-1]
				break
			}
		}
	}
}

// EnableTenantAccess adds a single or multiple tenantIDs to the given client
func (c *Client) EnableTenantAccess(tenantIDs ...string) {
	for i := range tenantIDs {
		found := false
		for j := range c.AllowedTenantAccess {
			if tenantIDs[i] == c.AllowedTenantAccess[j] {
				found = true
				break
			}
		}
		if !found {
			c.AllowedTenantAccess = append(c.AllowedTenantAccess, tenantIDs[i])
		}
	}
}

// RemoveTenants removes a single or multiple tenantIDs from the given client
func (c *Client) DisableTenantAccess(tenantIDs ...string) {
	for i := range tenantIDs {
		for j := range c.AllowedTenantAccess {
			if tenantIDs[i] == c.AllowedTenantAccess[j] {
				copy(c.AllowedTenantAccess[j:], c.AllowedTenantAccess[j+1:])
				c.AllowedTenantAccess[len(c.AllowedTenantAccess)-1] = ""
				c.AllowedTenantAccess = c.AllowedTenantAccess[:len(c.AllowedTenantAccess)-1]
				break
			}
		}
	}
}

// Equal enables checking equality as having a byte array in a struct stops allowing equality checks.
func (c Client) Equal(x Client) bool {
	if c.ID != x.ID {
		return false
	}

	if !strArrEquals(c.AllowedTenantAccess, x.AllowedTenantAccess) {
		return false
	}

	if c.Name != x.Name {
		return false
	}

	if string(c.Secret) != string(x.Secret) {
		return false
	}

	if !strArrEquals(c.RedirectURIs, x.RedirectURIs) {
		return false
	}

	if !strArrEquals(c.GrantTypes, x.GrantTypes) {
		return false
	}

	if !strArrEquals(c.ResponseTypes, x.ResponseTypes) {
		return false
	}

	if !strArrEquals(c.Scopes, x.Scopes) {
		return false
	}

	if c.Owner != x.Owner {
		return false
	}

	if c.PolicyURI != x.PolicyURI {
		return false
	}

	if c.TermsOfServiceURI != x.TermsOfServiceURI {
		return false
	}

	if c.ClientURI != x.ClientURI {
		return false
	}

	if c.LogoURI != x.LogoURI {
		return false
	}

	if !strArrEquals(c.Contacts, x.Contacts) {
		return false
	}

	if c.Public != x.Public {
		return false
	}

	if c.Disabled != x.Disabled {
		return false
	}

	return true
}

func (c Client) IsEmpty() bool {
	return c.Equal(Client{})
}

func strArrEquals(arr1 []string, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}
