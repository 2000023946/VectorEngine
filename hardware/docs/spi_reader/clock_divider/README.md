Good. Step 2 is **building the SPI clock generator**. This is the foundation because, as we discussed, **SPI communication cannot happen without the clock**.

Let's start from the hardware mindset.

---

# What are we trying to build?

Your FPGA has a main clock.

For example, your DE10-Lite board has:

```
FPGA Clock = 50 MHz
```

Meaning:

```
50,000,000 clock edges every second
```

But your accelerometer does not want a 50 MHz SPI clock. It might want:

```
SPI Clock = 1 MHz
```

So we need a circuit that converts:

```
50 MHz FPGA clock
        |
        v
 Clock Divider
        |
        v
1 MHz SPI clock
```

---

# Step 1: Understand frequency division

A clock is just a signal switching:

```
FPGA clock:

_|-|_|-|_|-|_|-|_|-|_
```

Each rising edge is a clock event.

If we count FPGA clock cycles:

```
FPGA cycles:

1
2
3
4
5
6
7
8
...
```

We can say:

"After N FPGA cycles, toggle the SPI clock."

Example:

Suppose:

```
FPGA clock = 10 MHz
Want SPI = 1 MHz
```

A SPI period has:

```
1 MHz = 1,000,000 cycles/sec

period = 1 / 1,000,000

period = 1 us
```

The clock has two halves:

```
HIGH time = 0.5 us
LOW time  = 0.5 us
```

The FPGA clock period:

```
10 MHz

period = 100 ns
```

How many FPGA cycles for 0.5 us?

```
500 ns / 100 ns = 5 cycles
```

So:

```
count 5 FPGA clocks
toggle SPI clock
count 5 FPGA clocks
toggle again
```

Result:

```
FPGA:

_|_|_|_|_|_|_|_|_|_|_|_|_


SPI:

____----____----____
```

---

# Step 2: Create the counter

In VHDL, hardware needs storage.

We need:

```
counter
```

that counts FPGA clock cycles.

Something like:

```vhdl
signal counter : integer := 0;
```

Every FPGA clock:

```
counter = counter + 1
```

---

# Step 3: Decide the divider value

For DE10-Lite:

```
FPGA clock = 50 MHz
```

Let's say we want:

```
SPI clock = 1 MHz
```

Formula:

[
divider = \frac{FPGA\ frequency}{2 \times SPI\ frequency}
]

Why divide by 2?

Because we toggle twice per SPI period:

```
LOW -> HIGH
HIGH -> LOW
```

So:

```
divider =
50 MHz / (2 * 1 MHz)

= 25
```

Meaning:

Every 25 FPGA clocks:

```
toggle spi_clk
```

---

# Step 4: Write the first VHDL version

Something like:

```vhdl
architecture rtl of spi_reader is

    signal count : integer := 0;
    signal spi_clk_reg : std_logic := '0';

begin

process(clk)
begin

    if rising_edge(clk) then

        if count = 24 then

            spi_clk_reg <= not spi_clk_reg;
            count <= 0;

        else

            count <= count + 1;

        end if;

    end if;

end process;


spi_clk <= spi_clk_reg;

end rtl;
```

---

Let's walk through the hardware:

Initially:

```
count = 0
spi_clk = 0
```

FPGA clock edge:

```
count = 1
```

Next:

```
count = 2
```

...

When:

```
count = 24
```

do:

```
spi_clk <= NOT spi_clk
```

So:

```
0 becomes 1
```

Then reset:

```
count = 0
```

Repeat.

---

# Step 5: What does this create?

Input:

```
FPGA clock:

50 MHz

_|-|_|-|_|-|_|-|_
```

Output:

```
SPI clock:

1 MHz

______----______----
```

Now your FPGA has a clock signal it can send to the accelerometer.

---

# But important: SPI clock should NOT always run

This is where the FSM comes later.

Right now this creates:

```
SPI clock always running
```

But a real SPI controller does:

```
IDLE:

spi_clk = 0
(no communication)


START:

enable clock


TRANSFER:

generate 8/16/64 cycles


DONE:

stop clock
```

So eventually your clock divider gets an enable:

```vhdl
spi_enable : in std_logic;
```

Then:

```text
if spi_enable = 1
    generate SPI clock

else
    hold SPI clock low
```

---

So Step 2 is really building a small hardware block:

```
              FPGA 50MHz

                   |
                   v

          +----------------+
          | Clock Divider |
          +----------------+

                   |
                   v

             SPI CLK output

                   |
                   v

             Accelerometer
```

Your first goal is simply:

**"Can I take my FPGA clock and create a slower controlled clock?"**

Once you have that, Step 3 (FSM) will use that clock to say:

* when to pull CS low
* when to send MOSI bits
* when to read MISO bits

The clock is the heartbeat; the FSM is the brain.
