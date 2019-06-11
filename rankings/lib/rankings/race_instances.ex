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

  alias Rankings.Result
  import Ecto.Query
  def get_results(id) do
    query = from(r in Result, where: r.race_instance_id == ^id, preload: [:runner])
    Repo.all(query)
  end

end
