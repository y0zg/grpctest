package zenkit

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const AuthHeaderScheme = "bearer"

func DevIdentity(ctx context.Context) (context.Context, error) {
	tenant := viper.GetString(AuthDevTenantConfig)
	user := viper.GetString(AuthDevUserConfig)
	email := viper.GetString(AuthDevEmailConfig)
	connection := viper.GetString(AuthDevConnectionConfig)
	scopes := viper.GetStringSlice(AuthDevScopesConfig)
	groups := viper.GetStringSlice(AuthDevGroupsConfig)
	roles := viper.GetStringSlice(AuthDevRolesConfig)
	clientid := viper.GetString(AuthDevClientIDConfig)
	ident := &devTenantIdentity{user, email, scopes, tenant, connection, groups, roles, clientid}
	addIdentityFieldsToTags(ctx, ident)
	return WithTenantIdentity(ctx, ident), nil
}

func UnverifiedIdentity(ctx context.Context) (context.Context, error) {
	// Set up the context
	meta := metautils.ExtractIncoming(ctx)
	ctx = meta.ToOutgoing(ctx)

	// Extract the identity from the metadata
	raw, err := grpc_auth.AuthFromMD(ctx, AuthHeaderScheme)
	if err != nil {
		return nil, wrapUnauthenticated(err)
	}
	ident, err := NewAuth0TenantIdentity(raw)
	if err != nil {
		return nil, wrapUnauthenticated(err)
	}

	addIdentityFieldsToTags(ctx, ident)

	// Add the identity to the context
	return WithTenantIdentity(ctx, ident), nil
}

func addIdentityFieldsToTags(ctx context.Context, ident TenantIdentity) {
	grpc_ctxtags.Extract(ctx).Set(LogTenantField, ident.Tenant()).Set(LogUserField, ident.ID())
}

func wrapUnauthenticated(err error) error {
	stat := status.New(codes.Unauthenticated, err.Error())
	return stat.Err()
}
