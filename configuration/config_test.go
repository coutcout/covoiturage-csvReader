// Package configuration_test is used to test application configuration
package configuration_test

import (
	"fmt"
	"me/coutcout/covoiturage/configuration"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadingConfiguration(t *testing.T) {
	t.Run("LoadingConfiguration", func(t *testing.T) {
		configPointer, err := configuration.NewConfig("./testdata/application-dev.yaml")
		assert.NoError(t, err)

		config := (*configPointer)

		testFields(t, config)

		assert.Equal(t, "127.0.0.1", config.Server.Host)
		assert.Equal(t, "8080", config.Server.Port)
		assert.Equal(t, int64(1000000), config.Journey.Import.MaxUploadFile)
		assert.Equal(t, 10, config.Journey.Parser.WorkerPoolSize)
	})

	t.Run("File does not exist", func(t *testing.T) {
		_, err := configuration.NewConfig("./testdata/application.yaml")
		assert.Error(t, err)
	})

	t.Run("File not formatted correctly", func(t *testing.T) {
		_, err := configuration.NewConfig("./testdata/application.yaml")
		assert.Error(t, err)
	})
}

func testFields(t *testing.T, obj interface{}) {
	structType := reflect.TypeOf(obj)
	fmt.Printf("Field: %s \n", structType)

	if structType.Kind() != reflect.Struct {
		return
	}

	// now go one by one through the fields and validate their value
	structVal := reflect.ValueOf(obj)
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		// Field(i) returns i'th value of the struct
		field := structVal.Field(i)
		fieldType := structType.Field(i)
		fieldInterface := field.Interface()

		if reflect.TypeOf(fieldInterface).Kind() == reflect.Struct {
			testFields(t, fieldInterface)
		} else {
			fmt.Printf("*** Field: %s - Value: ", fieldType.Name)
			if !field.IsValid() || field.IsZero() {
				fmt.Print("*Wrong Value*")
			} else {
				fmt.Printf("%v", field)
			}
			fmt.Print("\n")
		}
		assert.True(t, field.IsValid())
		assert.False(t, field.IsZero())
	}
}
