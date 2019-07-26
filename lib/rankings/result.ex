defmodule Rankings.Result do
  use Ecto.Schema
  import Ecto.Query

  schema "results" do
    field :distance, :integer
    field :unit, :string
    field :rating, :float
    field :time, :string
    field :scaled_time, :float
    field :time_float, :float
    field :date, :date
    field :gender, :string
    field :place, :integer
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

  def last_n_results(n) do
    q = from(r in Rankings.Result, limit: ^n, order_by: [asc: :scaled_time])
    Repo.all(q)
  end

  def get_top_25_men do
    mens = from(r in Rankings.Result, where: r.gender == "MALE", limit: 25, order_by: [desc: :rating])
    Repo.all(mens) |> Repo.preload([:runner, {:race_instance, :race}])
  end

  def get_top_25_women do
    womens = from(r in Rankings.Result, where: r.gender == "FEMALE", limit: 25, order_by: [desc: :rating])
    Repo.all(womens) |> Repo.preload([:runner, {:race_instance, :race}])
  end
end
