library ieee;
use ieee.std_logic_1164.all;

entity spi_shift_register is
  port (
    clk          : in std_logic;
    resetn       : in std_logic;
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

  -- Single RISING edge process ensures perfect Quartus synthesis
  process(clk)
  begin
    if rising_edge(clk) then
      -- 1. Synchronous Active-Low Reset
      if resetn = '0' then
        tx_reg    <= (others => '0');
        rx_reg    <= (others => '0');
        bit_count <= 0;
        done      <= '0';
        
      -- 2. Load Phase
      elsif load = '1' then
        tx_reg    <= tx_data;
        rx_reg    <= (others => '0');
        bit_count <= 0;
        done      <= '0';
        
      -- 3. Shift Phase
      elsif shift_enable = '1' then
        -- Shift out (MOSI) and capture in (MISO) simultaneously
        tx_reg <= tx_reg(6 downto 0) & '0';
        rx_reg <= rx_reg(6 downto 0) & spi_miso;
        
        -- Track the 8 bits
        if bit_count = 7 then
          done      <= '1';
          bit_count <= 0;
        else
          bit_count <= bit_count + 1;
          done      <= '0';
        end if;
        
      -- 4. Idle Phase
      else
        done <= '0';
      end if;
    end if;
  end process;

  -- Continuous assignments mapping internal registers to output ports
  spi_mosi <= tx_reg(7);
  rx_data  <= rx_reg;

end rtl;