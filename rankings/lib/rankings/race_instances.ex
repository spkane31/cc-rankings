defmodule Rankings.RaceInstance do
  use Ecto.Schema
  import Ecto.Query
  import Integer

  alias Rankings.{Result, Repo}
  alias Rankings

  schema "race_instances" do
    field :date, :date, null: false
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
    query = from(r in Result, where: r.race_instance_id == ^id, preload: [:runner], order_by: [desc: :rating])
    Repo.all(query)
  end

  def time_to_string(time) do
    minutes = round(time / 60)
    m = Integer.to_string(minutes)
    seconds = round(time - 60 * minutes)
    s = Integer.to_string(round(time - minutes * 60))
    if seconds < 10 do
      m <> ":0" <> s
    else
      m <> ":" <> s
    end

  end

  def get_average_time(id) do
    r = from(r in Result, where: r.race_instance_id == ^id)
    results = Repo.all(r)
    total = Enum.map(results, fn r -> r.time_float end) |> Enum.sum()
    avg = total / length(results)
    time_to_string(avg)
  end

end
