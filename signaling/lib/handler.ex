defmodule Signaling.Handler do
    @behaviour :cowboy_websocket

    def init(req, _state) do
        # headers = :cowboy_req.headers(req) |> IO.inspect
        %{"name" => name} = req.headers
        {:cowboy_websocket, req, %{name: name}}
    end
    
    def websocket_init(state) do
        state |> IO.inspect
        # Signaling.Session.add(state.name, {})
        Registry.register(Registry.Peers, state.name, :observer)
        {:ok, state}
    end

    def websocket_handle(_frame, state) do
        Signaling.Peers.broadcast_update()
        {:ok, state}
    end

    def websocket_info({:update_peers, peers}, state) do
        {:reply, {:text, peers}, state}
    end
    
end