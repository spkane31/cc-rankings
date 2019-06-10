defmodule Rankings.Runner do
  use Ecto.Schema
  import Ecto.Changeset

  schema "runners" do
    field :first_name, :string
    field :last_name, :string
    field :year, :string
    belongs_to :team, Rankings.Team
  end

  def changeset(struct, params) do
    struct
    |> cast(params, [:first_name, :last_name])
    |> validate_required([:first_name, :last_name])
    |> validate_length(:first_name, min: 1)
    |> validate_length(:last_name, min: 1)
  end
end
