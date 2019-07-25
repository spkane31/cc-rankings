defmodule RankingsWeb.SearchView do
  use RankingsWeb, :view

  alias Rankings.{Runner, Team, Race}

  def get_name(%Runner{first_name: first, last_name: last}) do
    first <> " " <> last
  end

  def get_name(%Race{name: name}) do
    name
  end

  def get_name(%Team{name: name}) do
    name
  end

  def get_name(%Race{name: name}) do
    name
  end
end
