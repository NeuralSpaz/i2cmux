// Copyright 2018 NeuralSpaz All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package i2cmux provides a API for using a I²C(i2c) multiplexer such as a NXP PCA9548A.
package i2cmux

import (
	"reflect"
	"testing"

	"periph.io/x/periph/conn/physic"
)

func TestNew(t *testing.T) {
	type args struct {
		name    string
		address uint16
	}
	tests := []struct {
		name    string
		args    args
		want    *Mux
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMux_RegisterChannel(t *testing.T) {
	type args struct {
		channel uint8
	}
	tests := []struct {
		name    string
		m       *Mux
		args    args
		want    Channel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.RegisterChannel(tt.args.channel)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mux.RegisterChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mux.RegisterChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_String(t *testing.T) {
	tests := []struct {
		name string
		c    Channel
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Channel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_Tx(t *testing.T) {
	type args struct {
		addr uint16
		w    []byte
		r    []byte
	}
	tests := []struct {
		name    string
		c       Channel
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Tx(tt.args.addr, tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Channel.Tx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChannel_SetSpeed(t *testing.T) {
	type args struct {
		f physic.Frequency
	}
	tests := []struct {
		name    string
		c       Channel
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SetSpeed(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("Channel.SetSpeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMux_tx(t *testing.T) {
	type args struct {
		channel Channel
		address uint16
		w       []byte
		r       []byte
	}
	tests := []struct {
		name    string
		m       *Mux
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.tx(tt.args.channel, tt.args.address, tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Mux.tx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
