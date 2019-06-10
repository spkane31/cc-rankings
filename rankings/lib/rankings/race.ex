defmodule Rankings.Race do
  use Ecto.Schema
  import Ecto.Changeset

  schema "races" do
    field :name, :string
    field :course, :string
    field :distance, :integer
    field :gender, :string
    field :correction, :float
    has_many :race_instances, Rankings.RaceInstance
  end

  def changeset(struct, params) do
    struct
    |> cast(params, [:name])
    |> validate_required([:name])
    |> validate_length(:name, min: 1)
  end
end
