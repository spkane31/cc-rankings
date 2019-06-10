defmodule Rankings.Result do
  use Ecto.Schema

  schema "races" do
    field :distance, :integer
    field :unit, :string
    field :rating, :float
    field :time, :string
    belongs_to :race_instance, Rankings.RaceInstance
  end
end
