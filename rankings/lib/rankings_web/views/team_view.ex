defmodule RankingsWeb.TeamView do
  use RankingsWeb, :view

  alias Rankings.Team
  alias Rankings.Runner

  def get_name(%Team{name: name}) do
    name
  end

  def get_runner_name(%Runner{first_name: first, last_name: last}) do
    first <> " " <> last
  end
end
