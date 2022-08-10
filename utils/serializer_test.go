package utils_test

import (
	"reflect"
	"testing"

	"github.com/DeeStarks/conoid/utils"
)

func TestCValidateConf(t *testing.T) {
	tests := []struct {
		conf     utils.AppConf
		expected utils.AppConf
	}{
		{
			conf: utils.AppConf{
				Name:          "test1",
				Type:          "server",
				RootDirectory: ".",
			},
			expected: utils.AppConf{},
		},
		{
			conf: utils.AppConf{
				Name: "test2",
				Type: "static",
				Listeners: []string{
					"127.0.0.1:8000",
				},
			},
			expected: utils.AppConf{},
		},
		{
			conf: utils.AppConf{
				Name: "test3",
				Type: "server",
				Listeners: []string{
					"127.0.0.1:8000",
				},
				RemoteServer: "124.567.785.34:33004",
				Tunnelled:    true,
			},
			expected: utils.AppConf{
				Name: "test3",
				Type: "server",
				Listeners: []string{
					"127.0.0.1:8000",
				},
				RemoteServer: "124.567.785.34:33004",
				Tunnelled:    false,
			},
		},
		{
			conf: utils.AppConf{
				Name:          "test4",
				Type:          "static",
				RootDirectory: "./myapp/",
				RemoteServer:  "",
				Tunnelled:     true,
			},
			expected: utils.AppConf{
				Name:          "test4",
				Type:          "static",
				Listeners:     []string{},
				RootDirectory: "./myapp/",
				RemoteServer:  "",
				Tunnelled:     true,
			},
		},
		{
			conf: utils.AppConf{
				Name:          "test5",
				Type:          "mytype",
				RootDirectory: "./myapp/",
				RemoteServer:  "",
				Tunnelled:     true,
			},
			expected: utils.AppConf{},
		},
	}

	for _, tc := range tests {
		conf, _ := utils.ValidateConf(tc.conf)
		if !reflect.DeepEqual(conf, tc.expected) {
			t.Errorf("Expected %v, got %v", tc.expected, conf)
		}
	}
}
