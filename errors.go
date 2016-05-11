package haproxyctl

import (
	"fmt"
)

type ErrorReport struct {
	Data []byte
}

func (c *ControlSocket) ShowErrors() (*ErrorReport, error) {
	cmd := []byte("show errors\n")

	res, err := c.roundTrip(cmd)
	if err != nil {
		return nil, err
	}

	if res.Len() == 0 {
		return nil, fmt.Errorf("zero bytss received")
	}

	r := &ErrorReport{
		Data: res.Bytes(),
	}

	return r, nil
}
