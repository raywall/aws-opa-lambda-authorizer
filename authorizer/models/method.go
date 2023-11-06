package models

import (
	"errors"
	"strings"
)

type MethodARN struct {
	ARN         string
	Environment string
	Method      string
	Path        string
}

func (method *MethodARN) Load(data string) error {
	parts := strings.Split(data, "/")

	if len(parts) < 4 {
		return errors.New("Method ARN with wrong number of slices")
	}

	method.ARN = parts[0]
	method.Environment = parts[1]
	method.Method = parts[2]
	method.Path = strings.Join(parts[3:], "/")

	return nil
}
