package utils

const (
	// ConfigurationFilePath is the constant path to the configuration file needed to start the application
	// written from root file since the application starts from `make run`
	ConfigurationFilePath = "local-config.yml"
	// PathPing stores the defualt address of storage directory of ping data
	PathPing = "storage/ping"
	// PathJitter stores the defualt address of storage directory of jitter data
	PathJitter = "storage/jitter"
	// PathFloodPing stores the defualt address of storage directory of flood ping data
	PathFloodPing = "storage/flood-ping"
	// PathReqResDelayMonitoring stores the defualt address of storage directory of req-res and monitoring data
	PathReqResDelayMonitoring = "storage/req-res-delay-monitoring"
)

// TypePingScrap as datatype for ping outputs
type TypePingScrap struct {
	Min, Avg, Max, Mdev float64
}

// TypeFloodPingScrap as datatype for flood ping outputs
type TypeFloodPingScrap struct {
	Min, Avg, Max, Mdev, PacketLoss float64
}

// Response struct
// This is the object that we return from resp_delay module
// Contains delay in response and the response length
type Response struct {
	Delay         int
	ResLength     int64
	ResStatusCode int
}

// YAMLBenchRoutesType defines the structure type for implementing the interface
type YAMLBenchRoutesType struct {
	Address string
	Config  *ConfigurationBR
}

// Interval sets a type for intervals between consecutive similar tests
type Interval struct {
	Test     string `yaml:"test"`
	Type     string `yaml:"type"`
	Duration int64  `yaml:"duration"`
}

// Headers store the header values(ofType and value) from the config file
type Headers struct {
	OfType string `yaml:"type"`
	Value  string `yaml:"value"`
}

// Params type for parameters passed along the url for specific route
type Params struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Routes sets routes mentioned in configuration file
type Routes struct {
	Method string    `yaml:"method"`
	URL    string    `yaml:"url"`
	Route  string    `yaml:"route"`
	Header []Headers `yaml:"headers"`
	Params []Params  `yaml:"params"`
}

// ResponseChangesConfig acts as a type for response-length configuration in config.yml
type ResponseChangesConfig struct {
	Mode float32 `yaml:"mode"`
	Mean float32 `yaml:"mean"`
}

// ServiceSignals type for defining current running states of various services supported
// by BR. Allowed only two values: `active` OR `passive`
type ServiceSignals struct {
	Ping                  string `yaml:"ping"`
	FloodPing             string `yaml:"flood-ping"`
	Jitter                string `yaml:"jitter"`
	ReqResDelayMonitoring string `yaml:"req-res-delay-or-monitoring"`
}

// UConfig type for storing utilities in config.yml as local DB
type UConfig struct {
	RespChanges    ResponseChangesConfig `yaml:"response-length"`
	ServicesSignal ServiceSignals        `yaml:"services-state"`
}

// ConfigurationBR sets a type for configuration file which also acts as a local DB
type ConfigurationBR struct {
	Password  string     `yaml:"password"`
	Interval  []Interval `yaml:"test_interval"`
	Routes    []Routes   `yaml:"routes"`
	UtilsConf UConfig    `yaml:"utils"`
}