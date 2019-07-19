defmodule Rankings.Edge do
  use Ecto.Schema
  import Ecto.{Query, Changeset}

  alias Rankings.{Repo, Race, Runner}

  schema "edges" do
    field :from_race_id, :integer
    field :to_race_id, :integer
    field :difference, :float
    field :gender, :string
    belongs_to :runner, Rankings.Runner
  end

  def get_runner_edges(id) do
    q = from(e in Rankings.Edge, where: e.runner == ^id, order_by: [desc: :date])
    Repo.all(q)
    |> Repo.preload([{:from_race_id, :race_id}, {:to_race_id, :race_id}])
  end

end
