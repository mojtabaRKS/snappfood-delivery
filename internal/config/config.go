package config

type AppEnv string

const (
	ProductionEnv AppEnv = "production"
	StageEnv      AppEnv = "stage"
	DevelopEnv    AppEnv = "develop"
	LocalEnv      AppEnv = "local"

	ORDER_QUEUE_KEY = "ORDER_QUEUE"
)

type (
	Config struct {
		AppEnv   AppEnv
		HTTP     HTTP
		Mysql    Mysql
		Redis    Redis
		TimeZone string
		AppDebug bool
	}

	HTTP struct {
		Host string
		Port int
	}

	Mysql struct {
		Host     string
		Port     string
		UserName string
		Password string
		Database string
	}

	Redis struct {
		Host     string
		Port     string
		Password string
		Database int
	}
)
