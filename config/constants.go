package config

const (
	// VERSION is the current cli version
	VERSION = "dev"
	// AuthHost is the production auth URL
	AuthHost = "https://api.catalyze.io"
	// AuthHostVersion is the version path for the auth host
	AuthHostVersion = ""
	// PaasHost is the production PaaS URL
	PaasHost = "https://paas-api.catalyze.io"
	// PaasHostVersion is the version path for the PaaS host
	PaasHostVersion = ""

	// AuthHostEnvVar is the env variable used to override AuthHost
	AuthHostEnvVar = "AUTH_HOST"
	// PaasHostEnvVar is the env variable used to override PaasHost
	PaasHostEnvVar = "PAAS_HOST"
	// AuthHostVersionEnvVar is the env variable used to override AuthHostVersion
	AuthHostVersionEnvVar = "AUTH_HOST_VERSION"
	// PaasHostVersionEnvVar is the env variable used to override PaasHostVersion
	PaasHostVersionEnvVar = "PAAS_HOST_VERSION"
	// CatalyzeUsernameEnvVar is the env variable used to override the username
	CatalyzeUsernameEnvVar = "CATALYZE_USERNAME"
	// CatalyzePasswordEnvVar is the env variable used to override the passowrd
	CatalyzePasswordEnvVar = "CATALYZE_PASSWORD"
	// CatalyzeEnvironmentEnvVar is the env variable used to override the environment used in the current command
	CatalyzeEnvironmentEnvVar = "CATALYZE_ENV"
)
