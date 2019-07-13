defmodule Rankings.Race do
  use Ecto.Schema
  import Ecto.Changeset
  import Ecto.Query

  alias Rankings

  schema "races" do
    field :name, :string
    field :course, :string
    field :distance, :integer
    field :gender, :string
    field :is_base, :boolean
    field :average, :float
    field :std_dev, :float
    field :correction_avg, :float
    field :correction_graph, :float
    has_many :instances, Rankings.RaceInstance
  end

  def changeset(struct, params) do
    struct
    |> cast(params, [:name])
    |> validate_required([:name])
    |> validate_length(:name, min: 1)
  end

  alias Rankings.Repo

  def get_race(id) do
    Repo.get(Rankings.Race, id)
  end

  def get_race!(id) do
    Repo.get!(Rankings.Race, id)
  end

  def list_races do
    Repo.all(Rankings.Race)
  end

  def get_n_races(n) do
    q = from(r in Rankings.Race, limit: ^n)
    Repo.all(q)
  end

  alias Rankings.Race
  def list_races(params) do
    race = get_in(params, ["race"])
    Race
    |> Race.search(race) |> Repo.all()
  end

  def search(query, name) do
    wildcard = "%#{name}%"
    from r in query,
    where: ilike(r.name, ^wildcard)
  end

  def get_runner_count(id) do
    r = get_race(id) |> Repo.preload([{:instances, :instance_results}])
    count = 0
    for instance <- @r do
      count = count + length(instance.results)
    end
    count
  end

  def get_winner(id) do
    r = from(r in Rankings.Result,
      where: r.race_instance_id == ^id and
        r.place == 1,
      limit: 1)
    Repo.one(r)
  end
end
