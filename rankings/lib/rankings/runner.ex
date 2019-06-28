defmodule Rankings.Runner do
  use Ecto.Schema
  import Ecto.Changeset

  alias Rankings
  alias Rankings.Result
  alias Rankings.Repo
  import Ecto.Query

  schema "runners" do
    field :first_name, :string
    field :last_name, :string
    field :year, :string
    belongs_to :team, Rankings.Team
    has_many :instances, Rankings.Result
  end

  def changeset(struct, params) do
    struct
    |> cast(params, [:first_name, :last_name])
    |> validate_required([:first_name, :last_name])
    |> validate_length(:first_name, min: 1)
    |> validate_length(:last_name, min: 1)
  end

  alias Rankings.Repo

  def get_runner(id) do
    Repo.get(Rankings.Runner, id)
  end

  def get_runner!(id) do
    Repo.get!(Rankings.Runner, id)
  end

  def list_runners do
    Repo.all(Rankings.Runner)
  end

  def last_n_runners(n) do
    q = from(r in Rankings.Runner, limit: ^n)
    Repo.all(q)
  end

  def get_team_name(id) do
    r = get_runner(id)
    r = Repo.preload(r, [:team])
    if r.team == nil do
      ""
    else
      r.team.name
    end
  end

  import Ecto.Query
  def get_athlete_results(id) do
    Repo.all(from r in Result, where: r.runner_id == ^id)
    |> Repo.preload([ {:race_instance, :race}])
  end

  def get_results(id) do
    r = get_runner(id)
    r = Repo.preload(r, [:results])
    if r.results == nil do
      ""
    else
      r.results
    end
  end
end
