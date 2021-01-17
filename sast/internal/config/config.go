package config

const Version = "0.0.1"

type Config struct {
	Debug      bool
	HTTPPort   uint64
	StaticPath string
}
