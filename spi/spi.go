// Package spi provides generic interfaces for interacting with SPI devices.
// It doesn't actually have any implementation of SPI usage for
// any particular hardware; other packages (such as linuxspi) will
// provide various implementations of these interfaces.
package spi

import (
	"io"
)

// Mode represents one of the four "standard" SPI modes.
type Mode uint

// BitOrder describes one of the two orders in which bits can be transmitted
// over a serial channel.
type BitOrder uint

const (
	Mode0 Mode = 0
	Mode1 Mode = 1
	Mode2 Mode = 2
	Mode3 Mode = 3
)

const (
	MsbFirst BitOrder = 0
	LsbFirst BitOrder = 1
)

// Configurator includes all of the configuration functions that are expected
// to be available on all SPI channels.
type Configurator interface {
	SetMode(mode Mode) error
	SetBitOrder(order BitOrder) error
	SetMaxSpeedHz(speed uint32) error
}

// ReadableDevice represents an SPI device that can only be written to.
type WritableDevice interface {
	Configurator
	io.Writer
}

// ReadableDevice represents an SPI device that can only be read.
type ReadableDevice interface {
	Configurator
	io.Reader
}

// Device represents a full-featured, bidirectional SPI channel to a
// particular device.
//
// Higher-level device drivers should prefer to depend on WritableDevice
// or ReadableDevice if data will pass only in one direction, so it is
// clear to the programmer how the SPI bus will be used and thus e.g. whether
// the MISO pin needs to be connected.
type Device interface {
	Configurator
	io.Writer
	io.Reader

	// Exhange writes data while simultaneously reading it (a full-duplex
	// transfer).
    //
	// The two provided slices must have the same length, since
	// the other device will provide one bit for each bit we provide.
	Exchange(outData []byte, inData []byte) (n int, err error)

	// Request writes data and then reads it (a half-duplex request/response
	// transfer).
	//
	// The two provided slices will often have differing lengths. Data
	// obtained while writing outData will be discarded, and undefined
	// data will be sent to the other device while filling inData.
	Request(outData []byte, inData []byte) (n int, err error)
}
