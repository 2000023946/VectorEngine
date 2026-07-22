library ieee;
use ieee.std_logic_1164.all;
entity spi_reader is
  port (
    clk   : in std_logic;
    reset : in std_logic;

    start : in std_logic;
    busy  : out std_logic;
    ready : out std_logic;

    data_out : out std_logic_vector(7 downto 0);

    spi_cs   : out std_logic;
    spi_clk  : out std_logic;
    spi_mosi : out std_logic;
    spi_miso : in std_logic
  );
end spi_reader;

architecture rtl of spi_reader is

  signal count       : integer   := 0;
  signal spi_clk_reg : std_logic := '0';

begin

  process (clk)
  begin

    if rising_edge(clk) then

      if count = 24 then

        spi_clk_reg <= not spi_clk_reg;
        count       <= 0;

      else

        count <= count + 1;

      end if;

    end if;

  end process;
  spi_clk <= spi_clk_reg;

end rtl;
