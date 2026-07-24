# SPI Reader

## Overview

The SPI Reader module handles communication between the FPGA and an external SPI sensor, such as an accelerometer.

The module manages the SPI communication process by generating the SPI clock, controlling the chip select signal, and managing the communication sequence using a finite state machine (FSM).

## Structure

The SPI Reader is divided into smaller hardware modules:

- `spi_reader.vhd`
  - Top-level module that connects all SPI components together.

- `spi_fsm.vhd`
  - Controls the SPI communication states and determines the order of operations.

- `spi_clock_divider.vhd`
  - Generates the SPI clock from the FPGA system clock.

## Communication Flow

The SPI reader follows this sequence:

1. Wait for a start signal.
2. Select the sensor using chip select (CS).
3. Send the read command.
4. Send the sensor register address.
5. Receive sensor data.
6. Signal that the data is ready.

## Documentation

More detailed information about the SPI Reader design, FSM states, timing, and internal operation can be found in:

`docs/spi_reader/`