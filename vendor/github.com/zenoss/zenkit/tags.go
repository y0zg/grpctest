package zenkit

import (
	"context"

	"go.opencensus.io/tag"
)

var (
	KeyTenant, _ = tag.NewKey(LogTenantField)
	KeyUser, _   = tag.NewKey(LogUserField)
)

func IdentityTaggedContext(ctx context.Context) context.Context {
	identity := ContextTenantIdentity(ctx)
	ctx, _ = tag.New(ctx, tag.Upsert(KeyTenant, identity.Tenant()), tag.Upsert(KeyUser, identity.ID()))
	return ctx
}
