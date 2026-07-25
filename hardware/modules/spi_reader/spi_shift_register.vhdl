library ieee;
use ieee.std_logic_1164.all;

entity spi_shift_register is
  port (
    clk          : in std_logic;
    reset        : in std_logic;
    load         : in std_logic;
    shift_enable : in std_logic;
    tx_data      : in std_logic_vector(7 downto 0);
    spi_miso     : in std_logic;
    spi_mosi     : out std_logic;
    rx_data      : out std_logic_vector(7 downto 0);
    done         : out std_logic
  );
end spi_shift_register;

architecture rtl of spi_shift_register is

  signal tx_reg    : std_logic_vector(7 downto 0);
  signal rx_reg    : std_logic_vector(7 downto 0);
  signal bit_count : integer range 0 to 7 := 0;

begin

  -- Process 1: Transmit (MOSI) and Control on the RISING edge
  process(clk)
  begin
    if rising_edge(clk) then
      if reset = '1' then
        tx_reg    <= (others => '0');
        bit_count <= 0;
        done      <= '0';
      elsif load = '1' then
        tx_reg    <= tx_data;
        bit_count <= 0;
        done      <= '0';
      elsif shift_enable = '1' then
        tx_reg <= tx_reg(6 downto 0) & '0';
        
        if bit_count = 7 then
          done      <= '1';
          bit_count <= 0;
        else
          bit_count <= bit_count + 1;
          done      <= '0';
        end if;
      else
        done <= '0';
      end if;
    end if;
  end process;

  -- Process 2: Receive (MISO) on the FALLING edge
  process(clk)
  begin
    if falling_edge(clk) then
      if reset = '1' then
        rx_reg <= (others => '0');
      -- Only sample data when we are actively shifting
      elsif shift_enable = '1' then
        rx_reg <= rx_reg(6 downto 0) & spi_miso;
      end if;
    end if;
  end process;

  -- Continuous assignments
  spi_mosi <= tx_reg(7);
  rx_data  <= rx_reg;

end rtl;