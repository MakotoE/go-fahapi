package fahapi

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
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
	o.CaptureOnError = isTrue(m["capture-on-error"])
	o.CapturePackets = isTrue(m["capture-packets"])
	o.CaptureRequests = isTrue(m["capture-requests"])
	o.CaptureResponses = isTrue(m["capture-responses"])
	o.CaptureSockets = isTrue(m["capture-sockets"])
	o.Cause = m["cause"]
	o.CertificateFile = m["certificate-file"]
	o.Checkpoint, err = strconv.Atoi(m["checkpoint"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.Child = isTrue(m["child"])
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
	o.CommandEnable = isTrue(m["command-enable"])
	o.CommandPort, err = strconv.Atoi(m["command-port"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.ConfigRotate = isTrue(m["config-rotate"])
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
	o.Daemon = isTrue(m["daemon"])
	o.DebugSockets = isTrue(m["debug-sockets"])
	o.DisableSleepWhenActive = isTrue(m["disable-sleep-when-active"])
	o.DisableViz = isTrue(m["disable-viz"])
	o.DumpAfterDeadline = isTrue(m["dump-after-deadline"])
	o.ExceptionLocations = isTrue(m["exception-locations"])
	o.ExitWhenDone = isTrue(m["exit-when-done"])
	o.ExtraCoreArgs = m["extra-core-args"]
	o.FoldAnon = isTrue(m["fold-anon"])
	o.Gpu = isTrue(m["gpu"])
	o.GpuIndex = m["gpu-index"]
	o.GpuUsage, err = strconv.Atoi(m["gpu-usage"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.GuiEnabled = isTrue(m["gui-enabled"])
	o.HttpAddresses = m["http-addresses"]
	o.HttpsAddresses = m["https-addresses"]
	o.Idle = isTrue(m["idle"])
	o.Log = m["log"]
	o.LogColor = isTrue(m["log-color"])
	o.LogCrlf = isTrue(m["log-crlf"])
	o.LogDate = isTrue(m["log-date"])
	o.LogDatePeriodically, err = strconv.Atoi(m["log-date-periodically"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.LogDomain = isTrue(m["log-domain"])
	o.LogDomainLevels = m["log-domain-levels"]
	o.LogHeader = isTrue(m["log-header"])
	o.LogLevel = isTrue(m["log-level"])
	o.LogNoInfoHeader = isTrue(m["log-no-info-header"])
	o.LogRedirect = isTrue(m["log-redirect"])
	o.LogRotate = isTrue(m["log-rotate"])
	o.LogRotateDir = m["log-rotate-dir"]
	o.LogRotateMax, err = strconv.Atoi(m["log-rotate-max"])
	if err != nil {
		return errors.WithStack(err)
	}
	o.LogShortLevel = isTrue(m["log-short-level"])
	o.LogSimpleDomains = isTrue(m["log-simple-domains"])
	o.LogThreadId = isTrue(m["log-thread-id"])
	o.LogThreadPrefix = isTrue(m["log-thread-prefix"])
	o.LogTime = isTrue(m["log-time"])
	o.LogToScreen = isTrue(m["log-to-screen"])
	o.LogTruncate = isTrue(m["log-truncate"])
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
	o.NoAssembly = isTrue(m["no-assembly"])
	o.OpenWebControl = isTrue(m["open-web-control"])
	o.OpenclIndex = m["opencl-index"]
	o.OsSpecies = m["os-species"]
	o.OsType = m["os-type"]
	o.Passkey = m["passkey"]
	o.Password = m["password"]
	o.PauseOnBattery = isTrue(m["pause-on-battery"])
	o.PauseOnStart = isTrue(m["pause-on-start"])
	o.Paused = isTrue(m["paused"])
	o.Pid = isTrue(m["pid"])
	o.PidFile = isTrue(m["pid-file"])
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
	o.ProxyEnable = isTrue(m["proxy-enable"])
	o.ProxyPass = m["proxy-pass"]
	o.ProxyUser = m["proxy-user"]
	o.Respawn = isTrue(m["respawn"])
	o.Service = isTrue(m["service"])
	o.ServiceDescription = m["service-description"]
	o.ServiceRestart = isTrue(m["service-restart"])
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
	o.Smp = isTrue(m["smp"])
	o.StackTraces = isTrue(m["stack-traces"])
	o.StallDetectionEnabled = isTrue(m["stall-detection-enabled"])
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
	o.WebEnable = isTrue(m["web-enable"])
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

func isTrue(s string) bool {
	return s == "true"
}

type slotQueueInfoRaw struct {
	ID             string `json:"id"`
	State          string `json:"state"`
	Error          string `json:"error"`
	Project        int    `json:"project"`
	Run            int    `json:"run"`
	Clone          int    `json:"clone"`
	Gen            int    `json:"gen"`
	Core           string `json:"core"`
	Unit           string `json:"unit"`
	PercentDone    string `json:"percentdone"`
	ETA            string `json:"eta"`
	PPD            string `json:"ppd"`
	CreditEstimate string `json:"creditestimate"`
	WaitingOn      string `json:"waitingon"`
	NextAttempt    string `json:"nextattempt"`
	TimeRemaining  string `json:"timeremaining"`
	TotalFrames    int    `json:"totalframes"`
	FramesDone     int    `json:"framesdone"`
	Assigned       string `json:"assigned"`
	Timeout        string `json:"timeout"`
	Deadline       string `json:"deadline"`
	WS             string `json:"ws"`
	CS             string `json:"cs"`
	Attempts       int    `json:"attempts"`
	Slot           string `json:"slot"`
	TPF            string `json:"tpf"`
	BaseCredit     string `json:"basecredit"`
}

type SlotQueueInfo struct {
	ID             string
	State          string // TODO create enum type (need official documentation to know every possible value)
	Error          string // This should also be enum
	Project        int
	Run            int
	Clone          int
	Gen            int
	Core           string
	Unit           string
	PercentDone    string
	ETA            time.Duration
	PPD            int
	CreditEstimate int
	WaitingOn      string
	NextAttempt    time.Duration
	TimeRemaining  time.Duration
	TotalFrames    int
	FramesDone     int
	Assigned       time.Time
	Timeout        time.Time
	Deadline       time.Time
	WS             string // This could be an IP address but I don't know for sure
	CS             string
	Attempts       int
	Slot           string
	TPF            time.Duration
	BaseCredit     int
}

func (s *SlotQueueInfo) fromRaw(r *slotQueueInfoRaw) error {
	var err error
	s.ID = r.ID
	s.State = r.State
	s.Error = r.Error
	s.Project = r.Project
	s.Run = r.Run
	s.Clone = r.Clone
	s.Gen = r.Gen
	s.Core = r.Core
	s.Unit = r.Unit
	s.PercentDone = r.PercentDone
	s.ETA, err = parseFAHDuration(r.ETA)
	if err != nil {
		return err
	}
	s.PPD, err = strconv.Atoi(r.PPD)
	if err != nil {
		return errors.WithStack(err)
	}
	s.CreditEstimate, err = strconv.Atoi(r.CreditEstimate)
	if err != nil {
		return errors.WithStack(err)
	}
	s.WaitingOn = r.WaitingOn
	s.NextAttempt, err = parseFAHDuration(r.NextAttempt)
	if err != nil {
		return err
	}
	s.TotalFrames = r.TotalFrames
	s.FramesDone = r.FramesDone
	s.Assigned, err = time.Parse(time.RFC3339, r.Assigned)
	if err != nil {
		return errors.WithStack(err)
	}
	s.Timeout, err = time.Parse(time.RFC3339, r.Timeout)
	if err != nil {
		return errors.WithStack(err)
	}
	s.Deadline, err = time.Parse(time.RFC3339, r.Deadline)
	if err != nil {
		return errors.WithStack(err)
	}
	s.WS = r.WS
	s.CS = r.CS
	s.Attempts = r.Attempts
	s.Slot = r.Slot
	s.TPF, err = parseFAHDuration(r.TPF)
	if err != nil {
		return err
	}
	s.BaseCredit, err = strconv.Atoi(r.BaseCredit)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

var parseFAHDurationReplacer = strings.NewReplacer(
	" ", "",
	"days", "d",
	"day", "d",
	"hours", "h",
	"hour", "h",
	"mins", "m",
	"min", "m",
	"secs", "s",
	"sec", "s",
)

func parseFAHDuration(s string) (time.Duration, error) {
	shortened := parseFAHDurationReplacer.Replace(s)
	dIndex := strings.IndexByte(shortened, 'd')
	days := 0.0
	if dIndex > -1 {
		daysTemp, err := strconv.ParseFloat(shortened[:dIndex], 64)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		days = daysTemp

		if dIndex >= len(shortened)-1 { // s only contains days
			return time.Hour * 24 * time.Duration(days), nil
		}
	}

	duration, err := time.ParseDuration(shortened[dIndex+1:])
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return duration + time.Hour*24*time.Duration(days), nil
}
