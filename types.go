package go_fahapi

import (
	"github.com/pkg/errors"
	"strconv"
)

type Options struct {
	Allow                  string
	CaptureDirectory       string
	CaptureOnError         bool
	CapturePackets         bool
	CaptureRequests        bool
	CaptureResponses       bool
	CaptureSockets         bool
	Cause                  string
	CertificateFile        string
	Checkpoint             int
	Child                  bool
	ClientSubtype          string
	ClientThreads          int
	ClientType             string
	CommandAddress         string
	CommandAllowNoPass     string
	Deny                   string
	CommandDenyNoPass      string
	CommandEnable          bool
	CommandPort            int
	ConfigRotate           bool
	ConfigRotateDir        string
	ConfigRotateMax        int
	ConnectionTimeout      int
	CorePriority           string
	CpuSpecies             string
	CpuType                string
	CpuUsage               int
	Cpus                   int
	CrlFile                string
	CudaIndex              string
	CycleRate              int
	Cycles                 int
	Daemon                 bool
	DebugSockets           bool
	DisableSleepWhenActive bool
	DisableViz             bool
	DumpAfterDeadline      bool
	ExceptionLocations     bool
	ExitWhenDone           bool
	ExtraCoreArgs          string
	FoldAnon               bool
	Gpu                    bool
	GpuIndex               string
	GpuUsage               int
	GuiEnabled             bool
	HttpAddresses          string
	HttpsAddresses         string
	Idle                   bool
	Log                    string
	LogColor               bool
	LogCrlf                bool
	LogDate                bool
	LogDatePeriodically    int
	LogDomain              bool
	LogDomainLevels        string
	LogHeader              bool
	LogLevel               bool
	LogNoInfoHeader        bool
	LogRedirect            bool
	LogRotate              bool
	LogRotateDir           string
	LogRotateMax           int
	LogShortLevel          bool
	LogSimpleDomains       bool
	LogThreadId            bool
	LogThreadPrefix        bool
	LogTime                bool
	LogToScreen            bool
	LogTruncate            bool
	MachineId              int
	MaxConnectTime         int
	MaxConnections         int
	MaxPacketSize          string
	MaxQueue               int
	MaxRequestLength       int
	MaxShutdownWait        int
	MaxSlotErrors          int
	MaxUnitErrors          int
	MaxUnits               int
	Memory                 string
	MinConnectTime         int
	NextUnitPercentage     int
	Priority               string
	NoAssembly             bool
	OpenWebControl         bool
	OpenclIndex            string
	OsSpecies              string
	OsType                 string
	Passkey                string
	Password               string
	PauseOnBattery         bool
	PauseOnStart           bool
	Paused                 bool
	Pid                    bool
	PidFile                bool
	Power                  Power
	PrivateKeyFile         string
	ProjectKey             int
	Proxy                  string
	ProxyEnable            bool
	ProxyPass              string
	ProxyUser              string
	Respawn                bool
	Service                bool
	ServiceDescription     string
	ServiceRestart         bool
	ServiceRestartDelay    int
	SessionCookie          string
	SessionLifetime        int
	SessionTimeout         int
	Smp                    bool
	StackTraces            bool
	StallDetectionEnabled  bool
	StallPercent           int
	StallTimeout           int
	Team                   int
	User                   string
	Verbosity              int
	WebAllow               string
	WebDeny                string
	WebEnable              bool
}

func (o *Options) fromMap(m map[string]string) error {
	var err error
	o.Allow = m["allow"]
	o.CaptureDirectory = m["capture-directory"]
	o.CaptureOnError = IsTrue(m["capture-on-error"])
	o.CapturePackets = IsTrue(m["capture-packets"])
	o.CaptureRequests = IsTrue(m["capture-requests"])
	o.CaptureResponses = IsTrue(m["capture-responses"])
	o.CaptureSockets = IsTrue(m["capture-sockets"])
	o.Cause = m["cause"]
	o.CertificateFile = m["certificate-file"]
	o.Checkpoint, err = strconv.Atoi(m["checkpoint"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.Child = IsTrue(m["child"])
	o.ClientSubtype = m["client-subtype"]
	o.ClientThreads, err = strconv.Atoi(m["client-threads"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.ClientType = m["client-type"]
	o.CommandAddress = m["command-address"]
	o.CommandAllowNoPass = m["command-allow-no-pass"]
	o.Deny = m["deny"]
	o.CommandDenyNoPass = m["command-deny-no-pass"]
	o.CommandEnable = IsTrue(m["command-enable"])
	o.CommandPort, err = strconv.Atoi(m["command-port"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.ConfigRotate = IsTrue(m["config-rotate"])
	o.ConfigRotateDir = m["config-rotate-dir"]
	o.ConfigRotateMax, err = strconv.Atoi(m["config-rotate-max"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.ConnectionTimeout, err = strconv.Atoi(m["connection-timeout"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.CorePriority = m["core-priority"]
	o.CpuSpecies = m["cpu-species"]
	o.CpuType = m["cpu-type"]
	o.CpuUsage, err = strconv.Atoi(m["cpu-usage"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.Cpus, err = strconv.Atoi(m["cpus"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.CrlFile = m["crl-file"]
	o.CudaIndex = m["cuda-index"]
	o.CycleRate, err = strconv.Atoi(m["cycle-rate"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.Cycles, err = strconv.Atoi(m["cycles"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.Daemon = IsTrue(m["daemon"])
	o.DebugSockets = IsTrue(m["debug-sockets"])
	o.DisableSleepWhenActive = IsTrue(m["disable-sleep-when-active"])
	o.DisableViz = IsTrue(m["disable-viz"])
	o.DumpAfterDeadline = IsTrue(m["dump-after-deadline"])
	o.ExceptionLocations = IsTrue(m["exception-locations"])
	o.ExitWhenDone = IsTrue(m["exit-when-done"])
	o.ExtraCoreArgs = m["extra-core-args"]
	o.FoldAnon = IsTrue(m["fold-anon"])
	o.Gpu = IsTrue(m["gpu"])
	o.GpuIndex = m["gpu-index"]
	o.GpuUsage, err = strconv.Atoi(m["gpu-usage"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.GuiEnabled = IsTrue(m["gui-enabled"])
	o.HttpAddresses = m["http-addresses"]
	o.HttpsAddresses = m["https-addresses"]
	o.Idle = IsTrue(m["idle"])
	o.Log = m["log"]
	o.LogColor = IsTrue(m["log-color"])
	o.LogCrlf = IsTrue(m["log-crlf"])
	o.LogDate = IsTrue(m["log-date"])
	o.LogDatePeriodically, err = strconv.Atoi(m["log-date-periodically"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.LogDomain = IsTrue(m["log-domain"])
	o.LogDomainLevels = m["log-domain-levels"]
	o.LogHeader = IsTrue(m["log-header"])
	o.LogLevel = IsTrue(m["log-level"])
	o.LogNoInfoHeader = IsTrue(m["log-no-info-header"])
	o.LogRedirect = IsTrue(m["log-redirect"])
	o.LogRotate = IsTrue(m["log-rotate"])
	o.LogRotateDir = m["log-rotate-dir"]
	o.LogRotateMax, err = strconv.Atoi(m["log-rotate-max"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.LogShortLevel = IsTrue(m["log-short-level"])
	o.LogSimpleDomains = IsTrue(m["log-simple-domains"])
	o.LogThreadId = IsTrue(m["log-thread-id"])
	o.LogThreadPrefix = IsTrue(m["log-thread-prefix"])
	o.LogTime = IsTrue(m["log-time"])
	o.LogToScreen = IsTrue(m["log-to-screen"])
	o.LogTruncate = IsTrue(m["log-truncate"])
	o.MachineId, err = strconv.Atoi(m["machine-id"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.MaxConnectTime, err = strconv.Atoi(m["max-connect-time"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.MaxConnections, err = strconv.Atoi(m["max-connections"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.MaxPacketSize = m["max-packet-size"]
	o.MaxQueue, err = strconv.Atoi(m["max-queue"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.MaxRequestLength, err = strconv.Atoi(m["max-request-length"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.MaxShutdownWait, err = strconv.Atoi(m["max-shutdown-wait"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.MaxSlotErrors, err = strconv.Atoi(m["max-slot-errors"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.MaxUnitErrors, err = strconv.Atoi(m["max-unit-errors"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.MaxUnits, err = strconv.Atoi(m["max-units"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.Memory = m["memory"]
	o.MinConnectTime, err = strconv.Atoi(m["min-connect-time"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.NextUnitPercentage, err = strconv.Atoi(m["next-unit-percentage"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.Priority = m["priority"]
	o.NoAssembly = IsTrue(m["no-assembly"])
	o.OpenWebControl = IsTrue(m["open-web-control"])
	o.OpenclIndex = m["opencl-index"]
	o.OsSpecies = m["os-species"]
	o.OsType = m["os-type"]
	o.Passkey = m["passkey"]
	o.Password = m["password"]
	o.PauseOnBattery = IsTrue(m["pause-on-battery"])
	o.PauseOnStart = IsTrue(m["pause-on-start"])
	o.Paused = IsTrue(m["paused"])
	o.Pid = IsTrue(m["pid"])
	o.PidFile = IsTrue(m["pid-file"])
	o.Power, err = NewPower(m["power"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.PrivateKeyFile = m["private-key-file"]
	o.ProjectKey, err = strconv.Atoi(m["project-key"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.Proxy = m["proxy"]
	o.ProxyEnable = IsTrue(m["proxy-enable"])
	o.ProxyPass = m["proxy-pass"]
	o.ProxyUser = m["proxy-user"]
	o.Respawn = IsTrue(m["respawn"])
	o.Service = IsTrue(m["service"])
	o.ServiceDescription = m["service-description"]
	o.ServiceRestart = IsTrue(m["service-restart"])
	o.ServiceRestartDelay, err = strconv.Atoi(m["service-restart-delay"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.SessionCookie = m["session-cookie"]
	o.SessionLifetime, err = strconv.Atoi(m["session-lifetime"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.SessionTimeout, err = strconv.Atoi(m["session-timeout"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.Smp = IsTrue(m["smp"])
	o.StackTraces = IsTrue(m["stack-traces"])
	o.StallDetectionEnabled = IsTrue(m["stall-detection-enabled"])
	o.StallPercent, err = strconv.Atoi(m["stall-percent"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.StallTimeout, err = strconv.Atoi(m["stall-timeout"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.Team, err = strconv.Atoi(m["team"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.User = m["user"]
	o.Verbosity, err = strconv.Atoi(m["verbosity"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.WebAllow = m["web-allow"]
	o.WebDeny = m["web-deny"]
	o.WebEnable = IsTrue(m["web-enable"])
	return nil
}

type Power string

const (
	PowerLight  Power = "LIGHT"
	PowerMedium Power = "MEDIUM"
	PowerFull   Power = "FULL"
)

func NewPower(s string) (Power, error) {
	if s == string(PowerLight) || s == string(PowerMedium) || s == string(PowerFull) {
		return Power(s), nil
	}

	return "", errors.Errorf("s is invalid: %s", s)
}

func IsTrue(s string) bool {
	return s == "true"
}
