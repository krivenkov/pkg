package builder

import (
	"testing"
	"time"

	"github.com/krivenkov/pkg/bus/franz"
	"github.com/krivenkov/pkg/mcfg"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("env", func(t *testing.T) {
		t.Setenv("TRANSPORT", string(TransportFranz))

		t.Setenv("FRANZ_ADDRESSES", "kafka01.dev.bred.team:9092,kafka02.dev.bred.team:9092,kafka03.dev.bred.team:9092")
		t.Setenv("FRANZ_READ_BACKOFF_MIN", "1s")
		t.Setenv("FRANZ_READ_BACKOFF_MAX", "2s")
		t.Setenv("FRANZ_HEARTBEAT_INTERVAL", "3s")
		t.Setenv("FRANZ_SESSION_TIMEOUT", "4s")
		t.Setenv("FRANZ_MAX_WAIT", "5s")
		t.Setenv("FRANZ_BATCH_TIMEOUT", "6s")
		t.Setenv("FRANZ_WATCH_PARTITION_CHANGES", "false")
		t.Setenv("FRANZ_FETCH_MAX_BYTES", "100")
		t.Setenv("FRANZ_BROKER_MAX_READ_BYTES", "200")

		var cfg Config
		err := mcfg.LoadFromEnv(&cfg)
		require.NoError(t, err)
		require.Equal(t, Config{
			Transport: TransportFranz,
			Franz: franz.Config{
				Addresses:             []string{"kafka01.dev.bred.team:9092", "kafka02.dev.bred.team:9092", "kafka03.dev.bred.team:9092"},
				ReadBackoffMin:        time.Second,
				ReadBackoffMax:        time.Second * 2,
				HeartbeatInterval:     time.Second * 3,
				SessionTimeout:        time.Second * 4,
				MaxWait:               time.Second * 5,
				BatchTimeout:          time.Second * 6,
				WatchPartitionChanges: false,
				FetchMaxBytes:         100,
				BrokerMaxReadBytes:    200,
			},
		}, cfg)
	})

	t.Run("defaults", func(t *testing.T) {
		t.Setenv("TRANSPORT", string(TransportFranz))
		t.Setenv("FRANZ_ADDRESSES", "addr1,addr2")

		var cfg Config
		err := mcfg.LoadFromEnv(&cfg)
		require.NoError(t, err)
		require.Equal(t, Config{
			Transport: TransportFranz,
			Franz: franz.Config{
				Addresses:             []string{"addr1", "addr2"},
				ReadBackoffMin:        time.Millisecond * 10,
				ReadBackoffMax:        time.Millisecond * 100,
				HeartbeatInterval:     time.Second * 3,
				SessionTimeout:        time.Second * 30,
				MaxWait:               time.Millisecond * 100,
				BatchTimeout:          time.Millisecond * 30,
				WatchPartitionChanges: true,
				FetchMaxBytes:         209715200,
				BrokerMaxReadBytes:    419430400,
			},
		}, cfg)
	})
}
