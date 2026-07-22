library ieee;
use ieee.std_logic_1164.all;


entity spi_reader is

port(

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



begin



-- SPI clock generator

clock_divider : entity work.spi_clock_divider

port map(

    clk => clk,
    reset => reset,

    spi_clk => spi_clk

);



-- SPI control FSM

fsm : entity work.spi_fsm

port map(

    clk => clk,
    reset => reset,

    start => start,


    spi_cs => spi_cs,
    busy => busy,
    ready => ready

);



-- temporary

data_out <= (others => '0');

spi_mosi <= '0';



end rtl;