package protocol

// 定义TR369协议中常用的常量

// DefaultTR369Subprotocol TR369协议的默认WebSocket子协议
const DefaultTR369Subprotocol = "tr-369-1-0"

// CommonEventTypes 常用的事件类型
const (
	EventBootstrapDone    = "bootstrap-done"
	EventConnectionRequest = "connection-request"
	EventReboot           = "reboot"
	EventValueChange      = "value-change"
	EventObjectCreation   = "object-creation"
	EventObjectDeletion   = "object-deletion"
	EventTransferComplete = "transfer-complete"
	EventDownloadComplete = "download-complete"
	EventUploadComplete   = "upload-complete"
	EventTransferFailure  = "transfer-failure"
)

// CommonParameterPaths 常用的参数路径
const (
	// 设备信息相关
	ParamDeviceInfoManufacturer     = "Device.DeviceInfo.Manufacturer"
	ParamDeviceInfoModelName        = "Device.DeviceInfo.ModelName"
	ParamDeviceInfoSerialNumber     = "Device.DeviceInfo.SerialNumber"
	ParamDeviceInfoFirmwareVersion  = "Device.DeviceInfo.SoftwareVersion"
	ParamDeviceInfoHardwareVersion  = "Device.DeviceInfo.HardwareVersion"
	ParamDeviceInfoProvisioningCode = "Device.DeviceInfo.ProvisioningCode"

	// 设备管理相关
	ParamManagementServerURL        = "Device.ManagementServer.URL"
	ParamManagementServerUsername   = "Device.ManagementServer.Username"
	ParamManagementServerPassword   = "Device.ManagementServer.Password"
	ParamManagementServerPeriodicInformInterval = "Device.ManagementServer.PeriodicInformInterval"
	ParamManagementServerConnectionRequestURL = "Device.ManagementServer.ConnectionRequestURL"

	// 连接状态相关
	ParamConnectionStatus        = "Device.Interface.X_ConnectionStatus"
	ParamConnectionUpTime        = "Device.Interface.X_UpTime"

	// 系统相关
	ParamSystemUptime            = "Device.DeviceInfo.UpTime"
	ParamSystemMemoryTotal       = "Device.System.Memory.Total"
	ParamSystemMemoryAvailable   = "Device.System.Memory.Available"
	ParamSystemCPUUtilization    = "Device.System.CPU.Utilization"
)

// DefaultParameterValues 默认参数值
var DefaultParameterValues = map[string]interface{}{
	ParamDeviceInfoManufacturer:               "ExampleVendor",
	ParamDeviceInfoModelName:                  "TR369Client",
	ParamDeviceInfoSerialNumber:               "SN123456",
	ParamDeviceInfoFirmwareVersion:            "1.0.0",
	ParamDeviceInfoHardwareVersion:            "1.0",
	ParamDeviceInfoProvisioningCode:           "PROV123",
	ParamManagementServerURL:                  "wss://tr069-server.example.com:7547",
	ParamManagementServerUsername:             "admin",
	ParamManagementServerPassword:             "password",
	ParamManagementServerPeriodicInformInterval: 3600,
	ParamConnectionStatus:                     "Connected",
	ParamSystemUptime:                         0,
	ParamSystemMemoryTotal:                    536870912, // 512MB in bytes
	ParamSystemMemoryAvailable:                268435456, // 256MB in bytes
	ParamSystemCPUUtilization:                 10.5,
}

// ErrorMessages 错误消息映射
var ErrorMessages = map[int]string{
	ErrorCodeNoError:              "No error",
	ErrorCodeInvalidParameterName: "Invalid parameter name",
	ErrorCodeInvalidParameterType: "Invalid parameter type",
	ErrorCodeInvalidParameterValue: "Invalid parameter value",
	ErrorCodeInternalError:        "Internal error",
	ErrorCodeResourceBusy:         "Resource busy",
}

// TR369Capabilities TR369协议支持的能力
var TR369Capabilities = []string{
	"GetParameterValues",
	"SetParameterValues",
	"GetParameterNames",
	"GetParameterAttributes",
	"SetParameterAttributes",
	"AddObject",
	"DeleteObject",
	"Reboot",
	"Download",
	"Upload",
	"ScheduleInform",
	"GetRPCMethods",
	"FactoryReset",
}

// MessageHeaderNames WebSocket消息头名称
const (
	HeaderNameAuthorization = "Authorization"
	HeaderNameSessionID     = "Session-ID"
	HeaderNameTransactionID = "Transaction-ID"
	HeaderNameCapabilities  = "Capabilities"
)

// InformEventStructured 结构化的INFORM事件类型
type InformEventStructured struct {
	EventCode  string `json:"eventCode"`
	CommandKey string `json:"commandKey,omitempty"`
}

// GetEventCode 获取事件代码
func GetEventCode(event string) string {
	eventCodes := map[string]string{
		EventBootstrapDone:    "0 BOOTSTRAP",
		EventConnectionRequest: "2 CONNECTION REQUEST",
		EventReboot:           "1 BOOT",
		EventValueChange:      "4 VALUE CHANGE",
	}
	
	if code, exists := eventCodes[event]; exists {
		return code
	}
	return "999 CUSTOM EVENT"
}