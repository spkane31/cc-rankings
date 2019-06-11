defmodule Rankings.Result do
  use Ecto.Schema

  schema "results" do
    field :distance, :integer
    # field :unit, :string # THIS NEEDS TO BE ADDED
    field :rating, :float
    field :time, :string
    belongs_to :runner, Rankings.Runner
    belongs_to :race_instance, Rankings.RaceInstance
  end

  alias Rankings.Repo

  def get_result(id) do
    Repo.get(Rankings.Result, id) |> Repo.preload(:runner)
  end

  def get_result!(id) do
    Repo.get!(Rankings.Result, id)
  end

  def list_results do
    Repo.all(Rankings.Result) |> Repo.preload(:runner)
  end
end
