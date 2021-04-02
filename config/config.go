package config

type Config struct {
	App  `yaml:"app"`
	HTTP `yaml:"http"`
	Log  `yaml:"logger"`
	PG   `yaml:"postgres"`
	RMQ  `yaml:"rabbitmq"`
}

type App struct {
	Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
	Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
}

type HTTP struct {
	Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
}

type Log struct {
	ZapLevel     string `env-required:"true" yaml:"zap_level"   env:"LOG_ZAP_LEVEL"`
	RollbarEnv   string `env-required:"true" yaml:"rollbar_env" env:"LOG_ROLLBAR_ENV"`
	RollbarToken string `env-required:"true"                    env:"LOG_ROLLBAR_TOKEN"`
}

type PG struct {
	PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
	URL     string `env-required:"true"                 env:"PG_URL"`
}

type RMQ struct {
	ServerExchange string `env-required:"true" yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
	ClientExchange string `env-required:"true" yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
	URL            string `env-required:"true"                            env:"RMQ_URL"`
}
