package centrifuge

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var errRedisStartup = errors.New("temporary subscription failure")

// runRedisBrokerWithStartupResults runs a broker whose startup channels already
// hold the given results, without launching Redis connections.
func runRedisBrokerWithStartupResults(t *testing.T, controlErr, pubSubErr error) []LogEntry {
	t.Helper()

	var entries []LogEntry
	node, err := New(Config{
		LogLevel: LogLevelError,
		LogHandler: func(entry LogEntry) {
			entries = append(entries, entry)
		},
	})
	require.NoError(t, err)

	controlCh := make(chan error, 1)
	controlCh <- controlErr
	pubSubCh := make(chan error, 1)
	pubSubCh <- pubSubErr
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
	close(broker.closeCh)

	require.NoError(t, broker.Run(nil))
	require.Empty(t, controlCh)
	require.Empty(t, pubSubCh)
	return entries
}

func requireRedisStartupErrorLogged(t *testing.T, entry LogEntry, message string) {
	t.Helper()
	require.Equal(t, LogLevelError, entry.Level)
	require.Equal(t, message, entry.Message)
	require.Equal(t, errRedisStartup.Error(), entry.Fields["error"])
	require.Equal(t, 0, entry.Fields["shard"])
}

func TestRedisBrokerRunStartupSuccess(t *testing.T) {
	entries := runRedisBrokerWithStartupResults(t, nil, nil)
	require.Empty(t, entries)
}

func TestRedisBrokerRunLogsControlStartupError(t *testing.T) {
	entries := runRedisBrokerWithStartupResults(t, errRedisStartup, nil)
	require.Len(t, entries, 1)
	requireRedisStartupErrorLogged(t, entries[0], "control pub/sub startup error")
}

func TestRedisBrokerRunLogsPubSubStartupError(t *testing.T) {
	entries := runRedisBrokerWithStartupResults(t, nil, errRedisStartup)
	require.Len(t, entries, 1)
	requireRedisStartupErrorLogged(t, entries[0], "pub/sub startup error")
	require.Equal(t, 0, entries[0].Fields["clusterShard"])
	require.Equal(t, 0, entries[0].Fields["pubSubShard"])
}

func TestRedisBrokerRunLogsAllStartupErrors(t *testing.T) {
	entries := runRedisBrokerWithStartupResults(t, errRedisStartup, errRedisStartup)
	require.Len(t, entries, 2)
	requireRedisStartupErrorLogged(t, entries[0], "control pub/sub startup error")
	requireRedisStartupErrorLogged(t, entries[1], "pub/sub startup error")
}
