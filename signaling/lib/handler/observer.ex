defmodule Signaling.Handler.Observer do
  @behaviour Signaling.Handler
  alias Signaling.Message
  alias Signaling.Credential

  def handle(%Message{:type => type, :payload => payload}, state) do
    case type do
      :connect -> connect(payload, state)
      :update -> handle_update(payload, state)
    end
  end

  defp handle_update(payload, state) do
    {:reply, {:text, payload}, state}
  end

  defp connect(%{"name" => name, "password" => password}, state) do
    with {:ok, {pid, receiver}} <- Signaling.Peer.lookup(name),
         {:ok} <- Credential.hash(password) |> Credential.check_password(receiver.password),
         {:ok, updated_peer} <-
           %{state.peer | type: :sender, password: receiver.password}
           |> Signaling.Peer.update() do
      pid
      |> send(
        Message.new!(
          %{"name" => state.peer.name, "password" => password},
          :connect
        )
      )

      {:ok, Map.put(state, :peer, updated_peer)}
    else
      {:error, reason} ->
        {[
           {:text, Message.error(reason)},
           {:close}
         ], state}
    end
  end
end
