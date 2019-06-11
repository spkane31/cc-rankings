defmodule RankingsWeb.ResultView do
  use RankingsWeb, :view

  alias Rankings.Result

  def get_distance(%Result{distance: distance}) do
    distance
  end

  def get_time(%Result{time: time}) do
    time
  end

  def get_rating(%Result{rating: rating}) do
    rating
  end

end
