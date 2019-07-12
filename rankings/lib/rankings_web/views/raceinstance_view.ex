defmodule RankingsWeb.RaceInstanceView do
  use RankingsWeb, :view

  alias Rankings.RaceInstance
  alias Rankings.Runner
  alias Rankings.Team

  import Float

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

  def get_runner_team(id) do
    runner = Repo.get(Runner, id) |> Repo.preload(:team)
    team = Repo.get(Team, runner.team_id)
    team.name
  end

  alias Rankings.Result
  def get_runner_id(%Result{runner_id: id}) do
    id
  end

  def get_rating(%Result{scaled_time: t, gender: g, race_instance: r, distance: d}) do
    if g == "MALE" do
      if d == 10000 or d == 8000 do
        rating = (1900 - t - r.race.correction_graph) / (8000.0 / 1609.0)
        Float.round(rating, 3)
      end
    else
      if d == 5000 or d == 6000 do
        rating = (1300 - t - r.race.correction_graph) / (5000.0 / 1609.0)
        Float.round(rating, 3)
      end
    end
  end

end
