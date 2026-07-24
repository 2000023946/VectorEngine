library ieee;
use ieee.std_logic_1164.all;
entity accelerometer_controller is

  port (

    clk   : in std_logic;
    reset : in std_logic;

    start : in std_logic;
    -- SPI Reader interface

    spi_start : out std_logic;

    spi_tx_data : out std_logic_vector(7 downto 0);

    spi_ready : in std_logic;

    spi_rx_data : in std_logic_vector(7 downto 0);

    -- Debug output

    device_id : out std_logic_vector(7 downto 0);

    done : out std_logic

  );

end accelerometer_controller;

architecture rtl of accelerometer_controller is
  type state_type is (

    IDLE,
    SEND_ADDRESS,
    WAIT_SPI,
    FINISH

  );
  signal state       : state_type := IDLE;
  signal id_register : std_logic_vector(7 downto 0);
begin

  process (clk)

  begin

    if rising_edge(clk) then
      if reset = '1' then

        state <= IDLE;

        id_register <= (others => '0');
      else
        case state is
          when IDLE =>

            if start = '1' then

              state <= SEND_ADDRESS;

            end if;

          when SEND_ADDRESS =>

            state <= WAIT_SPI;

          when WAIT_SPI =>

            if spi_ready = '1' then

              id_register <= spi_rx_data;

              state <= FINISH;

            end if;

          when FINISH =>

            state <= IDLE;

        end case;
      end if;
    end if;
  end process;

  process (state)

  begin
    spi_start <= '0';

    spi_tx_data <= (others => '0');

    done <= '0';

    case state is

      when IDLE =>

        null;

      when SEND_ADDRESS =>
        -- ADXL345 Device ID register

        spi_start <= '1';

        spi_tx_data <= "00000000";

      when WAIT_SPI =>

        null;

      when FINISH =>

        done <= '1';

    end case;
  end process;

  device_id <= id_register;
end rtl;