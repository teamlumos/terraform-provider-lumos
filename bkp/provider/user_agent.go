// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package provider

import "runtime"

type UserAgent struct {
	goVersion string
	osName string
	osVersion string
}

func NewUserAgent() UserAgent {
	ua := UserAgent{}
	ua.goVersion = runtime.Version()
	ua.osName = runtime.GOOS
	ua.osVersion = runtime.GOARCH

	return ua
}

func (ua UserAgent) String() string {
	userAgentString := "terraform-provider-lumos/0.0.0 "
	userAgentString += "golang/" + ua.goVersion + " "
	userAgentString += ua.osName + "/" + ua.osVersion

	return userAgentString
}