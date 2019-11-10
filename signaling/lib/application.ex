defmodule Signaling.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  def start(_type, _args) do
    children = [
      Plug.Cowboy.child_spec(
        scheme: :http,
        plug: Signaling.Router,
        options: [
          port: 4001,
          dispatch: dispatch()
        ]
      ),
      Registry.child_spec(
        keys: :unique,
        name: Signaling.Peer,
        partitions: System.schedulers_online()
      )
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: Signaling.Supervisor]
    Supervisor.start_link(children, opts)
  end

  defp dispatch do
    [
      {:_,
       [
         {"/ws/[...]", Signaling.Handler, []},
         {:_, Plug.Cowboy.Handler, {Signaling.Router, []}}
       ]}
    ]
  end
end
