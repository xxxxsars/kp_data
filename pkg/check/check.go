package check

import (
	"errors"
	"fmt"
)

func MapKeyExist(key string, data map[string]string) error {
	if _, ok := data[key]; ok {
		return nil
	}
	return errors.New(fmt.Sprintf("the key of the [%s] not in the xlsx", key))
}
