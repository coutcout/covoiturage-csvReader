// Package configuration_test is used to test application configuration
package configuration_test

import (
	"me/coutcout/covoiturage/configuration"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseParameter(t *testing.T) {

	var tests = []struct {
		args []string
		expectedParams *configuration.Parameters
		shouldHaveError bool
	}{
		{
			[]string{"-config", "./testdata/application-dev.yaml"},
			&configuration.Parameters{ConfigFilePath: "./testdata/application-dev.yaml"},
			false,
		},
		{
			[]string{"-config", },
			nil,
			true,
		},
		{
			[]string{"-config", "./testdata/unknown-file.yaml"},
			nil,
			true,
		},
		{
			[]string{"-config", "./testdata"},
			nil,
			true,
		},
	}

	for _, test := range tests {
		t.Run(strings.Join(test.args, " "), func(t *testing.T) {
			params, _, err := configuration.ParseFlag("test", test.args)
		
			if test.shouldHaveError{
				assert.Error(t, err)
				assert.Nil(t, params)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedParams, params)
			}
		})
	}
	


}