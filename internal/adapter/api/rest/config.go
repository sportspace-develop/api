package rest

type Config struct {
	TLSEnable   uint   `env:"TLS_ENABLE" envDefault:"0"`
	TLSCert     string `env:"TLS_CERT" envDefault:""`
	TLSKey      string `env:"TLS_KEY" envDefault:""`
	TLSHosts    string `env:"TLS_HOSTS" envDefault:""`
	TLSDirCache string `env:"TLS_DIR_CACHE"`
}
