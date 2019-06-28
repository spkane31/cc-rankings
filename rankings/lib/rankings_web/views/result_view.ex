defmodule RankingsWeb.ResultView do
  use RankingsWeb, :view

  alias Rankings.Result

  def get_distance(%Result{distance: distance}) do
    distance
  end

  def get_time(%Result{time: time}) do
    time
  end

  def get_rating(%Result{time_float: t, distance: d}) do
    if t != nil do
      if d == 6000 do
        (1300 - (t/1.21)) / (d / 1609)
      else if d == 5000 do
        (1300 - t) / (d / 1609)
      end
      end
    else
      0
    end
  end

end
