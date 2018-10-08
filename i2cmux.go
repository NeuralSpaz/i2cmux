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

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
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
}

// New creates a new Mux on an i2c bus given the bus name and mux i2c address
func New(name string, address uint16) (*Mux, error) {

	b, err := i2creg.Open(name)
	if err != nil {
		return nil, err
	}
	fmt.Println(b)

	m := Mux{
		numChannels:    8,
		maxClock:       400 * physic.KiloHertz,
		currentChannel: 0,
		bus:            b,
		address:        address,
	}
	err = m.bus.Tx(address, []byte{0x01}, nil)
	if err != nil {
		return nil, errors.New("Failed to init channel on mux: " + err.Error())
	}

	return &m, nil
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
func (c Channel) String() string { return strconv.Itoa(int(c.channel)) }

// Tx w or r can be omitted for one way communication
func (c Channel) Tx(addr uint16, w, r []byte) error {
	return c.mux.tx(c, addr, w, r)
}

// SetSpeed sets the clock of the Mux, Changing the channel clock will also change the underlying i2c bus clock
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
	fmt.Println("MUX TX", channel, address, w, r)
	if channel.channel != m.currentChannel {
		err := m.bus.Tx(m.address, []byte{1 << channel.channel}, nil)
		if err != nil {
			return errors.New("failed to write active channel on mux: " + err.Error())
		}
		m.currentChannel = channel.channel
	}
	if channel.clock != m.currentClock {
		err := m.bus.SetSpeed(channel.clock)
		if err != nil {
			return errors.New("failed to change speed on channel " + channel.String() + ": " + err.Error())
		}
		m.currentClock = channel.clock
	}
	return m.bus.Tx(address, w, r)
}
