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

  type state_type is (
    IDLE,
    ASSERT_CS,
    SEND_COMMAND,
    SEND_ADDRESS,
    RECEIVE_DATA,
    FINISH
  );

  signal state : state_type := IDLE;

  signal count       : integer := 0;
  signal spi_clk_reg : std_logic := '0';

begin


  ----------------------------------------------------
  -- SPI Clock Divider
  ----------------------------------------------------
  process(clk)
  begin

    if rising_edge(clk) then

      if reset = '1' then

        count       <= 0;
        spi_clk_reg <= '0';

      elsif count = 24 then

        spi_clk_reg <= not spi_clk_reg;
        count       <= 0;

      else

        count <= count + 1;

      end if;

    end if;

  end process;


  spi_clk <= spi_clk_reg;



  ----------------------------------------------------
  -- FSM State Register
  ----------------------------------------------------
  process(clk)
  begin

    if rising_edge(clk) then

      if reset = '1' then

        state <= IDLE;

      else

        case state is

          when IDLE =>

            if start = '1' then
              state <= ASSERT_CS;
            end if;


          when ASSERT_CS =>

            state <= SEND_COMMAND;


          when SEND_COMMAND =>

            state <= SEND_ADDRESS;


          when SEND_ADDRESS =>

            state <= RECEIVE_DATA;


          when RECEIVE_DATA =>

            state <= FINISH;


          when FINISH =>

            state <= IDLE;


        end case;

      end if;

    end if;

  end process;



  ----------------------------------------------------
  -- FSM Output Logic
  ----------------------------------------------------
  process(state)
  begin

    -- Default values
    spi_cs   <= '1';
    busy     <= '0';
    ready    <= '0';
    spi_mosi <= '0';


    case state is


      ------------------------------------------------
      -- Waiting for transaction
      ------------------------------------------------
      when IDLE =>

        spi_cs <= '1';
        busy   <= '0';
        ready  <= '0';



      ------------------------------------------------
      -- Select the sensor
      ------------------------------------------------
      when ASSERT_CS =>

        spi_cs <= '0';
        busy   <= '1';



      ------------------------------------------------
      -- Send READ command
      ------------------------------------------------
      when SEND_COMMAND =>

        spi_cs <= '0';
        busy   <= '1';



      ------------------------------------------------
      -- Send register address
      ------------------------------------------------
      when SEND_ADDRESS =>

        spi_cs <= '0';
        busy   <= '1';



      ------------------------------------------------
      -- Receive sensor data
      ------------------------------------------------
      when RECEIVE_DATA =>

        spi_cs <= '0';
        busy   <= '1';



      ------------------------------------------------
      -- Transaction complete
      ------------------------------------------------
      when FINISH =>

        spi_cs <= '1';
        busy   <= '0';
        ready  <= '1';



    end case;

  end process;


  ----------------------------------------------------
  -- Placeholder data output
  ----------------------------------------------------
  data_out <= (others => '0');


end rtl;