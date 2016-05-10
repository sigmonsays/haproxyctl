package haproxy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Stat struct {
	Pxname         string
	Svname         string
	Qcur           int
	Qmax           int
	Scur           int
	Smax           int
	Slim           int
	Stot           int
	Bin            int
	Bout           int
	Dreq           int
	Dresp          int
	Ereq           int
	Econ           int
	Eresp          int
	Wretr          int
	Wredis         int
	Status         string
	Weight         int
	Act            int
	Bck            int
	Chkfail        int
	Chkdown        int
	Lastchg        int
	Downtime       int
	Qlimit         int
	Pid            int
	Iid            int
	Sid            int
	Throttle       int
	Lbtot          int
	Tracked        int
	Type           int
	Rate           int
	Rate_lim       int
	Rate_max       int
	Check_status   string
	Check_code     int
	Check_duration int
	Hrsp_1Xx       int
	Hrsp_2Xx       int
	Hrsp_3Xx       int
	Hrsp_4Xx       int
	Hrsp_5Xx       int
	Hrsp_other     int
	Hanafail       int
	Req_rate       int
	Req_rate_Max   int
	Req_tot        int
	Cli_abrt       int
	Srv_abrt       int
	Comp_in        int
	Comp_out       int
	Comp_byp       int
	Comp_Rsp       int
	Lastsess       int
	Last_chk       string
	Last_agt       int
	Qtime          int
	Ctime          int
	Rtime          int
	Ttime          int
}

var statFields = []string{"pxname", "svname", "qcur", "qmax", "scur", "smax", "slim", "stot", "bin", "bout", "dreq", "dresp", "ereq", "econ", "eresp", "wretr", "wredis", "status", "weight", "act", "bck", "chkfail", "chkdown", "lastchg", "downtime", "qlimit", "pid", "iid", "sid", "throttle", "lbtot", "tracked", "type", "rate", "rate_lim", "rate_max", "check_status", "check_code", "check_duration", "hrsp_1xx", "hrsp_2xx", "hrsp_3xx", "hrsp_4xx", "hrsp_5xx", "hrsp_other", "hanafail", "req_rate", "req_rate_max", "req_tot", "cli_abrt", "srv_abrt", "comp_in", "comp_out", "comp_byp", "comp_rsp", "lastsess", "last_chk", "last_agt", "qtime", "ctime", "rtime", "ttime"}

func (c *ControlSocket) Stat() ([]*Stat, error) {
	cmd := []byte("show stat\n")

	buf, err := c.roundTrip(cmd)
	if err != nil {
		return nil, err
	}

	if buf.Len() == 0 {
		return nil, fmt.Errorf("zero bytss received")
	}

	res := make([]*Stat, 0)

	var line []byte
	rdr := bufio.NewReader(buf)
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

		tmp := strings.Split(string(line), ",")

		m := make(map[string]string, 0)
		l := len(statFields)
		for i := 0; i < l; i++ {

			k := strings.ToUpper(string(statFields[i][0])) + statFields[i][1:]
			v := tmp[i]
			m[k] = v
		}

		c.log("map %+v", m)

		stat := &Stat{}
		err = ScanMap(m, stat)
		if err != nil {
			return nil, err
		}

		res = append(res, stat)
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
