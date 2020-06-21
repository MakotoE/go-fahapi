package fahapi

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	Allow                  string     `json:"allow"`
	CaptureDirectory       string     `json:"capture-directory"`
	CaptureOnError         StringBool `json:"capture-on-error"`
	CapturePackets         StringBool `json:"capture-packets"`
	CaptureRequests        StringBool `json:"capture-requests"`
	CaptureResponses       StringBool `json:"capture-responses"`
	CaptureSockets         StringBool `json:"capture-sockets"`
	Cause                  string     `json:"cause"`
	CertificateFile        string     `json:"certificate-file"`
	Checkpoint             StringInt  `json:"checkpoint"`
	Child                  StringBool `json:"child"`
	ClientSubtype          string     `json:"client-subtype"`
	ClientThreads          StringInt  `json:"client-threads"`
	ClientType             string     `json:"client-type"`
	CommandAddress         string     `json:"command-address"`
	CommandAllowNoPass     string     `json:"command-allow-no-pass"`
	Deny                   string     `json:"deny"`
	CommandDenyNoPass      string     `json:"command-deny-no-pass"`
	CommandEnable          StringBool `json:"command-enable"`
	CommandPort            StringInt  `json:"command-port"`
	ConfigRotate           StringBool `json:"config-rotate"`
	ConfigRotateDir        string     `json:"config-rotate-dir"`
	ConfigRotateMax        StringInt  `json:"config-rotate-max"`
	ConnectionTimeout      StringInt  `json:"connection-timeout"`
	CorePriority           string     `json:"core-priority"`
	CpuSpecies             string     `json:"cpu-species"`
	CpuType                string     `json:"cpu-type"`
	CpuUsage               StringInt  `json:"cpu-usage"`
	Cpus                   StringInt  `json:"cpus"`
	CrlFile                string     `json:"crl-file"`
	CudaIndex              string     `json:"cuda-index"`
	CycleRate              StringInt  `json:"cycle-rate"`
	Cycles                 StringInt  `json:"cycles"`
	Daemon                 StringBool `json:"daemon"`
	DebugSockets           StringBool `json:"debug-sockets"`
	DisableSleepWhenActive StringBool `json:"disable-sleep-when-active"`
	DisableViz             StringBool `json:"disable-viz"`
	DumpAfterDeadline      StringBool `json:"dump-after-deadline"`
	ExceptionLocations     StringBool `json:"exception-locations"`
	ExitWhenDone           StringBool `json:"exit-when-done"`
	ExtraCoreArgs          string     `json:"extra-core-args"`
	FoldAnon               StringBool `json:"fold-anon"`
	Gpu                    StringBool `json:"gpu"`
	GpuIndex               string     `json:"gpu-index"`
	GpuUsage               StringInt  `json:"gpu-usage"`
	GuiEnabled             StringBool `json:"gui-enabled"`
	HttpAddresses          string     `json:"http-addresses"`
	HttpsAddresses         string     `json:"https-addresses"`
	Idle                   StringBool `json:"idle"`
	Log                    string     `json:"log"`
	LogColor               StringBool `json:"log-color"`
	LogCrlf                StringBool `json:"log-crlf"`
	LogDate                StringBool `json:"log-date"`
	LogDatePeriodically    StringInt  `json:"log-date-periodically"`
	LogDomain              StringBool `json:"log-domain"`
	LogDomainLevels        string     `json:"log-domain-levels"`
	LogHeader              StringBool `json:"log-header"`
	LogLevel               StringBool `json:"log-level"`
	LogNoInfoHeader        StringBool `json:"log-no-info-header"`
	LogRedirect            StringBool `json:"log-redirect"`
	LogRotate              StringBool `json:"log-rotate"`
	LogRotateDir           string     `json:"log-rotate-dir"`
	LogRotateMax           StringInt  `json:"log-rotate-max"`
	LogShortLevel          StringBool `json:"log-short-level"`
	LogSimpleDomains       StringBool `json:"log-simple-domains"`
	LogThreadId            StringBool `json:"log-thread-id"`
	LogThreadPrefix        StringBool `json:"log-thread-prefix"`
	LogTime                StringBool `json:"log-time"`
	LogToScreen            StringBool `json:"log-to-screen"`
	LogTruncate            StringBool `json:"log-truncate"`
	MachineId              StringInt  `json:"machine-id"`
	MaxConnectTime         StringInt  `json:"max-connect-time"`
	MaxConnections         StringInt  `json:"max-connections"`
	MaxPacketSize          string     `json:"max-packet-size"`
	MaxQueue               StringInt  `json:"max-queue"`
	MaxRequestLength       StringInt  `json:"max-request-length"`
	MaxShutdownWait        StringInt  `json:"max-shutdown-wait"`
	MaxSlotErrors          StringInt  `json:"max-slot-errors"`
	MaxUnitErrors          StringInt  `json:"max-unit-errors"`
	MaxUnits               StringInt  `json:"max-units"`
	Memory                 string     `json:"memory"`
	MinConnectTime         StringInt  `json:"min-connect-time"`
	NextUnitPercentage     StringInt  `json:"next-unit-percentage"`
	Priority               string     `json:"priority"`
	NoAssembly             StringBool `json:"no-assembly"`
	OpenWebControl         StringBool `json:"open-web-control"`
	OpenclIndex            string     `json:"opencl-index"`
	OsSpecies              string     `json:"os-species"`
	OsType                 string     `json:"os-type"`
	Passkey                string     `json:"passkey"`
	Password               string     `json:"password"`
	PauseOnBattery         StringBool `json:"pause-on-battery"`
	PauseOnStart           StringBool `json:"pause-on-start"`
	Paused                 StringBool `json:"paused"`
	Pid                    StringBool `json:"pid"`
	PidFile                string     `json:"pid-file"`
	Power                  Power      `json:"power"`
	PrivateKeyFile         string     `json:"private-key-file"`
	ProjectKey             StringInt  `json:"project-key"`
	Proxy                  string     `json:"proxy"`
	ProxyEnable            StringBool `json:"proxy-enable"`
	ProxyPass              string     `json:"proxy-pass"`
	ProxyUser              string     `json:"proxy-user"`
	Respawn                StringBool `json:"respawn"`
	Service                StringBool `json:"service"`
	ServiceDescription     string     `json:"service-description"`
	ServiceRestart         StringBool `json:"service-restart"`
	ServiceRestartDelay    StringInt  `json:"service-restart-delay"`
	SessionCookie          string     `json:"session-cookie"`
	SessionLifetime        StringInt  `json:"session-lifetime"`
	SessionTimeout         StringInt  `json:"session-timeout"`
	Smp                    StringBool `json:"smp"`
	StackTraces            StringBool `json:"stack-traces"`
	StallDetectionEnabled  StringBool `json:"stall-detection-enabled"`
	StallPercent           StringInt  `json:"stall-percent"`
	StallTimeout           StringInt  `json:"stall-timeout"`
	Team                   StringInt  `json:"team"`
	User                   string     `json:"user"`
	Verbosity              StringInt  `json:"verbosity"`
	WebAllow               string     `json:"web-allow"`
	WebDeny                string     `json:"web-deny"`
	WebEnable              StringBool `json:"web-enable"`
}

type StringBool bool

var _ json.Unmarshaler = (*StringBool)(nil)

func (s *StringBool) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte(`"true"`)) {
		*s = true
		return nil
	} else if bytes.Equal(b, []byte(`"false"`)) {
		*s = false
		return nil
	}

	return &json.UnmarshalTypeError{
		Value: string(b),
		Type:  reflect.TypeOf(s),
	}
}

type StringInt int

var _ json.Unmarshaler = (*StringInt)(nil)

func (i *StringInt) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if err := i.FromString(s); err != nil {
		return &json.UnmarshalTypeError{
			Value: string(b),
			Type:  reflect.TypeOf(i),
		}
	}

	return nil
}

func (i *StringInt) FromString(s string) error {
	integer, err := strconv.Atoi(s)
	*i = StringInt(integer)
	return errors.WithStack(err)
}

type Power string

var _ json.Unmarshaler = (*Power)(nil)

const (
	PowerNull   Power = ""
	PowerLight  Power = "LIGHT"
	PowerMedium Power = "MEDIUM"
	PowerFull   Power = "FULL"
)

func NewPower(s string) (Power, error) {
	uppercased := Power(strings.ToUpper(s))
	switch uppercased {
	case PowerNull, PowerLight, PowerMedium, PowerFull:
		return uppercased, nil
	}

	return PowerNull, errors.Errorf("s is invalid: %s", s)
}

func (p *Power) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	power, err := NewPower(s)
	if err != nil {
		return &json.UnmarshalTypeError{
			Value: string(b),
			Type:  reflect.TypeOf(p),
		}
	}

	*p = power
	return nil
}

type SlotQueueInfo struct {
	ID             string      `json:"id"`
	State          string      `json:"state"`
	Error          string      `json:"error"`
	Project        int         `json:"project"`
	Run            int         `json:"run"`
	Clone          int         `json:"clone"`
	Gen            int         `json:"gen"`
	Core           string      `json:"core"`
	Unit           string      `json:"unit"`
	PercentDone    string      `json:"percentdone"`
	ETA            FAHDuration `json:"eta"`
	PPD            StringInt   `json:"ppd"`
	CreditEstimate StringInt   `json:"creditestimate"`
	WaitingOn      string      `json:"waitingon"`
	NextAttempt    FAHDuration `json:"nextattempt"`
	TimeRemaining  FAHDuration `json:"timeremaining"`
	TotalFrames    int         `json:"totalframes"`
	FramesDone     int         `json:"framesdone"`
	Assigned       FAHTime     `json:"assigned"`
	Timeout        FAHTime     `json:"timeout"`
	Deadline       FAHTime     `json:"deadline"`
	WS             string      `json:"ws"`
	CS             string      `json:"cs"`
	Attempts       int         `json:"attempts"`
	Slot           string      `json:"slot"`
	TPF            FAHDuration `json:"tpf"`
	BaseCredit     StringInt   `json:"basecredit"`
}

// FAHDuration may be "unknowntime", which can be checked by calling duration.UnknownTime().
type FAHDuration time.Duration

var _ json.Unmarshaler = (*FAHDuration)(nil)

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

const unknowntime = FAHDuration(-1)

const unknowntimeStr = "unknowntime"

func ParseFAHDuration(s string) (FAHDuration, error) {
	shortened := parseFAHDurationReplacer.Replace(s)
	if shortened == unknowntimeStr {
		return unknowntime, nil
	}

	dIndex := strings.IndexByte(shortened, 'd')
	days := 0.0
	if dIndex > -1 {
		daysTemp, err := strconv.ParseFloat(shortened[:dIndex], 64)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		days = daysTemp

		if dIndex >= len(shortened)-1 { // s only contains days
			return FAHDuration(float64(time.Hour) * 24 * days), nil
		}
	}

	duration, err := time.ParseDuration(shortened[dIndex+1:])
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return FAHDuration(duration + time.Duration(float64(time.Hour)*24*days)), nil
}

func (f FAHDuration) UnknownTime() bool {
	return f == unknowntime
}

func (f FAHDuration) String() string {
	if f.UnknownTime() {
		return unknowntimeStr
	}

	return time.Duration(f).String()
}

func (f *FAHDuration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	duration, err := ParseFAHDuration(s)
	if err != nil {
		return &json.UnmarshalTypeError{
			Value: string(b),
			Type:  reflect.TypeOf(f),
		}
	}

	*f = duration
	return nil
}

type SimulationInfo struct {
	User            string  `json:"user"`
	Team            string  `json:"team"`
	Project         int     `json:"project"`
	Run             int     `json:"run"`
	Clone           int     `json:"clone"`
	Gen             int     `json:"gen"`
	CoreType        int     `json:"core_type"`
	Core            string  `json:"core"`
	TotalIterations int     `json:"total_iterations"`
	IterationsDone  int     `json:"iterations_done"`
	Energy          int     `json:"energy"`
	Temperature     int     `json:"temperature"`
	StartTime       FAHTime `json:"start_time"`
	Timeout         int     `json:"timeout"`
	Deadline        int     `json:"deadline"`
	ETA             int     `json:"eta"`
	Progress        float64 `json:"progress"`
	Slot            int     `json:"slot"`
}

// FAHTime can be invalid, which can be checked with time.Invalid().
type FAHTime time.Time

var _ json.Unmarshaler = (*FAHTime)(nil)

const invalidTime = "<invalid>"

func ParseFAHTime(s string) (FAHTime, error) {
	if s == invalidTime {
		return FAHTime{}, nil
	}

	t, err := time.Parse(time.RFC3339, s)
	return FAHTime(t), err
}

func (t FAHTime) Invalid() bool {
	return time.Time(t).IsZero()
}

func (t FAHTime) String() string {
	if t.Invalid() {
		return invalidTime
	}

	return time.Time(t).String()
}

func (t *FAHTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	fahTime, err := ParseFAHTime(s)
	if err != nil {
		return &json.UnmarshalTypeError{
			Value: string(b),
			Type:  reflect.TypeOf(t),
		}
	}
	*t = fahTime
	return nil
}

type SlotInfo struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"`
	Description string                 `json:"description"`
	Options     map[string]interface{} `json:"options"`
	Reason      string                 `json:"reason"`
	Idle        bool                   `json:"idle"`
}

type SlotOptions struct {
	MachineID string     `json:"machine-id"`
	Paused    StringBool `json:"paused"`
}

type Info struct {
	FAHClient struct {
		Version   string
		Author    string
		Copyright string
		Homepage  string
		Date      string
		Time      string
		Revision  string
		Branch    string
		Compiler  string
		Options   string
		Platform  string
		Bits      string
		Mode      string
		Args      string
		Config    string
	}
	CBang struct {
		Date     string
		Time     string
		Revision string
		Branch   string
		Compiler string
		Options  string
		Platform string
		Bits     string
		Mode     string
	}
	System struct {
		CPU        string
		CPUID      string
		CPUs       StringInt
		Memory     string
		FreeMemory string
		Threads    string
		OSVersion  string
		HasBattery string
		OnBattery  string
		UTCOffset  string
		PID        string
		CWD        string
		OS         string
		OSArch     string
		GPUs       StringInt
		// I don't have multiple GPUs so I can't test the "GPU 0" part
	}
	LibFAH struct {
		Date     string
		Time     string
		Revision string
		Branch   string
		Compiler string
		Options  string
		Platform string
		Bits     string
		Mode     string
	}
}

func (i *Info) FromSlice(src [][]interface{}) error {
	if len(src) < 4 ||
		src[0][0] != "FAHClient" ||
		src[1][0] != "CBang" ||
		src[2][0] != "System" ||
		src[3][0] != "libFAH" {
		return errors.New("src is invalid")
	}

	primaryFields := []interface{}{
		&i.FAHClient,
		&i.CBang,
		&i.System,
		&i.LibFAH,
	}

	for i, field := range primaryFields {
		if err := readSlice(src[i], field); err != nil {
			return err
		}
	}
	return nil
}

var stringtype = reflect.TypeOf("")

func readSlice(src []interface{}, dst interface{}) error {
	infoValue := reflect.ValueOf(dst).Elem()
	for _, item := range src[1:] {
		key := item.([]interface{})[0].(string)
		field := infoValue.FieldByName(strings.ReplaceAll(key, " ", ""))
		if field.IsValid() {
			value := item.([]interface{})[1].(string)
			if field.Type() == stringtype {
				field.Set(reflect.ValueOf(value))
			} else {
				result := field.Addr().
					MethodByName("FromString").
					Call([]reflect.Value{reflect.ValueOf(value)})

				if !result[0].IsNil() {
					return result[0].Interface().(error)
				}
			}
		} else {
			log.Printf("discarded info field: %s", key)
		}
	}
	return nil
}
