package haproxyctl

import (
	"fmt"
)

// just check for health of the service and return true or false
func (c *ControlSocket) Ping() (bool, error) {

	cmd := []byte("show info\n")

	res, err := c.roundTrip(cmd)
	if err != nil {
		return false, err
	}

	if res.Len() == 0 {
		return false, fmt.Errorf("zero bytss received")
	}

	return true, nil
}
