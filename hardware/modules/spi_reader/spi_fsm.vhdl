library ieee;
use ieee.std_logic_1164.all;

entity spi_fsm is
  port (
    clk          : in std_logic;
    resetn        : in std_logic;
    start        : in std_logic;
    spi_cs       : out std_logic;
    busy         : out std_logic;
    ready        : out std_logic;
    load         : out std_logic;
    shift_enable : out std_logic;
    done         : in std_logic
  );
end spi_fsm;

architecture rtl of spi_fsm is
  type state_type is (
    IDLE,
    ASSERT_CS,
    LOAD_ADDRESS,
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
      if resetn = '0' then
        state <= IDLE;
      else
        case state is
          when IDLE =>
            if start = '1' then
              state <= ASSERT_CS;
            end if;
            
          when ASSERT_CS =>
            state <= LOAD_ADDRESS;

          when LOAD_ADDRESS =>
            state <= SEND_ADDRESS;

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
    -- Default values
    spi_cs       <= '1';
    busy         <= '0';
    ready        <= '0';
    load         <= '0';
    shift_enable <= '0';

    case state is
      when IDLE =>
        spi_cs <= '1';
        busy   <= '0';

      when ASSERT_CS =>
        spi_cs <= '0';
        busy   <= '1';

      -- Load address/control byte into shift register
      when LOAD_ADDRESS =>
        spi_cs <= '0';
        busy   <= '1';
        load   <= '1';

      -- Shift address bits out on MOSI
      when SEND_ADDRESS =>
        spi_cs       <= '0';
        busy         <= '1';
        shift_enable <= '1';

      -- Keep shifting to receive sensor data on MISO
      when RECEIVE_DATA =>
        spi_cs       <= '0';
        busy         <= '1';
        shift_enable <= '1';

      when FINISH =>
        spi_cs <= '1';
        ready  <= '1';
        
    end case;
  end process;

end rtl;