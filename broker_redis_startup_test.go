package centrifuge

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedisBrokerRunStartupErrors(t *testing.T) {
	startupErr := errors.New("temporary subscription failure")
	for _, tt := range []struct {
		name       string
		controlErr error
		pubSubErr  error
		logs       int
	}{
		{name: "success"},
		{name: "control failure", controlErr: startupErr, logs: 1},
		{name: "pubsub failure", pubSubErr: startupErr, logs: 1},
		{name: "both failures", controlErr: startupErr, pubSubErr: startupErr, logs: 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var entries []LogEntry
			node, err := New(Config{
				LogLevel: LogLevelError,
				LogHandler: func(entry LogEntry) {
					entries = append(entries, entry)
				},
			})
			require.NoError(t, err)

			controlCh := make(chan error, 1)
			controlCh <- tt.controlErr
			pubSubCh := make(chan error, 1)
			pubSubCh <- tt.pubSubErr
			broker := &RedisBroker{
				node:   node,
				config: RedisBrokerConfig{UseLists: true},
				shards: []*shardWrapper{{
					shard:              &RedisShard{},
					controlPubSubStart: &controlPubSubStart{errCh: controlCh},
					pubSubStartChannels: [][]*pubSubStart{{
						{errCh: pubSubCh},
					}},
				}},
				closeCh: make(chan struct{}),
			}
			// Supply completed startup attempts without launching Redis connections.
			close(broker.closeCh)

			require.NoError(t, broker.Run(nil))
			require.Empty(t, controlCh)
			require.Empty(t, pubSubCh)
			require.Len(t, entries, tt.logs)
			for _, entry := range entries {
				require.Equal(t, LogLevelError, entry.Level)
				require.Equal(t, startupErr.Error(), entry.Fields["error"])
				require.Equal(t, 0, entry.Fields["shard"])
			}
		})
	}
}
