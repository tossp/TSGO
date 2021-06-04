package setting

import "github.com/getsentry/sentry-go"

func tsSentry() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://2920bb255d244c86aea629f6089e2f5a@bug.tossp.com:2087/16",
		Debug:            IsDev(),
		Release:          GitVersion,
		Environment:      GetString("mod"),
		AttachStacktrace: true,
	}); err != nil {
		panic(err)
	}
}
