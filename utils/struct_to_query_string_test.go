package utils_test

import (
	"fmt"
	"testing"

	"vobook/tests/assert"
	"vobook/utils"
)

func TestStructToQueryString(t *testing.T) {
	type A struct {
		B string `json:"b"`
		C int    `json:"c"`
	}
	st := A{
		B: "C",
		C: 1,
	}

	result, err := utils.StructToQueryString(st)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equals(t, fmt.Sprintf(`b=%s&c=%d`, st.B, st.C), result)
}
