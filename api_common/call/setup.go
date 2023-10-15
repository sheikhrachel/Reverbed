package call

// Call is a utility struct used for simplifying logging and metrics throughout the codebase
type Call struct {
	// Env is the environment the app is running in
	// - ex. cc.Env returns "dev"
	Env string
	// Region is the AWS region the app is running in
	// - ex. cc.Region returns "us-east-1"
	Region string
}

// New returns a new Call pointer, and initialises the statsd client for Datadog metrics on port 8125
func New(appEnv, appRegion string) Call {
	cc := Call{Env: appEnv, Region: appRegion}
	return cc
}
