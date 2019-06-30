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
end
