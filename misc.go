package haproxyctl

func (c *ControlSocket) ClearCounters(all bool) error {
	var clear string
	if all {
		clear = "clear counters all\n"
	} else {
		clear = "clear counters\n"
	}
	cmd := []byte(clear)

	_, err := c.roundTrip(cmd)
	if err != nil {
		return err
	}

	return nil
}
