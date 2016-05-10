package haproxy

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Info struct {
	Name                        string
	Version                     string
	Release_date                string
	Nbproc                      int
	Process_num                 int
	Pid                         int
	Uptime                      string
	Uptime_sec                  int
	Memmax_MB                   int
	Ulimit_n                    int `scan:"Ulimit-n"`
	Maxsock                     int
	Maxconn                     int
	Hard_maxconn                int
	CurrConns                   int
	CumConns                    int
	CumReq                      int
	MaxSslConns                 int
	CurrSslConns                int
	CumSslConns                 int
	Maxpipes                    int
	PipesUsed                   int
	PipesFree                   int
	ConnRate                    int
	ConnRateLimit               int
	MaxConnRate                 int
	SessRate                    int
	SessRateLimit               int
	MaxSessRate                 int
	SslRate                     int
	SslRateLimit                int
	MaxSslRate                  int
	SslFrontendKeyRate          int
	SslFrontendMaxKeyRate       int
	SslFrontendSessionReuse_pct int
	SslBackendKeyRate           int
	SslBackendMaxKeyRate        int
	SslCacheLookups             int
	SslCacheMisses              int
	CompressBpsIn               int
	CompressBpsOut              int
	CompressBpsRateLim          int
	ZlibMemUsage                int
	MaxZlibMemUsage             int
	Tasks                       int
	Run_queue                   int
	Idle_pct                    int
	Node                        string `scan:"node"`
	Description                 string `scan:"description"`
}

func (c *ControlSocket) Info() (*Info, error) {
	cmd := []byte("show info\n")

	buf, err := c.roundTrip(cmd)
	if err != nil {
		return nil, err
	}

	if buf.Len() == 0 {
		return nil, fmt.Errorf("zero bytss received")
	}

	var line []byte
	m := make(map[string]string, 0)
	rdr := bufio.NewReader(buf)
Reader:
	for {
		line, err = rdr.ReadBytes(NL)
		if err != nil {
			break Reader
		}

		tmp := strings.SplitN(string(line), ":", 2)
		if len(tmp) < 2 {
			continue
		}
		tmp[1] = strings.Trim(tmp[1], " \n")
		m[tmp[0]] = tmp[1]
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	info := &Info{}
	err = ScanMap(m, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}
