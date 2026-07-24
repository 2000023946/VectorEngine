library ieee;
use ieee.std_logic_1164.all;


entity spi_clock_divider is
    port(
        clk       : in std_logic;
        reset     : in std_logic;

        spi_clk   : out std_logic
    );
end spi_clock_divider;



architecture rtl of spi_clock_divider is

    signal count       : integer := 0;
    signal spi_clk_reg : std_logic := '0';

begin


process(clk)

begin

    if rising_edge(clk) then

        if reset = '1' then

            count       <= 0;
            spi_clk_reg <= '0';


        elsif count = 24 then

            count       <= 0;
            spi_clk_reg <= not spi_clk_reg;


        else

            count <= count + 1;

        end if;

    end if;

end process;


spi_clk <= spi_clk_reg;


end rtl;