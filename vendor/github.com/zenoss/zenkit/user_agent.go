package zenkit

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc/metadata"
)

var (
	userAgentOnce sync.Once
	userAgent     string
)

func setUserAgent() {
	globalConfig := GlobalConfig()
	name := globalConfig.GetString(ZINGProductNameConfig)
	version := globalConfig.GetString(ZINGProductVersionConfig)
	companyName := globalConfig.GetString(ZINGProductCompanyNameConfig)
	otherComments := globalConfig.GetString(ZINGProductOtherCommentsConfig)
	userAgent = fmt.Sprintf("%s/%s (GPN:%s;%s)", name, version, companyName, otherComments)
}

func getUserAgent() string {
	userAgentOnce.Do(setUserAgent)
	return userAgent
}

func WithUserAgent(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "user-agent", getUserAgent())
}
