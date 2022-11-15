// Copyright © 2022 Intel Corporation. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause

package driver

import (
	"device-scale/scale"
	"reflect"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/stretchr/testify/require"
)

func Test_scaleDevice_readWeight(t *testing.T) {

	config := scale.Config{
		PortName:        "/dev/tty.usbserial-test",
		BaudRate:        9600,
		DataBits:        7,
		StopBits:        1,
		MinimumReadSize: 1,
		ParityMode:      2,
		TimeOutMilli:    500,
	}

	testDevice := scale.InitializeMockDevice(&config)
	tests := []struct {
		name          string
		serialDevice  *scale.MockDevice
		testCaseIndex int
		want          map[string]interface{}
		wantErr       bool
	}{
		{
			name:          "valid case",
			serialDevice:  testDevice,
			testCaseIndex: 0,
			want:          map[string]interface{}{"status": "OK", "total": 2.494, "units": "LB"},
			wantErr:       false,
		},
		{
			name:          "status Scale at Zero",
			serialDevice:  testDevice,
			testCaseIndex: 1,
			want:          nil,
			wantErr:       false,
		},
		{
			name:          "invalid reading but status OK",
			serialDevice:  testDevice,
			testCaseIndex: 5,
			want:          map[string]interface{}{"status": "OK", "total": 0.0, "units": "LB"},
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.serialDevice.TestCase = tt.testCaseIndex
			device := &scaleDevice{
				serialDevice: tt.serialDevice,
			}
			got, err := device.readWeight()
			if (err != nil) != tt.wantErr {
				t.Errorf("scaleDevice.readWeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scaleDevice.readWeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newScaleDevice(t *testing.T) {
	tests := []struct {
		name       string
		serialPort string
		lc         logger.LoggingClient
		config     map[string]string
		isEmpty    bool
	}{
		{
			name:       "valid case",
			serialPort: "testSerialPort",
			lc:         logger.NewMockClient(),
			config:     getDefaultDriverConfig(),
			isEmpty:    false,
		},
		{
			name:       "nil config",
			serialPort: "testSerialPort",
			lc:         logger.NewMockClient(),
			config:     nil,
			isEmpty:    true,
		},
		{
			name:       "missing TimeOutMilli from config",
			serialPort: "testSerialPort",
			lc:         logger.NewMockClient(),
			config: map[string]string{
				"SimulatorPort": "8081",
				"ScaleID":       "123",
				"LaneID":        "123",
			},
			isEmpty: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newScaleDevice(tt.serialPort, tt.lc, tt.config)

			if tt.isEmpty {
				require.Nil(t, got)
			} else {
				require.NotEmpty(t, got)
			}
		})
	}
}