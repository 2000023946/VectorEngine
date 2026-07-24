library ieee;
use ieee.std_logic_1164.all;

entity spi_shift_register is
    port(

        clk : in std_logic;
        reset : in std_logic;

        -- Control
        load : in std_logic;
        shift_enable : in std_logic;

        -- Transmit byte
        tx_data : in std_logic_vector(7 downto 0);

        -- SPI lines
        spi_miso : in std_logic;
        spi_mosi : out std_logic;

        -- Received byte
        rx_data : out std_logic_vector(7 downto 0);

        -- Indicates 8 bits transferred
        done : out std_logic

    );
end spi_shift_register;