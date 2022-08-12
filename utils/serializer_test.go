package utils_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/DeeStarks/conoid/utils"
)

func TestValidateConf(t *testing.T) {
	wd, _ := os.Getwd()
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
				Tunnelled: true,
			},
			expected: utils.AppConf{
				Name: "test3",
				Type: "server",
				Listeners: []string{
					"127.0.0.1:8000",
				},
				Tunnelled: true,
			},
		},
		{
			conf: utils.AppConf{
				Name:          "test4",
				Type:          "static",
				RootDirectory: "./myapp/",
				Tunnelled:     true,
			},
			expected: utils.AppConf{
				Name:          "test4",
				Type:          "static",
				Listeners:     []string{},
				RootDirectory: filepath.Join(wd, "./myapp/"),
				Tunnelled:     true,
			},
		},
		{
			conf: utils.AppConf{
				Name:          "test5",
				Type:          "mytype",
				RootDirectory: "./myapp/",
				Tunnelled:     true,
			},
			expected: utils.AppConf{},
		},
		{
			conf: utils.AppConf{
				Name:          "test3",
				Type:          "server",
				RootDirectory: "./myapp/",
				Tunnelled:     true,
			},
			expected: utils.AppConf{},
		},
		{
			conf: utils.AppConf{
				Name:          "test3",
				Type:          "static",
				Listeners: []string{
					"127.0.0.1:8000",
				},
				Tunnelled:     true,
			},
			expected: utils.AppConf{},
		},
		{
			conf: utils.AppConf{
				Name:          "test3",
				Type:          "static",
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
