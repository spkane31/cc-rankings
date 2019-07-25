defmodule RankingsWeb.TeamView do
  use RankingsWeb, :view
  import Ecto.Query

  alias Rankings.{Team, Runner, Result, Repo}
  # alias Rankings.Runner

  def get_name(%Team{name: name}) do
    name
  end

  def get_runner_name(%Runner{first_name: first, last_name: last}) do
    first <> " " <> last
  end

  def best_performance(id) do
    q = from r in Result, where: r.runner_id == ^id, select: max(r.rating)
    r = Repo.one(q)
    Float.round(r, 2)
  end

end
