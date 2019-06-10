defmodule RankingsWeb.TeamController do
  use RankingsWeb, :controller

  alias Rankings.Team

  def index(conn, _params) do
    teams = Team.list_teams()
    render(conn, "index.html", teams: teams)
  end

  def show(conn, %{"id" => id}) do
    team = Team.get_team(id)
    render(conn, "show.html", team: team)
  end

end
