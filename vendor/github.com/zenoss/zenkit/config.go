package zenkit

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	LogLevelConfig       = "log.level"
	LogStackdriverConfig = "log.stackdriver"

	TracingEnabledConfig    = "tracing.enabled"
	TracingSampleRateConfig = "tracing.samplerate"

	MetricsEnabledConfig = "metrics.enabled"

	ServiceLabel = "service.label"

	ProfilingEnabledConfig           = "profiling.enabled"
	ProfilingServiceName             = "profiling.service.name"
	ProfilingMutexDisabledConfig     = "profiling.mutex.disabled"
	ProfilingHeapDisabledConfig      = "profiling.heap.disabled"
	ProfilingGoroutineDisabledConfig = "profiling.goroutine.disabled"

	AuthDisabledConfig      = "auth.disabled"
	AuthDevTenantConfig     = "auth.dev_tenant"
	AuthDevEmailConfig      = "auth.dev_email"
	AuthDevUserConfig       = "auth.dev_user"
	AuthDevConnectionConfig = "auth.dev_connection"
	AuthDevScopesConfig     = "auth.dev_scopes"
	AuthDevGroupsConfig     = "auth.dev_groups"
	AuthDevRolesConfig      = "auth.dev_roles"
	AuthDevClientIDConfig   = "auth.dev_clientid"

	GRPCMaxConcurrentRequests = "grpc.max_concurrent_requests"
	GRPCListenAddrConfig      = "grpc.listen_addr"
	GRPCHealthAddrConfig      = "grpc.health_addr"

	GCProjectIDConfig                    = "gcloud.project_id"
	GCDatastoreCredentialsConfig         = "gcloud.datastore.credentials"
	GCEmulatorBigtableConfig             = "gcloud.emulator.bigtable"
	GCEmulatorDatastoreEnabledConfig     = "gcloud.emulator.datastore.enabled"
	GCEmulatorDatastoreHostPortConfig    = "gcloud.emulator.datastore.host_port"
	GCEmulatorPubsubConfig               = "gcloud.emulator.pubsub"
	GCBigtableInstanceIDConfig           = "gcloud.bigtable.instance_id"
	GCBigtableApplicationProfileIDConfig = "gcloud.bigtable.application_profile_id"
	GCPubsubTopicConfig                  = "gcloud.pubsub.topic"
	GCBigtableSuffix                     = "gcloud.bigtable.suffix"
	GCBigtablePoolSize                   = "gcloud.bigtable.poolsize"
	GCMemstoreAddressConfig              = "gcloud.memorystore.address"
	GCMemstoreTTLConfig                  = "gcloud.memorystore.ttl"
	GCMemstoreLocalMaxLen                = "gcloud.memorystore.local_max_len"

	ServiceDialTimeoutConfig = "dial_timeout"

	ZINGAnomalyTableConfig           = "zing.bigtable.table.anomaly"
	ZINGDefinitionIDIndexTableConfig = "zing.bigtable.table.definition_id_index"
	ZINGFieldIndexTableConfig        = "zing.bigtable.table.field_index"
	ZINGItemDefinitionTableConfig    = "zing.bigtable.table.item_definition"
	ZINGItemInstanceTableConfig      = "zing.bigtable.table.item_instance"
	ZINGMetadataTableConfig          = "zing.bigtable.table.metadata"
	ZINGMetricsTableConfig           = "zing.bigtable.table.metrics"
	ZINGRecommendationsTableConfig   = "zing.bigtable.table.recommendations"
	ZINGTrendTableConfig             = "zing.bigtable.table.trend"
	ZINGQueryResultsTableConfig      = "zing.bigtable.table.query_results"

	ZINGProductNameConfig          = "zing.product.name"
	ZINGProductVersionConfig       = "zing.product.version"
	ZINGProductCompanyNameConfig   = "zing.product.company_name"
	ZINGProductOtherCommentsConfig = "zing.product.other_comments"
)

var (
	globalViper = viper.New()

	ErrNoServiceAddress = errors.New("no service address")
)

func init() {
	globalViper.AutomaticEnv()
	globalViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	globalViper.SetDefault(ServiceDialTimeoutConfig, 10*time.Second)

	globalViper.SetDefault(ZINGAnomalyTableConfig, "ANOMALY_V2")
	globalViper.SetDefault(ZINGDefinitionIDIndexTableConfig, "DEFINITION_ID_INDEX_V2")
	globalViper.SetDefault(ZINGFieldIndexTableConfig, "FIELD_INDEX_V2")
	globalViper.SetDefault(ZINGItemDefinitionTableConfig, "ITEM_DEFINITION_V2")
	globalViper.SetDefault(ZINGItemInstanceTableConfig, "ITEM_INSTANCE_V2")
	globalViper.SetDefault(ZINGMetadataTableConfig, "METADATA_V3")
	globalViper.SetDefault(ZINGMetricsTableConfig, "METRICS_V2")
	globalViper.SetDefault(ZINGRecommendationsTableConfig, "RECOMMENDATIONS_V2")
	globalViper.SetDefault(ZINGTrendTableConfig, "TREND_V2")
	globalViper.SetDefault(ZINGQueryResultsTableConfig, "QUERY_RESULTS_V2")

	globalViper.SetDefault(ZINGProductNameConfig, "Zenoss Cloud")
	globalViper.SetDefault(ZINGProductVersionConfig, "1.0")
	globalViper.SetDefault(ZINGProductCompanyNameConfig, "Zenoss")
}

func InitConfig(name string) {
	viper.SetDefault(LogLevelConfig, "info")
	viper.SetDefault(LogStackdriverConfig, true)
	viper.SetDefault(TracingEnabledConfig, true)
	viper.SetDefault(TracingSampleRateConfig, 1.0)
	viper.SetDefault(MetricsEnabledConfig, true)
	viper.SetDefault(AuthDevTenantConfig, "ACME")
	viper.SetDefault(AuthDevUserConfig, "zcuser@acme.example.com")
	viper.SetDefault(AuthDevEmailConfig, "zcuser@acme.example.com")
	viper.SetDefault(AuthDevClientIDConfig, "0123456789abcdef")
	viper.SetDefault(GRPCListenAddrConfig, ":8080")
	viper.SetDefault(GRPCHealthAddrConfig, ":8081")
	viper.SetDefault(GRPCMaxConcurrentRequests, 0)
	viper.SetDefault(GCBigtableInstanceIDConfig, "zenoss-zing-bt1")
	viper.SetDefault(GCProjectIDConfig, "zenoss-zing")
	viper.SetDefault(GCBigtableSuffix, "")
	viper.SetDefault(GCBigtablePoolSize, 4)
	viper.SetDefault(GCBigtableApplicationProfileIDConfig, "")
	viper.SetDefault(GCMemstoreAddressConfig, []string{})
	viper.SetDefault(GCMemstoreTTLConfig, "24h")
	viper.SetDefault(GCMemstoreLocalMaxLen, 1000000)

	viper.SetDefault(ProfilingEnabledConfig, false)
	viper.SetDefault(ProfilingServiceName, name)
	viper.SetDefault(ProfilingMutexDisabledConfig, false)
	viper.SetDefault(ProfilingHeapDisabledConfig, false)
	viper.SetDefault(ProfilingGoroutineDisabledConfig, false)

	viper.SetDefault(ServiceLabel, name)

	viper.SetEnvPrefix(name)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	viper.AutomaticEnv()
}

func ServiceAddress(svc string) (string, error) {
	host := svc
	port := globalViper.GetString(svc + "_SERVICE_PORT")
	if host == "" || port == "" {
		return "", ErrNoServiceAddress
	}
	return fmt.Sprintf("%s:%s", host, port), nil
}

func GlobalConfig() *viper.Viper {
	return globalViper
}
