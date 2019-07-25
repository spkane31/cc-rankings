defmodule RankingsWeb.RunnerView do
  use RankingsWeb, :view
  import Float
  import Ecto.Query
  import Statistics

  alias Rankings.{Runner, Result, Edge, Repo, RaceInstance}
  alias Date

  def get_name(%Runner{first_name: first, last_name: last}) do
    first <> " " <> last
  end

  def get_team_name(%Runner{team: team}) do
    team
  end

  def get_rating(%Result{scaled_time: t, distance: d, race_instance: r}, %Runner{gender: g}) do
    if g == "MALE" do
      if d == 10000 or d == 8000 do
        rating = (1900 - t - r.race.correction_graph) / (8000.0 / 1609.0)
        Float.round(rating, 3)
      end
    else
      if d == 5000 or d == 6000 do
        rating = (1350 - t - r.race.correction_graph) / (5000 / 1609)
        Float.round(rating, 3)
      end
    end
  end

  def stringify(date) do
    Date.to_string(date)
  end

  def average_difference(from_id, to_id) do
    round(Edge.get_avg(from_id, to_id), 2)
  end

  def get_time(race_id, runner_id) do
    q = from r in Result, where: r.runner_id == ^runner_id and r.race_instance_id == ^race_id
    r = Repo.one(q)
    r.time
  end

  def get_scaled_time(race_id, runner_id) do
    q = from r in Result, where: r.runner_id == ^runner_id and r.race_instance_id == ^race_id
    r = Repo.one(q)
    RaceInstance.time_to_string(r.scaled_time)
  end

  def percentile(from_id, to_id, diff) do
    std_dev = Edge.get_std_dev(from_id, to_id)
    avg = Edge.get_avg(from_id, to_id)

    c = 1 / (std_dev * :math.sqrt(2 * 3.14159))

    e = :math.pow((diff - avg), 2) / (2 * :math.pow(std_dev, 2))

    round(:math.exp(-e) / c, 2)

  end

end
