package config

func NewConfig() *Config {
	c := &Config{
		Global:   NewGlobal(),
		Default:  NewDefault(),
		Frontend: make([]*Frontend, 0),
		Backend:  make([]*Backend, 0),
	}
	return c
}

type Config struct {
	Global   *Global
	Default  *Default
	Frontend []*Frontend
	Backend  []*Backend
}
