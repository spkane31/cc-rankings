defmodule RankingsWeb.TeamView do
  use RankingsWeb, :view

  alias Rankings.Team

  def get_name(%Team{name: name}) do
    name
  end
end
