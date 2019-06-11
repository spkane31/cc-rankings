defmodule Rankings.RaceInstance do
  use Ecto.Schema
  import Ecto.Changeset

  alias Rankings

  schema "race_instances" do
    field :date, :string, null: false
    belongs_to :race, Rankings.Race
    has_many :results, Rankings.Result
  end

  alias Rankings.Repo

  def get_instance(id) do
    Repo.get(Rankings.RaceInstance, id)
  end

  def get_instance!(id) do
    Repo.get!(Rankings.RaceInstance, id)
  end

  def list_race_instances do
    Repo.all(Rankings.RaceInstance)
  end

  def get_results(id) do
    instance = Repo.get(Rankings.RaceInstance, id) |> Repo.preload(:results)
    results = instance.results
    results = Repo.preload(results, :runner)
    results
  end

end
