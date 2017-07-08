package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {

	p, err := Parse("examples/simple.cfg")
	if err != nil {
		t.Errorf("parse: %s", err)
	}

	buf, err := json.MarshalIndent(p.Config, "", "  ")
	if err != nil {
		t.Errorf("marshal: %s", err)
	}
	fmt.Printf("\n%s\n", buf)

}
