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

  def get_rating(%Result{time_float: t, distance: d, race_instance: r}, %Runner{gender: g}) do
    if g == "MENS" do
      if d == 10000 do
        rating = (1900 - (t / 1.268) - r.race.correction_avg) / (d / 1609)
        Float.round(rating, 3)
      else
        (1900 - t - r.race.correction_avg) / (d / 1609)
      end
    else
      if d == 5000 do
        rating = (1300 - t - r.race.correction_avg) / (d / 1609)
      else
        rating = (1300 - (t/1.213) - r.race.correction_avg) / (d / 1609)
        Float.round(rating, 3)
      end

    end
  end

  def stringify(date) do
    Date.to_string(date)
  end

end
