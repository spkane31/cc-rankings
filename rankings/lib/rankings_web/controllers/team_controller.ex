defmodule RankingsWeb.TeamController do
  use RankingsWeb, :controller

  alias Rankings.Team
  alias Rankings.Repo

  def index(conn, _params) do
    teams = Team.list_teams()
    render(conn, "index.html", teams: teams)
  end

  def show(conn, %{"id" => id}) do
    team = Team.get_team(id) |> Repo.preload(:runners)
    runners = team.runners
    render(conn, "show.html", team: team, runners: runners)
  end

end
