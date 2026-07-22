The **SPI shift register** is the part that actually moves the bits. The FSM tells it **when** to send or receive, but the shift register handles **how the bits move one at a time**.

Remember, SPI does not send a full byte at once. It sends:

```
1 clock cycle = 1 bit transferred
```

So if you want to send:

```
10110010
```

you need 8 SPI clock cycles.

---

## Transmitting data (MOSI)

The FPGA wants to send a command or register address to the sensor.

Example:

```
Command = 10110010
```

First, load it into the transmit shift register:

```
TX Register:

[1][0][1][1][0][0][1][0]
```

Then every SPI clock:

### Clock 1

Output the first bit:

```
MOSI = 1
```

Shift:

```
[0][1][1][0][0][1][0][x]
```

---

### Clock 2

Output:

```
MOSI = 0
```

Shift again.

---

This repeats until all 8 bits are sent.

After 8 clocks:

```
10110010
```

has completely reached the sensor.

---

## Receiving data (MISO)

Now the sensor sends data back.

Example:

Sensor sends:

```
01010101
```

Every SPI clock, the FPGA samples MISO.

Clock 1:

```
MISO = 0
```

Store it.

Clock 2:

```
MISO = 1
```

Store it.

After 8 clocks:

```
RX Register:

[0][1][0][1][0][1][0][1]
```

Now:

```
data_out = 01010101
```

---

## How it connects to the FSM

The FSM controls the shift register.

Example:

```
FSM State              Shift Register

IDLE                   Do nothing

SEND_COMMAND           Load command
                       Shift bits to MOSI

SEND_ADDRESS           Load address
                       Shift bits to MOSI

RECEIVE_DATA           Shift bits from MISO
                       Store received data

FINISH                 Output complete byte
```

---

## Signals the shift register needs

It will probably have:

### Inputs

`tx_data`

The byte we want to send.

Example:

```
10110010
```

---

`load`

Tells it:

> "Put this new byte into the register."

---

`shift_enable`

Tells it:

> "SPI clock happened, move one bit."

---

`miso`

The bit coming from the sensor.

---

### Outputs

`mosi`

The bit going to the sensor.

---

`rx_data`

The completed received byte.

---

`done`

Tells the FSM:

> "I shifted 8 bits, the byte is complete."

---

The important design idea:

* **FSM = brain** → decides the operation order.
* **Shift register = hands** → moves the bits.
* **Clock divider = heartbeat** → provides the timing.

Together:

```
Clock Divider
      |
      v
   SPI Clock
      |
      v
     FSM
      |
      v
Shift Register
      |
      +---- MOSI ---> Sensor
      |
      +<--- MISO ---- Sensor
```

The next step is designing the `spi_shift_register.vhd` entity interface so it can connect cleanly to your FSM.
