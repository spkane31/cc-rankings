defmodule Rankings.Result do
  use Ecto.Schema
  # import Ecto.Query

  schema "results" do
    field :distance, :integer
    field :unit, :string
    field :rating, :float
    field :time, :string
    field :scaled_time, :float
    belongs_to :runner, Rankings.Runner
    belongs_to :race_instance, Rankings.RaceInstance
  end

  alias Rankings.Repo

  def get_result(id) do
    Repo.get(Rankings.Result, id) |> Repo.preload([{:runner, :team}])
  end

  def get_result!(id) do
    Repo.get!(Rankings.Result, id)
  end

  def list_results do
    Repo.all(Rankings.Result) |> Repo.preload(:runner)
  end
end
