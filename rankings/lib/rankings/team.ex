defmodule Rankings.Team do
  use Ecto.Schema

  schema "teams" do
    field :name, :string, null: false
    field :region, :string
    field :conference, :string
    has_many :runners, Rankings.Runner
  end
end
