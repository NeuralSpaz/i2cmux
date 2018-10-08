// Copyright 2018 NeuralSpaz All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package i2cmux provides a API for using a I²C(i2c) multiplexer such as a NXP PCA9548A.
package i2cmux

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
)

const (
	// PCA9548AAddr is the default address
	pca9548aAddr     = 0x70
	pca9548aChannels = 8
)

// Mux is i2c mux such NXP PCA9548A
type Mux struct {
	sync.Mutex
	numChannels    uint8
	maxClock       physic.Frequency
	currentClock   physic.Frequency
	currentChannel uint8
	bus            i2c.Bus
	address        uint16
	debug          bool
	resetPin       gpio.PinIO
}

// New creates a new Mux on an i2c bus given the bus name and mux i2c address
func New(name string, opts ...func(*Mux) error) (*Mux, error) {

	b, err := i2creg.Open(name)
	if err != nil {
		return nil, err
	}

	m := Mux{
		numChannels:    pca9548aChannels,
		maxClock:       400 * physic.KiloHertz,
		currentChannel: 0x00,
		bus:            b,
		address:        pca9548aAddr,
		resetPin:       "",
	}

	for _, option := range opts {
		option(&m)
	}
	fmt.Printf("%+#v\n", m)

	if m.debug {
		fmt.Printf("using i2c port %s\n", b)
		if p, ok := b.(i2c.Pins); ok {
			fmt.Printf("SDA: %s\n", p.SDA())
			fmt.Printf("SCL: %s\n", p.SCL())
		}
	}

	err = m.bus.Tx(m.address, []byte{0x01}, nil) // enables on channel on
	if err != nil {
		return nil, errors.New("Failed to init channel on mux: " + err.Error())
	}

	return &m, nil
}

// Address sets the i2c address in not using the default address of 0x40
func Address(address uint16) func(*Mux) error {
	return func(m *Mux) error {
		m.address = address
		return nil
	}
}

// Debug sets the enables mux debug output
func Debug() func(*Mux) error {
	return func(m *Mux) error {
		return m.enableDebug(true)
	}
}

func (m *Mux) enableDebug(enable bool) error {
	m.debug = enable
	return nil
}

// Channels sets the number of channels on the mux
func Channels(channels uint8) func(*Mux) error {
	return func(m *Mux) error {
		m.numChannels = channels
		return nil
	}
}

// Reset use if you have a reset pin
func Reset(pin PinIO) func(*Mux) error {
	return func(m *Mux) error {
		if m.resetPin.Name() != "" {
			m.resetPin = pin
		}
		return m.reset()
	}
}

func (m *Mux) reset() error {
	m.resetPin.Out(gpio.Low)
	m.resetPin.Out(gpio.High)
	time.Sleep(time.Millisecond * 100)
	return nil
}

// Channel implements the i2c.Bus interface through the Mux
type Channel struct {
	mux     *Mux
	channel uint8
	clock   physic.Frequency
}

// RegisterChannel returns a i2c.Bus
func (m *Mux) RegisterChannel(channel uint8) (Channel, error) {
	if channel >= m.numChannels {
		return Channel{}, errors.New("Channel number must be between 0 and " + strconv.Itoa(int(m.numChannels-1)))
	}
	return Channel{mux: m, channel: channel, clock: 100 * physic.KiloHertz}, nil
}

// String returns the channel number
func (c Channel) String() string { return "channel:" + strconv.Itoa(int(c.channel)) }

// Tx w or r can be omitted for one way communication
func (c Channel) Tx(addr uint16, w, r []byte) error {
	return c.mux.tx(c, addr, w, r)
}

// Scan the channel for i2c devices, returns a slice of i2c addresses
func (c Channel) Scan() []uint16 {
	addresses := make([]uint16, 0)
	for i := uint16(0); i < 0x77; i++ {
		r := []byte{0x00}
		err := c.Tx(i, nil, r)
		if err == nil {
			addresses = append(addresses, i)
		}
	}
	return addresses
}

// SetSpeed sets the clock of the Mux, Changing the channel clock will also change the underlying i2c bus clock
// and is not supported on all hosts.
func (c Channel) SetSpeed(f physic.Frequency) error {
	if f > c.mux.maxClock {
		return errors.New("maximum mux bus speed is " + c.mux.maxClock.String())
	}
	c.clock = f
	return nil
}

// Tx raw tx on mux
func (m *Mux) tx(channel Channel, address uint16, w, r []byte) error {
	m.Lock()
	defer m.Unlock()
	if m.debug {
		fmt.Println("MUX TX", channel.String(), address, w, r)
	}
	if address == m.address {
		return errors.New("failed to write device address conflicts with mux address")
	}
	if channel.channel != m.currentChannel {
		if m.debug {
			fmt.Printf("MUX Changing Channel from %d to %d\n", m.currentChannel, channel.channel)
		}
		err := m.bus.Tx(m.address, []byte{uint8(1 << channel.channel)}, nil)
		if err != nil {
			return errors.New("failed to write active channel on mux: " + err.Error())
		}
		m.currentChannel = channel.channel
	}
	// TODO make this work independent of platform failed on raspberry pi
	if channel.clock != m.currentClock {
		err := m.bus.SetSpeed(channel.clock)
		if err != nil {

			//	return errors.New("failed to change speed on channel " + channel.String() + ": " + err.Error())
		}
		m.currentClock = channel.clock
	}
	return m.bus.Tx(address, w, r)
}
