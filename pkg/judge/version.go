package judge

import (
	"fmt"
	goversion "github.com/hashicorp/go-version"
)

type Version struct {
	*goversion.Version
}

func (v *Version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *Version) UnmarshalText(text []byte) error {
	if err := v.Set(string(text)); err != nil {
		return err
	}

	return nil
}

func (v *Version) String() string {
	if v != nil && v.Version != nil {
		return v.Version.String()
	}
	return ""
}

func NewVersion(s string) (*Version, error) {
	if v, err := goversion.NewVersion(s); err != nil {
		return nil, err
	} else {
		return &Version{Version: v}, nil
	}
}

func NewFromGoVersion(source *goversion.Version) (*Version, error) {
	if source == nil {
		return nil, fmt.Errorf("source version can't be nil")
	}

	return &Version{Version: source}, nil
}

func (v *Version) Set(s string) error {
	goVersion, err := goversion.NewVersion(s)
	if err != nil {
		return err
	}
	v.Version = goVersion
	return nil
}

func (v *Version) Type() string {
	return "string"
}
