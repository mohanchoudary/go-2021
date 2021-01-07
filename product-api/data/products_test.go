package data

import (
	"testing"
)

func TestNew(t *testing.T) {
	p := &Product{
		ID:    123,
		SKU:   "Aws-eee-wer",
		Price: 2.3,
	}
	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}

}
