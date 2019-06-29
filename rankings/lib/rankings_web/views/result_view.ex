defmodule RankingsWeb.ResultView do
  use RankingsWeb, :view

  alias Rankings.Result

  def get_distance(%Result{distance: distance}) do
    distance
  end

  def get_time(%Result{time: time}) do
    time
  end

  def get_rating(%Result{time_float: t, distance: d, race_instance: i}) do
    correction = 0
    if i.race.correction != nil do
      correction = i.race.correction
    end
    if t != nil do
      if d == 6000 && i.race.gender == "WOMENS" do
        (1300 - (t/1.21) - correction) / (d / 1609)
      else if d == 5000 && i.race.gender == "WOMENS" do
        (1300 - t - correction) / (d / 1609)
      else if d == 8000 && i.race.gender == "MENS" do
        (1900 - t - correction) / (d / 1609)
      else if d == 10000 && i.race.gender == "MENS" do
        (1900 - (t/1.268) - correction) / (d / 1609)
      end
      end
      end
      end
    else
      0
    end
  end

end
