library ieee;
use ieee.std_logic_1164.all;
entity spi_fsm is

  port (

    clk   : in std_logic;
    reset : in std_logic;

    start : in std_logic;
    -- SPI control
    spi_cs : out std_logic;

    busy  : out std_logic;
    ready : out std_logic;
    -- Shift register control
    load         : out std_logic;
    shift_enable : out std_logic;

    done : in std_logic

  );

end spi_fsm;
architecture rtl of spi_fsm is
  type state_type is (
    IDLE,
    ASSERT_CS,
    SEND_COMMAND,
    SEND_ADDRESS,
    RECEIVE_DATA,
    FINISH
  );
  signal state : state_type := IDLE;

begin

  -- State transition logic
  process (clk)

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

            if done = '1' then

              state <= SEND_ADDRESS;

            end if;

          when SEND_ADDRESS =>

            if done = '1' then

              state <= RECEIVE_DATA;

            end if;

          when RECEIVE_DATA =>

            if done = '1' then

              state <= FINISH;

            end if;

          when FINISH =>

            state <= IDLE;

        end case;
      end if;
    end if;

  end process;

  -- Output logic

  process (state)

  begin
    spi_cs <= '1';
    busy   <= '0';
    ready  <= '0';

    case state is

      when IDLE =>

        spi_cs <= '1';
        busy   <= '0';

      when ASSERT_CS =>

        spi_cs <= '0';
        busy   <= '1';

      when SEND_COMMAND =>

        spi_cs <= '0';
        busy   <= '1';

      when SEND_ADDRESS =>

        spi_cs <= '0';
        busy   <= '1';

      when RECEIVE_DATA =>

        spi_cs <= '0';
        busy   <= '1';

      when FINISH =>

        spi_cs <= '1';
        ready  <= '1';

    end case;
  end process;

end rtl;