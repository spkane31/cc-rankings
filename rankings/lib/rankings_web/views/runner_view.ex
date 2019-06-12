defmodule RankingsWeb.RunnerView do
  use RankingsWeb, :view

  alias Rankings.Runner

  def get_first_name(%Runner{first_name: first}) do
    first
    |> String.split(" ")
    |> Enum.at(0)
  end

  def get_last_name(%Runner{last_name: last}) do
    last
  end

  def get_name(%Runner{first_name: first, last_name: last}) do
    first <> " " <> last
  end

  def get_team_name(%Runner{team: team}) do
    team
  end

  # def get_results(%Runner{results: results}) do
  #   results
  # end
end
