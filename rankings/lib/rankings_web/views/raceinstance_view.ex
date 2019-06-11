defmodule RankingsWeb.RaceInstanceView do
  use RankingsWeb, :view

  alias Rankings.RaceInstance

  def get_date(%RaceInstance{date: date}) do
    date
  end
end
