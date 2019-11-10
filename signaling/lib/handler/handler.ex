defmodule Signaling.Handler do
  @behaviour :cowboy_websocket
  alias Signaling.Handler.Waiter
  alias Signaling.Handler.Observer
  alias Signaling.Handler.Receiver
  alias Signaling.Handler.Sender
  alias Signaling.Message

  @callback handle(Message.t(), map()) :: any()

  def init(req, _state) do
    %{"name" => name} = req.headers

    {:cowboy_websocket, req,
     %{
       peer: Signaling.Peer.new(:waiter, name)
     }}
  end

  def websocket_init(state = %{:peer => peer}) do
    response =
      case Signaling.Peer.register(peer) do
        {:ok} ->
          [{:text, Message.new!(nil, :register, :encode)}]

        {:error} ->
          [
            {:text, Message.error("#{peer.name} is already taken")},
            :close
          ]
      end

    IO.inspect({response, state})
    {response, state}
  end

  def websocket_handle({:text, raw_message}, state) do
    case Jason.decode(raw_message) do
      {:ok, message} ->
        %Message{type: String.to_atom(message["type"]), payload: message["payload"]}
        |> route(state)

      {:error, _} ->
        {[
           {:text, Message.error("invalid message format")},
           :close
         ], state}
    end
  end

  def websocket_handle(_frame, state) do
    {[
       {:text, Message.error("invalid message type")},
       :close
     ], state}
  end

  def websocket_info(message, state) do
    route(message, state)
  end

  defp route(message, state = %{:peer => peer}) do
    IO.inspect(%{
      message: message,
      state: state
    })

    case peer.type do
      :waiter -> Waiter.handle(message, state)
      :observer -> Observer.handle(message, state)
      :receiver -> Receiver.handle(message, state)
      :sender -> Sender.handle(message, state)
    end
  end
end
