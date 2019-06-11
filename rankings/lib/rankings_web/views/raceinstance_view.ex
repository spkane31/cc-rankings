defmodule RankingsWeb.RaceInstanceView do
  use RankingsWeb, :view

  alias Rankings.RaceInstance
  alias Rankings.Runner

  def get_date(%RaceInstance{date: date}) do
    date
  end

  def get_name(%Runner{first_name: first, last_name: last}) do
    first <> " " <> last
  end

  def get_team(%Runner{team: team}) do
    team
  end

  alias Rankings.Repo
  def get_runner_name(id) do
    runner = Repo.get(Runner, id)
    runner.first_name <> " " <> runner.last_name
  end

end
