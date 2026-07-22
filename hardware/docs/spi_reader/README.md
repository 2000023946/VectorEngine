Good. Let's design the **SPI accelerometer reader module** like you would in a real FPGA system. Think of this module as having two sides:

1. **External interface (API)** → how the rest of the FPGA talks to it
2. **Internal implementation** → the state machine that actually generates SPI transactions

---

# 1. SPI Reader Module API

The SPI reader's job is:

> "Whenever I need acceleration data, communicate with the sensor, receive the bytes, and output a clean digital value."

The rest of your FPGA should **not care about SPI timing**. It should just see:

* request data
* wait
* receive data

---

## External API

A clean interface:

```vhdl
entity spi_accel_reader is
    port(
        -- FPGA clock
        clk          : in  std_logic;
        reset        : in  std_logic;

        -- Control Interface
        start        : in  std_logic;
        busy         : out std_logic;
        data_ready   : out std_logic;

        -- Sensor Data
        accel_x      : out std_logic_vector(15 downto 0);
        accel_y      : out std_logic_vector(15 downto 0);
        accel_z      : out std_logic_vector(15 downto 0);

        -- SPI Physical Interface
        spi_cs       : out std_logic;
        spi_clk      : out std_logic;
        spi_mosi     : out std_logic;
        spi_miso     : in  std_logic
    );
end spi_accel_reader;
```

---

Now let's break down each signal.

---

# Control Interface

## `start`

Input from your controller.

Meaning:

> "Go read the accelerometer now."

Example:

```
start = 1
```

The SPI module begins a transaction.

You pulse this:

```
      ______
_____|      |_____
     start
```

Not held high forever.

---

## `busy`

Output.

Means:

> "I am currently talking to the sensor."

Example:

```
start
 |
 v

busy:
____████████____
```

While busy:

* do not send another request
* ignore start

---

## `data_ready`

Output.

Means:

> "The new acceleration values are available."

Example:

```
busy:

████████

data_ready:

        _
_______| |_____
```

A single clock pulse.

---

## Sensor Data

```
accel_x
accel_y
accel_z
```

Example:

```
X = 0x01AF
Y = 0xFF20
Z = 0x03A0
```

These are the values you push into your FIFO.

---

# SPI Physical Interface

Now the actual SPI wires.

## CS (Chip Select)

Tells the sensor:

> "I am talking to you."

Normally:

```
Idle:

CS = 1
```

Transaction:

```
CS = 0
```

Example:

```
CS

____
    |________________|
```

---

## SPI Clock

The FPGA generates the clock.

Example:

```
spi_clk

___     ___     ___
   |___|   |___|
```

Each edge transfers one bit.

---

## MOSI

Master Out Slave In.

FPGA → Sensor.

Used for:

* register address
* read command

Example:

```
FPGA sends:

10100001
```

---

## MISO

Master In Slave Out.

Sensor → FPGA.

Used for:

* sensor data

Example:

```
Sensor sends:

00011010
```

---

# 2. SPI Transaction

Let's say the accelerometer uses:

```
Register 0x32 = X-axis low byte
Register 0x33 = X-axis high byte
```

A read operation looks like:

```
CS LOW

FPGA -> Sensor

[READ COMMAND]
[REGISTER ADDRESS]

Sensor -> FPGA

[DATA BYTE]

CS HIGH
```

Physically:

```
        CS
        |
________|----------------

MOSI:

READ
10110010
REGISTER
00110010


MISO:

xxxxxxxx
00011100
```

---

# 3. Internal Architecture

Inside the module:

```
             start
               |
               v
        +--------------+
        | Control FSM |
        +--------------+
               |
               |
        +--------------+
        | SPI Shift    |
        | Register     |
        +--------------+
               |
               |
        +--------------+
        | Clock        |
        | Divider      |
        +--------------+
               |
               |
          SPI Pins
```

You mainly need:

1. FSM
2. Clock divider
3. Shift registers

---

# Clock Divider

Your FPGA might run:

```
50 MHz
```

But SPI might need:

```
1 MHz
```

So:

```
50 MHz FPGA clock

        |
        v

divide by 50

        |
        v

1 MHz SPI clock
```

Implementation:

```vhdl
if counter = DIVIDER then
    spi_clk <= not spi_clk;
end if;
```

---

# 4. State Machine

The FSM controls the entire transaction.

Something like:

```
IDLE

 |
 | start
 v

ASSERT_CS

 |
 v

SEND_COMMAND

 |
 v

SEND_ADDRESS

 |
 v

READ_DATA

 |
 v

STORE_RESULT

 |
 v

DEASSERT_CS

 |
 v

DONE

 |
 v

IDLE
```

---

# State 1: IDLE

Waiting.

Outputs:

```
CS = 1
busy = 0
data_ready = 0
```

Nothing happens.

---

# State 2: ASSERT_CS

When:

```
start = 1
```

Move here.

Actions:

```
CS = 0
busy = 1
```

Now sensor listens.

---

# State 3: SEND_COMMAND

Load shift register:

```
10110000
```

Example:

```
Read command
```

Then shift:

```
bit 7
 |
 v

10110000

 ^
 |
 MOSI
```

Every SPI clock:

```
shift left
```

---

# State 4: SEND_ADDRESS

Example:

Need register:

```
0x32
```

Load:

```
00110010
```

Shift out.

---

# State 5: READ_DATA

Now sensor sends data.

FPGA samples MISO.

Example:

Clock edges:

```
SPI CLK:

_|-|_|-|_|-|_

MISO:

1 0 1 0 1 1 0 0
```

Shift into register:

```
received <= received & miso
```

After 8 clocks:

```
10101100
```

---

# State 6: STORE_RESULT

Save:

```
accel_x <= received_data;
```

---

# State 7: DONE

Generate pulse:

```
data_ready = 1
```

Tell FIFO:

> "New sample available."

Then return:

```
IDLE
```

---

# Full Timing Example

```
start
 |
 |
 v

CS     ______|----------------|_____

CLK       _-_-_-_-_-_-_-_-_

MOSI      CMD ADDRESS

MISO              DATA DATA DATA


busy   ____|----------------|____


ready                    _
________________________| |____
```

---

# How it connects to FIFO later

Your SPI module becomes a producer:

```
SPI MODULE

accel_x
accel_y
accel_z
    |
    |
data_ready
    |
    v

FIFO WRITE ENABLE
```

The FIFO doesn't know SPI exists.

It only sees:

```
write_enable = data_ready
write_data = accel packet
```

---

So conceptually:

**SPI module = hardware driver**

API:

```
start  ---> "read sensor"
busy    ---> "I'm working"
ready   ---> "data available"
data    ---> "here is acceleration"
```

Internally:

```
FSM controls:
    CS
    MOSI
    CLK
    MISO sampling
```

This is exactly the type of abstraction you want in FPGA design: the outside sees a simple software-like API, while the inside handles the ugly timing.
