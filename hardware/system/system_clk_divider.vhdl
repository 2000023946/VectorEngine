library ieee;
use ieee.std_logic_1164.all;
entity system_clk_divider is
  port (
    clk    : in std_logic;
    resetn : in std_logic;

    system_clk : out std_logic
  );
end system_clk_divider;

architecture rtl of system_clk_divider is

  signal count          : integer   := 0;
  signal system_clk_reg : std_logic := '0';

begin
  process (clk)

  begin

    if rising_edge(clk) then

      if resetn = '0' then

        count          <= 0;
        system_clk_reg <= '0';
      elsif count = 24 then

        count          <= 0;
        system_clk_reg <= not system_clk_reg;
      else

        count <= count + 1;

      end if;

    end if;

  end process;
  system_clk <= system_clk_reg;
end rtl;