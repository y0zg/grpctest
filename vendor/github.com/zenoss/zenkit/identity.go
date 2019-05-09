package zenkit

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"gopkg.in/square/go-jose.v2/jwt"
)

type key int

const (
	identityKey key = iota + 1
)

var (
	ErrorNoSubject    = errors.New("no subject present on token")
	ErrorNoScopes     = errors.New("no scopes present on token")
	ErrorNoTenant     = errors.New("no tenant present on token")
	ErrorNoConnection = errors.New("no connection present on token")
)

// TenantIdentity is an identity in a multi-tenant application
type TenantIdentity interface {
	ID() string
	Email() string
	Scopes() []string
	Tenant() string
	Connection() string
	HasScope(string) bool
	ClientID() string
}

// Cast the identity to IdentityGroups if you need group information from the identity.
// groups := ident.(IdentityGroups)
type IdentityGroups interface {
	Groups() []string
	HasGroup(string) bool
}

// Cast the identity to IdentityRoles if you need role information from the identity.
// roles := ident.(IdentityRoles)
type IdentityRoles interface {
	Roles() []string
	HasRole(string) bool
}

func WithTenantIdentity(ctx context.Context, identity TenantIdentity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}

func ContextTenantIdentity(ctx context.Context) TenantIdentity {
	if v := ctx.Value(identityKey); v != nil {
		return v.(TenantIdentity)
	}
	return nil
}

// NewAuth0TenantIdentity creates an Auth0TenantIdentity for the tokenClaims
func NewAuth0TenantIdentity(token string) (TenantIdentity, error) {
	var claims auth0TenantClaims
	if err := ParseUnverified(token, &claims); err != nil {
		return nil, errors.Wrap(err, "unable to parse token into Auth0 tenant claims")
	}
	if claims.Subject == "" {
		return nil, ErrorNoSubject
	}
	if claims.ScopesValue == "" && claims.ScopeValue == "" {
		return nil, ErrorNoScopes
	}
	if claims.TenantValue == "" {
		return nil, ErrorNoTenant
	}
	if claims.ConnectionValue == "" {
		return nil, ErrorNoConnection
	}
	return &claims, nil
}

type auth0TenantClaims struct {
	jwt.Claims
	ScopeValue      string   `json:"scope,omitempty"`
	ScopesValue     string   `json:"scopes,omitempty"`
	TenantValue     string   `json:"https://dev.zing.ninja/tenant,omitempty"`
	EmailValue      string   `json:"https://dev.zing.ninja/email,omitempty"`
	ConnectionValue string   `json:"https://dev.zing.ninja/connection,omitempty"`
	ClientIDValue   string   `json:"https://dev.zing.ninja/clientid,omitempty"`
	GroupsValue     []string `json:"https://zenoss.com/groups,omitempty"`
	RolesValue      []string `json:"https://zenoss.com/roles,omitempty"`
}

// ID gets the user id for the identity
func (c *auth0TenantClaims) ID() string {
	parts := strings.Split(c.Claims.Subject, "|")
	return parts[len(parts)-1]
}

// Scopes gets the scopes/permissions the identity has been granted
func (c *auth0TenantClaims) Scopes() []string {
	// jwts can contain scope or scopes, let's work with either
	// scope(s) is a space delimited list
	if c.ScopeValue != "" {
		return strings.Split(c.ScopeValue, " ")
	}
	return strings.Split(c.ScopesValue, " ")
}

// Tenant gets the tenant the identity belogs to
func (c *auth0TenantClaims) Tenant() string {
	return c.TenantValue
}

// Tenant gets the tenant the identity belogs to
func (c *auth0TenantClaims) Email() string {
	return c.EmailValue
}

// Connection gets the connection the identity was provided by
func (c *auth0TenantClaims) Connection() string {
	return c.ConnectionValue
}

// HasScope checks if the identity has the scope
func (c *auth0TenantClaims) HasScope(scope string) bool {
	return StringInSlice(scope, c.Scopes())
}

// Groups gets the groups for the identity
func (c *auth0TenantClaims) Groups() []string {
	return c.GroupsValue
}

func (c *auth0TenantClaims) HasGroup(group string) bool {
	return StringInSlice(group, c.Groups())
}

// Groups gets the groups for the identity
func (c *auth0TenantClaims) Roles() []string {
	return c.RolesValue
}

func (c *auth0TenantClaims) HasRole(role string) bool {
	return StringInSlice(role, c.Roles())
}

// Gets the client data for the identity
func (c *auth0TenantClaims) ClientID() string {
	return c.ClientIDValue
}

// StringSliceEquals compare two string slices for equality
func StringSliceEquals(lhs []string, rhs []string) bool {
	if lhs == nil && rhs == nil {
		return true
	}

	if lhs == nil && rhs != nil {
		return false
	}

	if lhs != nil && rhs == nil {
		return false
	}

	if len(lhs) != len(rhs) {
		return false
	}

	for i := range lhs {
		if lhs[i] != rhs[i] {
			return false
		}
	}

	return true
}

// StringInSlice returns whether or not a string is in a string slice
func StringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

type devTenantIdentity struct {
	id         string
	email      string
	scopes     []string
	tenant     string
	connection string
	groups     []string
	roles      []string
	clientid   string
}

func (c *devTenantIdentity) ID() string         { return c.id }
func (c *devTenantIdentity) Email() string      { return c.email }
func (c *devTenantIdentity) Scopes() []string   { return c.scopes }
func (c *devTenantIdentity) Tenant() string     { return c.tenant }
func (c *devTenantIdentity) Connection() string { return c.connection }
func (c *devTenantIdentity) Groups() []string   { return c.groups }
func (c *devTenantIdentity) Roles() []string    { return c.roles}
func (c *devTenantIdentity) ClientID() string   { return c.clientid }

func (c *devTenantIdentity) HasScope(scope string) bool {
	return StringInSlice(scope, c.Scopes())
}

func (c *devTenantIdentity) HasGroup(group string) bool {
	return StringInSlice(group, c.Groups())
}

func (c *devTenantIdentity) HasRole(role string) bool {
	return StringInSlice(role, c.Roles())
}

func IdentityTagsStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = IdentityTaggedContext(stream.Context())
		return handler(srv, wrapped)
	}
}

func IdentityTagsUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(IdentityTaggedContext(ctx), req)
	}
}
