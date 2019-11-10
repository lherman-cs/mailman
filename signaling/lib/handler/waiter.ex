defmodule Signaling.Handler.Waiter do
  @behaviour Signaling.Handler
  alias Signaling.Message

  def handle(%Message{:type => type, :payload => payload}, state) do
    case type do
      :create_session -> create_session(payload, state)
      :observe -> observe(state)
      _ -> {:ok, state}
    end
  end

  defp create_session(%{"password" => password}, state) do
    resp =
      %{state.peer | type: :receiver, password: password}
      |> Signaling.Peer.update()

    case resp do
      {:ok, updated_peer} ->
        Signaling.Peer.broadcast_update()

        {[{:text, Message.new!(nil, :create_session, :encode)}],
         Map.put(state, :peer, updated_peer)}

      # TODO: Better error message
      _ ->
        {:close}
    end
  end

  defp observe(state) do
    resp =
      %{state.peer | type: :observer}
      |> Signaling.Peer.update()

    case resp do
      {:ok, updated_peer} ->
        {[{:text, Message.new!(nil, :observe, :encode)}], Map.put(state, :peer, updated_peer)}

      # TODO: Better error message
      _ ->
        {:close}
    end
  end
end
