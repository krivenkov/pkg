// Package global contain mostly common application logic
package global

type Info struct {
	AppName string
}

func NewInfo(appName string) func() Info {
	return func() Info {
		return Info{
			AppName: appName,
		}
	}
}
