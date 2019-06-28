defmodule Rankings.Team do
  use Ecto.Schema
  import Ecto.Changeset
  import Ecto.Query

  alias Rankings
  alias Rankings.Repo

  schema "teams" do
    field :name, :string, null: false
    field :region, :string
    field :conference, :string
    has_many :runners, Rankings.Runner
  end

  def changeset(struct, params) do
    struct
    |> cast(params, [:name])
    |> validate_required([:name])
    |> validate_length(:name, min: 1)
  end

  alias Rankings.Repo

  def get_team(id) do
    Repo.get(Rankings.Team, id)
  end

  def get_team!(id) do
    Repo.get!(Rankings.Team, id)
  end

  def get_last_n(n) do
    q = from(r in Rankings.Team, limit: ^n)
    Repo.all(q)

  end

  def list_teams do
    Repo.all(Rankings.Team)
  end
end
