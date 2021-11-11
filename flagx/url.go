package flagx

import (
	"flag"
	"net/url"
)

type AbsURL struct {
	url.URL
}

var _ flag.Value = (*AbsURL)(nil)

func (a *AbsURL) Set(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	a.URL = *u
	return nil
}
