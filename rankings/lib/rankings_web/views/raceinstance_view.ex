defmodule RankingsWeb.RaceInstanceView do
  use RankingsWeb, :view

  alias Rankings.RaceInstance
  alias Rankings.Runner

  def get_date(%RaceInstance{date: date}) do
    date
  end

  def get_name(%Runner{first_name: first, last_name: last}) do
    first <> " " <> last
  end
end
