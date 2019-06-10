defmodule RankingsWeb.RaceView do
  use RankingsWeb, :view

  alias Rankings.Race

  def get_name(%Race{name: name}) do
    name
  end
end
