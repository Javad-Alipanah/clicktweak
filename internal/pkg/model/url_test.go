package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlValidate(t *testing.T) {
	valid := &Url{
		Url: "alipanah.me/resume",
	}

	invalid := &Url{
		Url:        "hps://alipanah.me/resume",
		Suggestion: "abcdefgh",
	}

	invalid2 := &Url{
		Url:        "alipanah.me/resume",
		Suggestion: "abcdefgh",
	}

	assert.NoError(t, valid.Validate())
	assert.Error(t, invalid.Validate())
	assert.Error(t, invalid2.Validate())
}
