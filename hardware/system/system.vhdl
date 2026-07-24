library ieee;
use ieee.std_logic_1164.all;
entity system is

  port (

    clk   : in std_logic;
    reset : in std_logic;

    start    : in std_logic;
    leds     : out std_logic_vector(7 downto 0);
    spi_cs   : out std_logic;
    spi_clk  : out std_logic;
    spi_mosi : out std_logic;
    spi_miso : in std_logic
  );

end system;

architecture rtl of system is
  signal spi_start_signal : std_logic;

  signal spi_ready_signal : std_logic;

  signal spi_tx_data_signal :
  std_logic_vector(7 downto 0);

  signal spi_rx_data_signal :
  std_logic_vector(7 downto 0);
  signal device_id_signal :
  std_logic_vector(7 downto 0);

begin
  accel_controller : entity work.accelerometer_controller

    port map
    (

      clk => clk,

      reset => reset,

      start     => start,
      spi_start => spi_start_signal,

      spi_tx_data => spi_tx_data_signal,
      spi_ready   => spi_ready_signal,

      spi_rx_data => spi_rx_data_signal,
      device_id   => device_id_signal,

      done => open

    );

  spi : entity work.spi_reader

    port map
    (

      clk => clk,

      reset    => reset,
      start    => spi_start_signal,
      busy     => open,
      ready    => spi_ready_signal,
      data_out => spi_rx_data_signal,
      spi_cs   => spi_cs,

      spi_clk => spi_clk,

      spi_mosi => spi_mosi,

      spi_miso => spi_miso

    );

  leds <= device_id_signal;
end rtl;