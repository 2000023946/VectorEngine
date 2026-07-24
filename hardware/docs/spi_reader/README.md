## SPI Reader Module Description

The `spi_reader` is the top-level module that controls communication between the FPGA and an accelerometer sensor using the SPI protocol. It is built from three smaller modules: a clock divider, an FSM controller, and a shift register.

### 1. SPI Clock Divider

The FPGA runs at a high-frequency clock, but the sensor requires a slower SPI clock. The clock divider reduces the FPGA clock frequency and generates the `spi_clk` signal used by the SPI bus.

```
FPGA Clock
     |
     v
Clock Divider
     |
     v
SPI Clock
```

The generated SPI clock is shared with both the FSM and shift register so all SPI operations happen at the correct timing.

---

### 2. SPI FSM Controller

The FSM controls the order of the SPI transaction. It does not move data; it only decides what operation should happen.

The sequence is:

```
IDLE
 |
 | start = 1
 v
ASSERT_CS
 |
 v
LOAD_COMMAND
 |
 v
SEND_COMMAND
 |
 v
LOAD_ADDRESS
 |
 v
SEND_ADDRESS
 |
 v
RECEIVE_DATA
 |
 v
FINISH
```

During each state, the FSM controls:

* `spi_cs` → selects the sensor
* `load` → tells the shift register to load a byte
* `shift_enable` → tells the shift register to transfer bits
* `done` → tells the FSM when 8 bits have been transferred

---

### 3. SPI Shift Register

The shift register performs the actual SPI communication.

It handles:

* Sending data through `spi_mosi`
* Receiving data through `spi_miso`

When the FSM sends:

```
load = 1
```

the shift register loads the command/address byte.

When the FSM sends:

```
shift_enable = 1
```

the shift register moves one bit:

```
TX Register:

10110010

   |
   v

Send first bit through MOSI
```

At the same time, it receives one bit from the sensor:

```
MISO --> RX Register
```

After 8 SPI clock cycles:

```
done = 1
```

is sent back to the FSM.

---

## Complete Data Flow

The full transaction works like this:

```
FPGA
 |
 | start
 v

SPI FSM
 |
 | load command
 v

Shift Register
 |
 | MOSI
 v

Accelerometer
 |
 | MISO
 v

Shift Register
 |
 | rx_data
 v

spi_reader output
```

---

## Module Responsibilities

| Module               | Responsibility                    |
| -------------------- | --------------------------------- |
| `spi_clock_divider`  | Creates slower SPI clock          |
| `spi_fsm`            | Controls SPI transaction sequence |
| `spi_shift_register` | Sends and receives bits           |
| `spi_reader`         | Connects all modules together     |

The final `spi_reader` acts as a complete SPI controller. The rest of the FPGA design only needs to provide `start` and read `data_out`; the SPI communication is handled internally.
