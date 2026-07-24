The **SPI shift register** is the part that actually moves the bits between the FPGA and the accelerometer.

The FSM is the "manager" telling it what to do, but the shift register is the "worker" that physically sends and receives bits.

Its job is:

1. Store the byte we want to send (`tx_reg`)
2. Send one bit at a time through MOSI
3. Receive one bit at a time through MISO
4. Store the received byte (`rx_reg`)

---

## Main signals

Your shift register has:

```vhdl
load
shift_enable
tx_data
spi_miso
spi_mosi
rx_data
done
```

### `load`

The FSM uses this when it wants to put a new byte into the transmit register.

Example:

```text
FSM:
load = 1
tx_data = 10110010
```

The shift register does:

```
tx_reg = 10110010
```

Now it is ready to send.

---

### `shift_enable`

This tells the shift register:

> "Move one bit now."

Every SPI clock while `shift_enable = 1`, it does one transfer.

---

## Transmitting data (MOSI)

Suppose:

```
tx_reg = 10110010
```

SPI sends the **leftmost bit first**.

The first bit:

```
10110010
^
|
send this bit
```

So:

```
MOSI = 1
```

Then the register shifts:

```
10110010
 |
 v

01100100
```

Next clock:

```
01100100
 ^
 |
 send 0
```

Continue:

```
10110010
01100100
11001000
10010000
00100000
01000000
10000000
00000000
```

After 8 clocks, the whole byte has been transmitted.

---

## Receiving data (MISO)

At the same time, SPI is full duplex.

While you send a bit:

```
FPGA --------------> Sensor
       MOSI
```

the sensor sends a bit back:

```
FPGA <-------------- Sensor
       MISO
```

Your code:

```vhdl
rx_reg <= rx_reg(6 downto 0) & spi_miso;
```

means:

"Shift everything left and put the new incoming bit on the right."

Example:

Starting:

```
rx_reg = 00000000
```

Receive:

```
1
```

becomes:

```
00000001
```

Receive:

```
0
```

becomes:

```
00000010
```

Receive:

```
1
```

becomes:

```
00000101
```

After 8 bits:

```
rx_reg = sensor data
```

---

## Bit counter

You have:

```vhdl
signal bit_count : integer range 0 to 7;
```

This counts how many bits have moved.

Example:

```
Clock 1 -> bit_count = 0
Clock 2 -> bit_count = 1
Clock 3 -> bit_count = 2
...
Clock 8 -> bit_count = 7
```

When:

```vhdl
bit_count = 7
```

the byte is complete:

```vhdl
done <= '1';
```

The FSM sees:

```
done = 1
```

and moves to the next state.

---

## Full operation example

The FSM says:

```
LOAD_COMMAND
```

Shift register:

```
tx_reg = 10110010
```

Then FSM says:

```
SEND_COMMAND
shift_enable = 1
```

The shift register does:

```
Clock 1:
MOSI = 1

Clock 2:
MOSI = 0

Clock 3:
MOSI = 1

...

Clock 8:
MOSI = 0
```

At the same time:

```
MISO bits -> rx_reg
```

After 8 clocks:

```
done = 1
```

---

So the simple way to remember:

```
FSM = decides WHAT happens

Shift Register = moves the bits

Clock Divider = decides WHEN bits move
```

The shift register is basically the bridge between the FPGA's internal parallel data (`10110010`) and SPI's serial data line (one bit at a time).
