library ieee;
use ieee.std_logic_1164.all;
use ieee.numeric_std.all;

entity system is
  port (
    clk       : in  std_logic; -- Standard 50MHz board clock
    reset_btn : in  std_logic; -- Physical button 
    miso      : in  std_logic;
    mosi      : out std_logic;
    sclk      : out std_logic;
    cs        : out std_logic;
    leds      : out std_logic_vector(7 downto 0)
  );
end system;

architecture rtl of system is
  -- Unified internal reset (inverts the physical active-low button)
  signal rst : std_logic;

  -- Clock divider to get a safe ~780kHz SPI enable pulse
  signal clk_div  : unsigned(5 downto 0) := (others => '0');
  signal spi_tick : std_logic := '0';

  -- SPI State Machine
  type state_t is (IDLE, TRANSFER, DONE);
  signal state : state_t := IDLE;

  -- 16 bits * 2 phases (falling/rising edges) = 32 phases total
  signal bit_counter : integer range 0 to 31 := 0; 
  
  -- x"8000" = 10000000 (Read Reg 0x00) followed by 00000000 (Dummy byte to keep clock running)
  signal tx_reg : std_logic_vector(15 downto 0) := x"8000"; 
  signal rx_reg : std_logic_vector(15 downto 0) := (others => '0');

  signal sclk_internal : std_logic := '1'; -- SPI Mode 3 idles high

begin

  -- 1. Create a true Active-High internal reset
  rst <= not reset_btn; 

  -- 2. Generate a 1-cycle enable pulse to drive the SPI state machine safely
  process(clk)
  begin
    if rising_edge(clk) then
      clk_div <= clk_div + 1;
      if clk_div = 0 then
        spi_tick <= '1';
      else
        spi_tick <= '0';
      end if;
    end if;
  end process;

  -- 3. SPI State Machine (Advances only when spi_tick is high)
  process(clk)
  begin
    if rising_edge(clk) then
      if rst = '1' then
        state         <= IDLE;
        sclk_internal <= '1';
        cs            <= '1';
        mosi          <= '0';
        rx_reg        <= (others => '0');
        tx_reg        <= x"8000";
        bit_counter   <= 0;
        leds          <= (others => '0');
        
      elsif spi_tick = '1' then
        case state is
          when IDLE =>
            cs            <= '1';
            sclk_internal <= '1';
            bit_counter   <= 0;
            tx_reg        <= x"8000"; 
            state         <= TRANSFER; -- Auto-start reading immediately

          when TRANSFER =>
            cs <= '0'; -- Pull CS low to begin transaction

            -- EVEN phases (0, 2, 4...) -> Falling edge of SCLK -> Shift MOSI out
            if (bit_counter mod 2) = 0 then
              sclk_internal <= '0';
              mosi          <= tx_reg(15 - (bit_counter / 2));
              bit_counter   <= bit_counter + 1;

            -- ODD phases (1, 3, 5...) -> Rising edge of SCLK -> Sample MISO in
            else
              sclk_internal <= '1';
              rx_reg(15 - (bit_counter / 2)) <= miso;

              if bit_counter = 31 then
                state <= DONE; -- 16 bits finished
              else
                bit_counter <= bit_counter + 1;
              end if;
            end if;

          when DONE =>
            cs            <= '1';
            sclk_internal <= '1';
            
            -- The ADXL345 sends the Device ID during the second byte (lower 8 bits)
            leds <= rx_reg(7 downto 0); 
            
            -- Stay in DONE forever. Pressing the physical reset button will restart the process.
        end case;
      end if;
    end if;
  end process;

  -- Drive the physical SCLK pin
  sclk <= sclk_internal;

end rtl;