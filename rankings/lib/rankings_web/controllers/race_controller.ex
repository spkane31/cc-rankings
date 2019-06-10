defmodule RankingsWeb.RaceController do
  use RankingsWeb, :controller

  alias Rankings.Race
  alias Rankings.Repo

  def index(conn, _params) do
    races = Repo.all(Race)
    render(conn, "index.html", races: races)
  end

  def show(conn, %{"id" => id}) do
    race = Race.get_race(id)
    render(conn, "show.html", race: race)
  end
end
