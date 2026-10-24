module github.com/centrifugal/centrifuge

go 1.26.0

require (
	github.com/FZambia/eagle v0.2.0
	github.com/centrifugal/protocol v0.9.3
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/igm/sockjs-go/v3 v3.0.3
	github.com/prometheus/client_golang v1.25.0
	github.com/redis/rueidis v1.0.79
	github.com/segmentio/encoding v0.5.4
	github.com/stretchr/testify v1.12.1
	golang.org/x/sync v0.24.0
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.9.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.3 // indirect
	github.com/prometheus/common v0.72.1 // indirect
	github.com/prometheus/procfs v0.23.0 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.49.0 // indirect
)

replace github.com/centrifugal/protocol => github.com/rocketsciencegg/centrifuge-protocol v0.9.3
