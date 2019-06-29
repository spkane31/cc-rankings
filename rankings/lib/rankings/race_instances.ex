defmodule Rankings.RaceInstance do
  use Ecto.Schema
  import Ecto.Query

  alias Rankings.Result
  alias Rankings

  schema "race_instances" do
    field :date, :string, null: false
    field :average, :float, default: 0
    field :std_dev, :float, default: 0
    belongs_to :race, Rankings.Race
    has_many :instance_results, Rankings.Result
  end

  alias Rankings.Repo

  def get_instance(id) do
    Repo.get(Rankings.RaceInstance, id)
  end

  def get_instance!(id) do
    Repo.get!(Rankings.RaceInstance, id)
  end

  def get_n_instances(n) do
    q = from(i in Rankings.RaceInstance, limit: ^n)
    Repo.all(q)
  end

  def list_race_instances do
    Repo.all(Rankings.RaceInstance)
  end
  def get_results(id) do
    query = from(r in Result, where: r.race_instance_id == ^id, preload: [:runner])
    Repo.all(query)
  end

end
