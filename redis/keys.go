package redis

const (
	CODE_KEY     = "Code_%s_%s"  // appKey + code
	CODE_EXP     = 300           // 5min
	RESOURCE_KEY = "Resource_%s" // TOKEN_ + appKey

	PLATFORM_INFO_KEY = "Platform_info_%s" // app_key
)
