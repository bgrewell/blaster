package internal

type TcpFlow struct {
	StartTime      int64  `json:"start_time,omitempty"`
	Duration       int64  `json:"duration,omitempty"`
	PacketSize     int    `json:"packet_size,omitempty"`
	RateBitsPerSec int64  `json:"rate_bits_per_sec,omitempty"`
	Scheduler      string `json:"scheduler,omitempty"`
}
