defmodule Rankings.RaceInstance do
  use Ecto.Schema

  schema "race_instances" do
    field :date, :string, null: false
    belongs_to :race, Rankings.Race
  end

end
