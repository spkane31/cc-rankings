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
      (1900 - t) / (d / 1.609)
    else
      0
    end
  end

end
