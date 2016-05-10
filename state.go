package haproxy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type ServerState struct {
	Be_id                      int
	Be_name                    string
	Srv_id                     string
	Srv_name                   string
	Srv_addr                   string
	Srv_op_state               string
	Srv_admin_state            string
	Srv_uweight                string
	Srv_iweight                string
	Srv_time_since_last_change string
	Srv_check_status           string
	Srv_check_result           string
	Srv_check_health           string
	Srv_check_state            string
	Srv_agent_state            string
	Bk_f_forced_id             string
	Srv_f_forced_id            string
}

var stateFields = []string{
	"Be_id",
	"Be_name",
	"Srv_id",
	"Srv_name",
	"Srv_addr",
	"Srv_op_state",
	"Srv_admin_state",
	"Srv_uweight",
	"Srv_iweight",
	"Srv_time_since_last_change",
	"Srv_check_status",
	"Srv_check_result",
	"Srv_check_health",
	"Srv_check_state",
	"Srv_agent_state",
	"Bk_f_forced_id",
	"Srv_f_forced_id",
}

func (c *ControlSocket) ServerState() ([]*ServerState, error) {
	cmd := []byte("show servers state\n")

	buf, err := c.roundTrip(cmd)
	if err != nil {
		return nil, err
	}

	if buf.Len() == 0 {
		return nil, fmt.Errorf("zero bytss received")
	}

	res := make([]*ServerState, 0)

	var line []byte
	rdr := bufio.NewReader(buf)

	version_line, err := rdr.ReadBytes(NL)
	if err != nil {
		return nil, fmt.Errorf("reading version line: %s", err)
	}

	if string(version_line) != "1\n" {
		return nil, fmt.Errorf("unsupported version line: %s", version_line)
	}

Reader:
	for {
		line, err = rdr.ReadBytes(NL)
		if err != nil {
			break Reader
		}

		if line[0] == '#' {
			continue
		}
		if bytes.Compare(line, []byte{NL}) == 0 {
			break
		}

		tmp := strings.Split(string(line), " ")

		m := make(map[string]string, 0)
		l := len(stateFields)

		for i := 0; i < l; i++ {

			k := strings.ToUpper(string(stateFields[i][0])) + stateFields[i][1:]
			v := tmp[i]
			m[k] = v
		}

		state := &ServerState{}
		err = ScanMap(m, state)
		if err != nil {
			return nil, err
		}

		res = append(res, state)
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
