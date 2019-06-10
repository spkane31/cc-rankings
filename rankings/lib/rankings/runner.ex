defmodule Rankings.Runner do
  use Ecto.Schema
  import Ecto.Changeset

  alias Rankings

  schema "runners" do
    field :first_name, :string
    field :last_name, :string
    field :year, :string
    belongs_to :team, Rankings.Team
    has_many :results, Rankings.Result
  end

  def changeset(struct, params) do
    struct
    |> cast(params, [:first_name, :last_name])
    |> validate_required([:first_name, :last_name])
    |> validate_length(:first_name, min: 1)
    |> validate_length(:last_name, min: 1)
  end

  alias Rankings.Repo

  def get_runner(id) do
    Repo.get(Rankings.Runner, id)
  end

  def get_runner!(id) do
    Repo.get!(Rankings.Runner, id)
  end

  def list_runners do
    Repo.all(Rankings.Runner)
  end
end
