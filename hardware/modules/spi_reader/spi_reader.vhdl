library ieee;
use ieee.std_logic_1164.all;

entity spi_reader is
  port (
    clk      : in std_logic;
    reset    : in std_logic;
    start    : in std_logic;
    tx_data  : in std_logic_vector(7 downto 0);
    busy     : out std_logic;
    ready    : out std_logic;
    data_out : out std_logic_vector(7 downto 0);
    spi_cs   : out std_logic;
    spi_clk  : out std_logic;
    spi_mosi : out std_logic;
    spi_miso : in std_logic
  );
end spi_reader;

architecture rtl of spi_reader is

  ------------------------------------------------
  -- FSM <-> Shift Register signals
  ------------------------------------------------
  signal load_signal         : std_logic;
  signal shift_enable_signal : std_logic;
  signal done_signal         : std_logic;

  ------------------------------------------------
  -- Data signals
  ------------------------------------------------
  signal rx_data_signal : std_logic_vector(7 downto 0);

begin

  ------------------------------------------------
  -- Gated SPI Clock (Mode 3)
  ------------------------------------------------
  -- The clock idles HIGH ('1').
  -- It only toggles when shift_enable is active.
  -- By using 'not clk', the sensor captures data precisely 
  -- half a clock cycle after our shift register updates it.

  spi_clk <= not clk when shift_enable_signal = '1' else
    '1';

  ------------------------------------------------
  -- SPI FSM
  ------------------------------------------------
  fsm : entity work.spi_fsm
    port map
    (
      clk          => clk, -- Using the incoming slow clock
      reset        => reset,
      start        => start,
      spi_cs       => spi_cs,
      busy         => busy,
      ready        => ready,
      load         => load_signal,
      shift_enable => shift_enable_signal,
      done         => done_signal
    );

  ------------------------------------------------
  -- SPI SHIFT REGISTER
  ------------------------------------------------
  shift_register : entity work.spi_shift_register
    port map
    (
      clk          => clk, -- Using the incoming slow clock
      reset        => reset,
      load         => load_signal,
      shift_enable => shift_enable_signal,
      tx_data      => tx_data,
      spi_miso     => spi_miso,
      spi_mosi     => spi_mosi,
      rx_data      => rx_data_signal,
      done         => done_signal
    );

  ------------------------------------------------
  -- Sensor received data
  ------------------------------------------------
  data_out <= rx_data_signal;

end rtl;