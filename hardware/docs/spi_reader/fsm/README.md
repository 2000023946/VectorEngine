Great. Step 3 is where the SPI controller starts becoming a real communication module. We are going to build the **FSM (Finite State Machine)** that controls the SPI transaction.

Remember:

* Step 2 = "Can we create a clock?"
* Step 3 = "When should we use that clock and what should happen?"
* Step 4 = "Actually move bits"

Right now we are only building the **brain**, not the data movement yet.

---

# Step 3 Goal

We want this behavior:

When another module says:

```text
start = 1
```

the SPI controller does:

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
RECEIVE_DATA
 |
 v
FINISH
 |
 v
IDLE
```

At the end:

```text
data_ready = 1
```

---

# First: Define the states

In VHDL, an FSM is usually an enumerated type.

Example:

```vhdl
type state_type is (
    IDLE,
    ASSERT_CS,
    SEND_COMMAND,
    SEND_ADDRESS,
    RECEIVE_DATA,
    FINISH
);

signal state : state_type;
```

This creates a variable that can only be one of these states.

Hardware interpretation:

You are creating a register that stores:

```
000 = IDLE
001 = ASSERT_CS
010 = SEND_COMMAND
011 = SEND_ADDRESS
100 = RECEIVE_DATA
101 = FINISH
```

---

# Second: Decide outputs for each state

This is the most important part.

The FSM controls outputs.

## IDLE

Waiting.

```
CS = 1
busy = 0
data_ready = 0
```

Meaning:

> No communication.

---

## ASSERT_CS

Start transaction.

```
CS = 0
busy = 1
```

Meaning:

> Sensor selected.

---

## SEND_COMMAND

Prepare command.

```
CS = 0
busy = 1
```

Later we add:

```
MOSI = command bits
```

---

## SEND_ADDRESS

Still talking.

```
CS = 0
busy = 1
```

---

## RECEIVE_DATA

Still selected.

```
CS = 0
busy = 1
```

---

## FINISH

Transaction complete.

```
CS = 1
busy = 0
data_ready = 1
```

---

# Third: State transition logic

This answers:

> "When do we move to the next state?"

Example:

```
IDLE

if start = 1

go to ASSERT_CS
```

---

```
ASSERT_CS

after 1 SPI clock

go to SEND_COMMAND
```

---

```
SEND_COMMAND

after 8 bits

go to SEND_ADDRESS
```

---

```
SEND_ADDRESS

after 8 bits

go to RECEIVE_DATA
```

---

```
RECEIVE_DATA

after 8 bits

go to FINISH
```

---

# Before VHDL, write the FSM like this

```
                 start
                   |
                   v

              +---------+
              |  IDLE   |
              +---------+
                   |
                   |
              +---------+
              |ASSERT_CS|
              +---------+
                   |
                   |
          +----------------+
          |SEND_COMMAND    |
          +----------------+
                   |
                   |
          +----------------+
          |SEND_ADDRESS    |
          +----------------+
                   |
                   |
          +----------------+
          |RECEIVE_DATA   |
          +----------------+
                   |
                   |
          +---------+
          | FINISH  |
          +---------+
                   |
                   v
                 IDLE
```

---

# Now let's make the VHDL skeleton

This is not the full SPI yet.

Just the FSM.

```vhdl
library ieee;
use ieee.std_logic_1164.all;


entity spi_controller is
    port(
        clk        : in std_logic;
        reset      : in std_logic;

        start      : in std_logic;

        cs         : out std_logic;
        busy       : out std_logic;
        data_ready : out std_logic
    );
end spi_controller;



architecture rtl of spi_controller is


type state_type is (
    IDLE,
    ASSERT_CS,
    SEND_COMMAND,
    SEND_ADDRESS,
    RECEIVE_DATA,
    FINISH
);


signal state : state_type;



begin


process(clk)
begin

    if rising_edge(clk) then

        if reset = '1' then

            state <= IDLE;

        else


            case state is


                when IDLE =>

                    if start = '1' then
                        state <= ASSERT_CS;
                    end if;



                when ASSERT_CS =>

                    state <= SEND_COMMAND;



                when SEND_COMMAND =>

                    state <= SEND_ADDRESS;



                when SEND_ADDRESS =>

                    state <= RECEIVE_DATA;



                when RECEIVE_DATA =>

                    state <= FINISH;



                when FINISH =>

                    state <= IDLE;


            end case;


        end if;

    end if;

end process;



end rtl;
```

---

Right now this FSM moves every FPGA clock:

```
IDLE
 |
1 clock
 |
ASSERT_CS
 |
1 clock
 |
SEND_COMMAND
```

That is NOT correct for SPI yet.

We will fix that by adding:

* SPI clock enable
* bit counter

---

# The next improvement

Currently:

```
SEND_COMMAND
```

lasts 1 FPGA clock.

But we need:

```
SEND_COMMAND

clock 1 -> bit 7
clock 2 -> bit 6
clock 3 -> bit 5
...
clock 8 -> bit 0

then move on
```

So next we add:

## Bit Counter

Something like:

```
bit_count = 0

SPI clock happens

bit_count++

if bit_count == 7:
    next state
```

---

So for now, your task for Step 3 is:

1. Create the FSM entity
2. Create the states
3. Make the state transitions
4. Simulate the state changing

Do **not** add MOSI/MISO yet.

The first waveform you want is:

```
start

____|‾‾‾|________


state:

IDLE
 |
ASSERT_CS
 |
SEND_COMMAND
 |
SEND_ADDRESS
 |
RECEIVE_DATA
 |
FINISH
 |
IDLE
```

Once that works, the next step is attaching the SPI clock and making each state last the correct number of SPI cycles.
