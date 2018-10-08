// Copyright 2018 NeuralSpaz All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package i2cmux_test

import (
	"fmt"
	"log"

	"github.com/NeuralSpaz/i2cmux"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"
	"periph.io/x/periph/host"
)

func Example() {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	mux, err := i2cmux.New("/dev/i2c-1", i2cmux.Address(0x40))
	if err != nil {
		log.Fatalln(err)
	}
	// register the channel with the mux
	ch, err := mux.RegisterChannel(0)
	if err != nil {
		log.Fatalln(err)
	}
	// use the mux channel like any other i2c.Bus example here taken from the
	// bmp180/280 example
	d, err := bmxx80.NewI2C(ch, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		log.Fatalf("failed to initialize bme280: %v", err)
	}
	e := physic.Env{}
	if err := d.Sense(&e); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%8s %10s %9s\n", e.Temperature, e.Pressure, e.Humidity)

}
