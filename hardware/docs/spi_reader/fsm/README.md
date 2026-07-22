The FSM (Finite State Machine) is basically the **controller that decides what the SPI module should do and when**. It moves through a sequence of states that represent the steps of talking to the sensor.

The flow is:

**IDLE:** The SPI module is waiting. Nothing is happening.

**ASSERT_CS:** The FPGA pulls `spi_cs` low, which tells the sensor "I want to communicate with you."

**SEND_COMMAND:** The FPGA sends the command (for example, "read data") to the sensor through MOSI.

**SEND_ADDRESS:** The FPGA sends the register address it wants to read from (for example, the X-axis data register).

**RECEIVE_DATA:** The FPGA receives the sensor data through MISO while generating clock cycles.

**FINISH:** The communication is complete. The FPGA releases `spi_cs`, signals that data is ready, and returns to IDLE.

The outputs are controlled by the current state:

* **spi_cs:** Selects the sensor. `0` means the sensor is active and listening; `1` means it is inactive.
* **busy:** Tells other modules that the SPI controller is currently communicating (`1`) or available (`0`).
* **ready:** A pulse/signal telling the rest of the FPGA that new sensor data has been received.
* **spi_mosi:** Data line from FPGA → sensor (commands and addresses).
* **spi_miso:** Data line from sensor → FPGA (sensor measurements).

So the FSM is the **brain**, and the outputs are the **signals it controls to communicate with the sensor**.
