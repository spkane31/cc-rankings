defmodule Rankings.Edge do
  use Ecto.Schema
  import Ecto.{Query, Changeset}

  alias Rankings.{Repo, Race, Runner, Edge, RaceInstance}

  schema "edges" do
    belongs_to :from_race, RaceInstance
    belongs_to :to_race, RaceInstance
    field :difference, :float
    field :gender, :string
    belongs_to :runner, Rankings.Runner

    timestamps()
  end

  def get_runner_edges(id) do
    q = from(e in Rankings.Edge, where: e.runner == ^id, order_by: [desc: :date])
    Repo.all(q)
    |> Repo.preload([{:from_race_id, :race_id}, {:to_race_id, :race_id}])
  end

  def slowest_n_edges(n) do
    q = from(e in Edge, limit: ^n, order_by: [desc: :difference])
    Repo.all(q)
    |> Repo.preload([:runner])
  end

  def get_edge(id) do
    Repo.get(Edge, id)
    |> Repo.preload([:runner, {:from_race, :race}, {:to_race, :race}])
  end

  def get_edges(from_id, to_id) do
    q = from e in Edge, where: e.from_race_id == ^from_id and e.to_race_id == ^to_id, order_by: [desc: :difference]
    Repo.all(q)
    |> Repo.preload([:runner, {:from_race, :race}, {:to_race, :race}])
  end

  def count_edges(from_id, to_id) do
    q = from e in Edge, where: e.from_race_id == ^from_id and e.to_race_id == ^to_id, select: count()
    Repo.all(q)
  end

  def get_avg(from_id, to_id) do
    q = from e in Edge, where: e.from_race_id == ^from_id and e.to_race_id == ^to_id, select: avg(e.difference)
    Repo.one(q)
  end


end
