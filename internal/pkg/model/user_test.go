package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserValidate(t *testing.T) {
	valid := &User{
		Email:    "javadalipanah@gmail.com",
		UserName: "javad",
		Password: "gotcha!!",
	}

	invalidMail := &User{
		Email:    "http://javadalipanah@gmail.c",
		UserName: "javad",
		Password: "gotcha!!",
	}

	invalidUsername := &User{
		Email:    "javadalipanah@gmail.com",
		UserName: "inv",
		Password: "gotcha!!",
	}

	invalidPass := &User{
		Email:    "javadalipanah@gmail.com",
		UserName: "javad",
		Password: "gotcha!",
	}

	assert.NoError(t, valid.Validate())
	assert.Error(t, invalidMail.Validate())
	assert.Error(t, invalidUsername.Validate())
	assert.Error(t, invalidPass.Validate())
}
