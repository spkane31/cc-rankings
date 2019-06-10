defmodule Rankings.Result do
  use Ecto.Schema

  schema "results" do
    field :distance, :integer
    # field :unit, :string # THIS NEEDS TO BE ADDED
    field :rating, :float
    field :time, :string
    belongs_to :race_instance, Rankings.RaceInstance
  end
end
