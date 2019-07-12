defmodule RankingsWeb.RaceController do
  use RankingsWeb, :controller

  alias Rankings.Race
  alias Rankings.Repo

  def index(conn, _params) do
    races = Race.get_n_races(25) #Repo.all(Race)
    render(conn, "index.html", races: races)
  end

  def show(conn, %{"id" => id}) do
    race = Race.get_race(id) |> Repo.preload([{:instances, :instance_results}])
    race_instances = race.instances
    render(conn, "show.html", race: race, race_instances: race_instances)
  end
end
