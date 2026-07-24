library ieee;
use ieee.std_logic_1164.all;
entity spi_reader is

  port (

    clk      : in std_logic;
    reset    : in std_logic;
    start    : in std_logic;
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
  -- Internal SPI clock
  ------------------------------------------------

  signal spi_clk_internal : std_logic;

  ------------------------------------------------
  -- FSM <-> Shift Register signals
  ------------------------------------------------

  signal load_signal : std_logic;

  signal shift_enable_signal : std_logic;

  signal done_signal : std_logic;

  ------------------------------------------------
  -- Data signals
  ------------------------------------------------

  signal tx_data_signal : std_logic_vector(7 downto 0);

  signal rx_data_signal : std_logic_vector(7 downto 0);

begin

  ------------------------------------------------
  -- SPI CLOCK DIVIDER
  ------------------------------------------------

  clock_divider : entity work.spi_clock_divider

    port map
    (

      clk => clk,

      reset => reset,

      spi_clk => spi_clk_internal

    );

  -- Send SPI clock to outside sensor

  spi_clk <= spi_clk_internal;

  ------------------------------------------------
  -- SPI FSM
  ------------------------------------------------

  fsm : entity work.spi_fsm

    port map
    (

      -- FSM now runs on SPI clock

      clk    => spi_clk_internal,
      reset  => reset,
      start  => start,
      spi_cs => spi_cs,
      busy   => busy,

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
      -- Shift register uses SPI clock

      clk          => spi_clk_internal,
      reset        => reset,
      load         => load_signal,
      shift_enable => shift_enable_signal,
      tx_data      => tx_data_signal,
      spi_miso     => spi_miso,
      spi_mosi     => spi_mosi,
      rx_data      => rx_data_signal,
      done         => done_signal

    );

  ------------------------------------------------
  -- Temporary transmit byte
  ------------------------------------------------

  -- Later replaced with:
  -- Accelerometer command/address

  tx_data_signal <= "10110010";

  ------------------------------------------------
  -- Sensor received data
  ------------------------------------------------

  data_out <= rx_data_signal;

end rtl;