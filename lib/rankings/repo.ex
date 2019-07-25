defmodule Rankings.Repo do
  use Ecto.Repo,
    otp_app: :rankings,
    adapter: Ecto.Adapters.Postgres
end
