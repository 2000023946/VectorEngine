The SPI FSM (Finite State Machine) is the **controller** of the SPI reader. Its job is to control the order of operations. It does **not** send bits itself; it tells the shift register **when to load and when to shift**.

The FSM has these states:

```text
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
 |
 v
IDLE
```

---

### 1. IDLE

The FPGA is waiting.

Outputs:

```vhdl
spi_cs = 1
busy = 0
```

Meaning:

* Sensor is not selected.
* No SPI communication is happening.

When:

```vhdl
start = 1
```

the FSM begins a transaction.

---

### 2. ASSERT_CS

The FSM selects the sensor.

It does:

```vhdl
spi_cs = 0
```

In SPI, chip select is active low.

Meaning:

```
spi_cs = 0
```

tells the accelerometer:

> "I am talking to you now."

---

### 3. LOAD_COMMAND

The FSM tells the shift register:

```vhdl
load = 1
```

The shift register loads the command byte.

Example:

```
tx_data = 10110010
```

The shift register stores:

```
tx_reg = 10110010
```

The FSM itself does not care what the bits are.

---

### 4. SEND_COMMAND

Now the FSM tells the shift register:

```vhdl
shift_enable = 1
```

The shift register starts moving bits.

Each SPI clock:

```
MOSI sends 1 bit
MISO receives 1 bit
```

After 8 bits:

```
done = 1
```

The FSM moves forward.

---

### 5. LOAD_ADDRESS

Same idea as the command.

The FSM says:

```
load = 1
```

The shift register loads the sensor register address.

Example:

```
Read X-axis register
Address = 0x28
```

---

### 6. SEND_ADDRESS

The address is shifted out:

```
FPGA ---> Sensor

MOSI:
0
1
0
1
0
0
0
0
```

After 8 bits:

```
done = 1
```

Move to receive.

---

### 7. RECEIVE_DATA

Now the FPGA wants data back.

The shift register keeps generating SPI clocks.

During each clock:

```
FPGA sends dummy bit
        |
        v
Sensor sends data bit
        |
        v
MISO
```

SPI is full duplex, so even while receiving, something must be transmitted.

After 8 bits:

```
rx_data = sensor value
done = 1
```

---

### 8. FINISH

The transaction is complete.

The FSM does:

```
spi_cs = 1
ready = 1
```

Meaning:

```
Sensor released
Data available
```

Then it returns to:

```
IDLE
```

---

### Overall idea:

The FSM is like a manager:

```
FSM
 |
 | "load this byte"
 v
Shift Register

 |
 | "send 8 bits"
 v

Shift Register

 |
 | "finished"
 v

FSM moves to next step
```

The FSM controls **when** things happen.
The shift register controls **how bits move**.
