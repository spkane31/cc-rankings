defmodule Rankings.Result do
  use Ecto.Schema

  schema "races" do
    field :distance, :integer
    field :unit, :string
    field :rating, :float
    field :time, :string
  end
end
