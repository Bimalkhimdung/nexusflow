module github.com/nexusflow/nexusflow/services/search-service

go 1.24.0

toolchain go1.24.6

require (
	github.com/elastic/go-elasticsearch/v8 v8.11.0
	github.com/nexusflow/nexusflow/pkg/config v0.0.0
	github.com/nexusflow/nexusflow/pkg/logger v0.0.0
	github.com/nexusflow/nexusflow/pkg/proto v0.0.0
	google.golang.org/grpc v1.77.0
)

require (
	github.com/elastic/elastic-transport-go/v8 v8.3.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.1 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.18.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20240119083558-1b970713d09a // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251111163417-95abcf5c77ba // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/nexusflow/nexusflow/pkg/config => ../../pkg/config
	github.com/nexusflow/nexusflow/pkg/kafka => ../../pkg/kafka
	github.com/nexusflow/nexusflow/pkg/logger => ../../pkg/logger
	github.com/nexusflow/nexusflow/pkg/proto => ../../pkg/proto
)
