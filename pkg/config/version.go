package config

import (
	goversion "github.com/hashicorp/go-version"
)

type Version struct {
	*goversion.Version
}

func (v *Version) String() string {
	if v.Version != nil {
		return v.Version.String()
	}
	return ""
}

func (v *Version) Set(s string) error {
	version, err := goversion.NewVersion(s)
	if err != nil {
		return err
	}
	v.Version = version
	return nil
}

func (v *Version) SetFromVersion(new *goversion.Version) error {
	v.Version = new
	return nil
}

func (v *Version) Type() string {
	return "string"
}

func NewVersion() *Version {
	return &Version{}
}
