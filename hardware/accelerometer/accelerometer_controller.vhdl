library ieee;
use ieee.std_logic_1164.all;


entity accelerometer_controller is

    port(

        clk   : in std_logic;
        reset : in std_logic;


        -- SPI Reader interface

        spi_start : out std_logic;

        spi_tx_data : out std_logic_vector(7 downto 0);

        spi_ready : in std_logic;

        spi_rx_data : in std_logic_vector(7 downto 0);



        -- Sensor outputs

        x_data : out std_logic_vector(15 downto 0);
        y_data : out std_logic_vector(15 downto 0);
        z_data : out std_logic_vector(15 downto 0)

    );

end accelerometer_controller;