defmodule RankingsWeb.RunnerView do
  use RankingsWeb, :view
  import Float

  alias Rankings.Runner
  alias Rankings.Result
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

end
