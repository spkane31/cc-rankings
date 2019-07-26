defmodule Rankings.RaceInstance do
  use Ecto.Schema
  import Ecto.Query

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
    minutes = trunc(time / 60)
    m = Integer.to_string(minutes)
    seconds = trunc(time - (60 * minutes))
    s = Integer.to_string(seconds)
    milli = trunc(10 * (time - trunc(time)))
    ms = Integer.to_string(milli)

    if seconds < 10 do
      m <> ":0" <> s <> "." <> ms
    else
      m <> ":" <> s <> "." <> ms
    end

  end

  def get_average_time(id) do
    q = from(r in Result, select: avg(r.time_float), where: r.race_instance_id == ^id)
    Repo.one(q)
    |> time_to_string()
  end

end
