package models

// Errors is a wrapper around an array of errors from the API
type Errors struct {
	ReportedErrors *[]ReportedError `json:"errors"`
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	Code           int              `json:"code"`
}

// ReportedError is the standard error model sent back from the API
type ReportedError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Login is used for making an authentication request
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Invite represents an invitation to an environment
type Invite struct {
	ID            string `json:"id,omitempty"`
	Code          string `json:"code,omitempty"`
	Email         string `json:"email"`
	EnvironmentID string `json:"environmentId,omitempty"`
}

// User is an authenticated User
type User struct {
	Username     string `json:"username"`
	SessionToken string `json:"sessionToken"`
	UsersID      string `json:"usersId"`
}

// EncryptionStore holds the values for encryption on backup/import jobs
type EncryptionStore struct {
	Key string `json:"key"`
	IV  string `json:"iv"`
}

// Import creation object
type Import struct {
	Location  string            `json:"location"`
	Key       string            `json:"key"`
	IV        string            `json:"iv"`
	WipeFirst bool              `json:"wipeBeforeImport"`
	Options   map[string]string `json:"options"`
}

// TempURL holds a URL for uploading or downloading files from a temporary URL
type TempURL struct {
	URL string `json:"url"`
}

// Task task
type Task struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// EnvironmentData is the data blob inside of an Environment
type EnvironmentData struct {
	Name      string     `json:"name"`
	Services  *[]Service `json:"services"`
	Namespace string     `json:"namespace"`
	DNSName   string     `json:"dns_name"`
}

// EnvironmentUsers users who have access to an environment
type EnvironmentUsers struct {
	Users []string `json:"users"`
}

// Environment environment
type Environment struct {
	ID    string           `json:"environmentId"`
	State string           `json:"state"`
	Data  *EnvironmentData `json:"data"`
	PodID string           `json:"podId"`
	Name  string           `json:"name"`
}

// Job job
type Job struct {
	ID          string           `json:"id"`
	Type        string           `json:"type"`
	Status      string           `json:"status,omitempty"`
	Backup      *EncryptionStore `json:"backup,omitempty"`
	Restore     *EncryptionStore `json:"restore,omitempty"`
	CreatedAt   string           `json:"created_at"`
	MetricsData *[]MetricsData   `json:"metrics"`
}

// Service service
type Service struct {
	ID           string            `json:"id"`
	Type         string            `json:"type"`
	Label        string            `json:"label"`
	Size         interface{}       `json:"size"`
	BuildStatus  string            `json:"build_status"`
	DeployStatus string            `json:"deploy_status"`
	Name         string            `json:"name"`
	EnvVars      map[string]string `json:"environmentVariables"`
	Source       string            `json:"source"`
	LBIP         string            `json:"load_balancer_ip,omitempty"`
}

// ServiceSize holds size information for a service
type ServiceSize struct {
	ServiceID string `json:"service"`
	RAM       string `json:"ram"`
	Storage   string `json:"storage"`
	Behavior  string `json:"behavior"`
	Type      string `json:"type"`
	CPU       string `json:"cpu"`
}

// PodMetadata podmetadata
type PodMetadata struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	PHISafe              bool   `json:"phiSafe"`
	ImportRequiresLength bool   `json:"importRequiresLength"`
}

// Settings holds various settings for the current context. All items with
// `json:"-"` are never persisted to disk but used in memory for the current
// command.
type Settings struct {
	BaasHost        string                   `json:"-"`
	PaasHost        string                   `json:"-"`
	Username        string                   `json:"-"`
	Password        string                   `json:"-"`
	EnvironmentID   string                   `json:"-"` // the id of the environment used for the current command
	ServiceID       string                   `json:"-"` // the id of the service used for the current command
	EnvironmentName string                   `json:"-"` // the name of the environment used for the current command
	SessionToken    string                   `json:"token"`
	UsersID         string                   `json:"user_id"`
	Environments    map[string]AssociatedEnv `json:"environments"`
	Default         string                   `json:"default"`
}

// AssociatedEnv holds information about an associated environment
type AssociatedEnv struct {
	EnvironmentID string `json:"environmentId"`
	ServiceID     string `json:"serviceId"`
	Directory     string `json:"dir"`
	Name          string `json:"name"`
}

// Breadcrumb is stored in a local git repo to make a link back to the
// global list of associated envs
type Breadcrumb struct {
	EnvName       string `json:"env_name"`
	EnvironmentID string `json:"environmentId,omitempty"` // for backwards compatibility
	ServiceID     string `json:"serviceId,omitempty"`     // for backwards compatibility
	SessionToken  string `json:"token,omitempty"`         // for backwards compatibility
	UsersID       string `json:"user_id,omitempty"`       // for backwards compatibility
}

// ConsoleCredentials hold the keys necessary for connecting to a console service
type ConsoleCredentials struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

// Metrics holds all metrics data for an entire environment or a single service
type Metrics struct {
	ServiceName  string `json:"serviceName"`
	ServiceType  string `json:"serviceType"`
	ServiceID    string `json:"serviceId"`
	ServiceLabel string `json:"serviceLabel"`
	Jobs         *[]Job `json:"jobs"`
}

// MetricsData is a container for each type of metrics: network, memory, etc.
type MetricsData struct {
	Network *NetworkData `json:"network"`
	Memory  *MinMaxAvg   `json:"memory"`
	DiskIO  *DiskIOData  `json:"diskio"`
	TS      int64        `json:"ts"`
	Name    string       `json:"name"`
	CPU     *CPUData     `json:"cpu"`
}

// NetworkData holds metrics data for the network category
type NetworkData struct {
	TXKb      float64 `json:"tx_kb"`
	TXPackets float64 `json:"tx_packets"`
	TXDropped float64 `json:"tx_dropped"`
	TXErrors  float64 `json:"tx_errors"`
	RXKb      float64 `json:"rx_kb"`
	RXPackets float64 `json:"rx_packets"`
	RXDropped float64 `json:"rx_dropped"`
	RXErrors  float64 `json:"rx_errors"`
}

// MinMaxAvg is a generic metrics data structure holding the minimum, maximum,
// and average values.
type MinMaxAvg struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
	Avg float64 `json:"ave"`
}

// DiskIOData is a data structure holding metrics values for the DiskIO category
type DiskIOData struct {
	Read  float64 `json:"read"`
	Async float64 `json:"async"`
	Write float64 `json:"write"`
	Sync  float64 `json:"sync"`
}

// CPUData is a data structure holding metrics values for the CPU category
type CPUData struct {
	Usage float64    `json:"usage"`
	Load  *MinMaxAvg `json:"load"`
}

// LogQuery holds data for querying Kibana's elastic search
type LogQuery struct {
	Fields []string     `json:"fields"`
	Query  *Query       `json:"query"`
	Filter *FilterRange `json:"filter"`
	Sort   *LogSort     `json:"sort"`
	From   int          `json:"from"`
	Size   int          `json:"size"`
}

// LogSort tells the LogQuery how to sort the results
type LogSort struct {
	Timestamp map[string]string `json:"@timestamp"`
	Message   map[string]string `json:"message"`
}

// Query holds query terms
type Query struct {
	Wildcard map[string]string `json:"wildcard"`
}

// FilterRange holds query filter range
type FilterRange struct {
	Range *RangeTimestamp `json:"range"`
}

// RangeTimestamp holds range timestamps
type RangeTimestamp struct {
	Timestamp map[string]string `json:"@timestamp"`
}

// Logs hold the log values from a successful LogQuery
type Logs struct {
	Hits *Hits `json:"hits"`
}

// Hits contain arrays of log data
type Hits struct {
	Total    int64      `json:"total"`
	MaxScore float64    `json:"max_score"`
	Hits     *[]LogHits `json:"hits"`
}

// LogHits contain ordering data for logs
type LogHits struct {
	Index  string              `json:"_index"`
	Type   string              `json:"_type"`
	ID     string              `json:"_id"`
	Score  float64             `json:"_score"`
	Fields map[string][]string `json:"fields"`
}
