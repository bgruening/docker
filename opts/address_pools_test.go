package opts

import (
	"testing"
)

func TestAddressPoolOpt(t *testing.T) {
	poolopt := &PoolsOpt{}
	addresspool := "base=175.30.0.0/16,size=16"
	invalidAddresspoolString := "base=175.30.0.0/16,size=16, base=175.33.0.0/16,size=24"

	if err := poolopt.Set(addresspool); err != nil {
		t.Fatal(err)
	}

	if err := poolopt.Set(invalidAddresspoolString); err == nil {
		t.Fatal(err)
	}
}
