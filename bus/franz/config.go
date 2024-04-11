package franz

import "time"

type Config struct {
	Addresses             []string      `json:"addresses" yaml:"addresses" env:"ADDRESSES" envSeparator:"," validate:"notEmpty"`
	ReadBackoffMin        time.Duration `json:"read_backoff_min" yaml:"read_backoff_min" env:"READ_BACKOFF_MIN" default:"10ms"`
	ReadBackoffMax        time.Duration `json:"read_backoff_max" yaml:"read_backoff_max" env:"READ_BACKOFF_MAX" default:"100ms"`
	HeartbeatInterval     time.Duration `json:"heartbeat_interval" yaml:"heartbeat_interval" env:"HEARTBEAT_INTERVAL" default:"3s"`
	SessionTimeout        time.Duration `json:"session_timeout" yaml:"session_timeout" env:"SESSION_TIMEOUT" default:"30s"`
	MaxWait               time.Duration `json:"max_wait" yaml:"max_wait" env:"MAX_WAIT" default:"100ms"`
	BatchTimeout          time.Duration `json:"batch_timeout" yaml:"batch_timeout" env:"BATCH_TIMEOUT" default:"30ms"`
	WatchPartitionChanges bool          `json:"watch_partition_changes" yaml:"watch_partition_changes" env:"WATCH_PARTITION_CHANGES" default:"true"`
	FetchMaxBytes         int32         `json:"fetch_max_bytes" yaml:"fetch_max_bytes" env:"FETCH_MAX_BYTES" default:"209715200"`
	BrokerMaxReadBytes    int32         `json:"broker_max_read_bytes" yaml:"broker_max_read_bytes" env:"BROKER_MAX_READ_BYTES" default:"419430400"`
}
