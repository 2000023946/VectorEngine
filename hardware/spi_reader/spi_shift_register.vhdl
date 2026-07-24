library ieee;
use ieee.std_logic_1164.all;

entity spi_shift_register is
  port (

    clk   : in std_logic;
    reset : in std_logic;

    -- Control
    load         : in std_logic;
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

architecture rtl of spi_shift_register is

  signal tx_reg : std_logic_vector(7 downto 0);
  signal rx_reg : std_logic_vector(7 downto 0);

  signal bit_count : integer range 0 to 7 := 0;

begin

  process (clk)
  begin

    if rising_edge(clk) then

      if reset = '1' then

        tx_reg <= (others => '0');
        rx_reg <= (others => '0');

        bit_count <= 0;

        done <= '0';

      elsif load = '1' then

        tx_reg <= tx_data;

        rx_reg <= (others => '0');

        bit_count <= 0;

        done <= '0';

      end if;

    end if;

  end process;