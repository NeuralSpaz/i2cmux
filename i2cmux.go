// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/golang/exp/blob/master/LICENSE

// Multiplex modifications Copyright 2017 NeuralSpaz

// Package i2cmux allows users to read from and write to a slave I2C device with a multiplexer.
package i2cmux

import (
	"golang.org/x/exp/io/i2c/driver"
)

type Multiplexer interface {
	SetPort(uint8) error
}

// //
// type pca9548a struct {
// 	address uint8
// 	bus     i2c.Device
// 	// Mutex protects port from changing during a read / write
// 	sync.Mutex
// 	port uint8
// }

// // SetPort switches the multiplexer to desired port
// func (p *pca9548a) SetPort(port uint8) error {
// 	p.Lock()
// 	if p.port == port {
// 		p.Unlock()
// 		return nil
// 	}

// 	if port < 0 || port > 7 {
// 		p.Unlock()
// 		return fmt.Errorf("error setting port to %d : port must be be 0-7", port)
// 	}
// 	if err := p.bus.Write([]byte{byte(port)}); err != nil {
// 		p.Unlock()
// 		return err
// 	}
// 	p.port = port
// 	p.Unlock()
// 	return nil
// }

const tenbitMask = 1 << 12

// Device represents an I2C device. Devices must be closed once
// they are no longer in use.
type Device struct {
	conn driver.Conn
	mux  Multiplexer
	port uint8
}

// TenBit marks an I2C address as a 10-bit address.
func TenBit(addr int) int {
	return addr | tenbitMask
}

// Read reads len(buf) bytes from the device.
func (d *Device) Read(buf []byte) error {
	if err := d.mux.SetPort(d.port); err != nil {
		return err
	}
	return d.conn.Tx(nil, buf)
}

// ReadReg is similar to Read but it reads from a register.
func (d *Device) ReadReg(reg byte, buf []byte) error {
	if err := d.mux.SetPort(d.port); err != nil {
		return err
	}
	return d.conn.Tx([]byte{reg}, buf)
}

// Write writes the buffer to the device. If it is required to write to a
// specific register, the register should be passed as the first byte in the
// given buffer.
func (d *Device) Write(buf []byte) (err error) {
	if err := d.mux.SetPort(d.port); err != nil {
		return err
	}
	return d.conn.Tx(buf, nil)
}

// WriteReg is similar to Write but writes to a register.
func (d *Device) WriteReg(reg byte, buf []byte) (err error) {
	if err := d.mux.SetPort(d.port); err != nil {
		return err
	}
	// TODO(jbd): Do not allocate, not optimal.
	return d.conn.Tx(append([]byte{reg}, buf...), nil)
}

// Tx is raw access to the device.
func (d *Device) Tx(w, r []byte) (err error) {
	if err := d.mux.SetPort(d.port); err != nil {
		return err
	}
	return d.conn.Tx(w, r)
}

// Close closes the device and releases the underlying sources.
func (d *Device) Close() error {
	if err := d.mux.SetPort(d.port); err != nil {
		return err
	}
	return d.conn.Close()
}

// Open opens a connection to an I2C device.
// All devices must be closed once they are no longer in use.
// For devices that use 10-bit I2C addresses, addr can be marked
// as a 10-bit address with TenBit.
func Open(o driver.Opener, addr int, mux Multiplexer, port uint8) (*Device, error) {
	// first go around will setup the multiplexer
	unmasked, tenbit := resolveAddr(addr)
	conn, err := o.Open(unmasked, tenbit)
	if err != nil {
		return nil, err
	}
	return &Device{conn: conn, mux: mux, port: port}, nil
}

// resolveAddr returns whether the addr is 10-bit masked or not.
// It also returns the unmasked address.
func resolveAddr(addr int) (unmasked int, tenbit bool) {
	return addr & (tenbitMask - 1), addr&tenbitMask == tenbitMask
}
