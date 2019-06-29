defmodule RankingsWeb.RunnerView do
  use RankingsWeb, :view

  alias Rankings.Runner
  # alias Rankings.Result

  def get_name(%Runner{first_name: first, last_name: last}) do
    first <> " " <> last
  end

  def get_team_name(%Runner{team: team}) do
    team
  end
end
