package xbmc

type Platform struct {
	OS      string
	Arch    string
	Version string
	Kodi    int
	Build   string
}

func GetPlatform() *Platform {
	retVal := Platform{}
	executeJSONRPCEx("GetPlatform", &retVal, nil)
	return &retVal
}
