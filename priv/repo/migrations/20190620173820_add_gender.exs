defmodule Rankings.Repo.Migrations.AddGender do
  use Ecto.Migration

  def change do
    alter table(:runners) do
      add :gender, :string
    end
  end
end
