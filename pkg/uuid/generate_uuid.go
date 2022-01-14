package uuid

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func GenerateUUID() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed when generate uuid : %s", err)
	}
	return uid.String(), nil
}
