defmodule Rankings.Race do
  use Ecto.Schema
  import Ecto.Changeset
  import Ecto.Query

  alias Rankings

  schema "races" do
    field :name, :string
    field :course, :string
    field :distance, :integer
    field :gender, :string
    field :correction, :float
    field :is_base, :boolean
    has_many :instances, Rankings.RaceInstance
  end

  def changeset(struct, params) do
    struct
    |> cast(params, [:name])
    |> validate_required([:name])
    |> validate_length(:name, min: 1)
  end

  alias Rankings.Repo

  def get_race(id) do
    Repo.get(Rankings.Race, id)
  end

  def get_race!(id) do
    Repo.get!(Rankings.Race, id)
  end

  def list_races do
    Repo.all(Rankings.Race)
  end

  def get_n_races(n) do
    q = from(r in Rankings.Team, limit: ^n)
    Repo.all(q)
  end
end
