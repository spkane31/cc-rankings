defmodule RankingsWeb.RaceView do
  use RankingsWeb, :view

  alias Rankings.{Race, RaceInstance}

  def get_name(%Race{name: name}) do
    name
  end

  def time_to_string(time) do
    min = Float.round(time / 60, 0)
    min
  end

  def get_winner(id) do
    r = Race.get_winner(id)
    r.time
  end

  def get_average_time(id) do
    r = RaceInstance.get_average_time(id)
    # r.total / r.count
  end
end
