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
    field :gender, :string
    field :elo_rating, :float
    field :speed_rating, :float
    belongs_to :team, Rankings.Team
    has_many :instances, Rankings.Result
    has_many :edges, Rankings.Edge
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
    q = from(r in Rankings.Runner, limit: ^n, order_by: :speed_rating)
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
    Repo.all(from r in Result, where: r.runner_id == ^id, order_by: [desc: :date])
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

  alias Rankings.Runner

  def list_runners(params) do
    first = get_in(params, ["first"])
    last = get_in(params, ["last"])
    if first == nil and last == nil do
      nil
    else
      Runner |> Runner.search(first, last) |> Repo.all() |> Repo.preload(:team)
    end
  end

  def search(query, first, last) do
    wildcard_first = "%#{first}"
    wildcard_last = "%#{last}"

    from r in query,
    where: ilike(r.first_name, ^wildcard_first) and ilike(r.last_name, ^wildcard_last)
  end
end
