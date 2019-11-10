defmodule Signaling.Handler.Receiver do
  @behaviour Signaling.Handler
  alias Signaling.Message

  def handle(%Message{:type => type, :payload => payload}, state) do
    case type do
      :connect -> ask_confirmation(payload, state)
      :confirm -> answer_sender(payload, state)
      _ -> {:ok, state}
    end
  end

  defp ask_confirmation(%{"name" => name}, state) do
    {[
       {:text, %{name: name} |> Message.new!(:confirm, :encode)}
     ], state}
  end

  def answer_sender(%{"name" => sender, "accept" => accept}, state) do
    case Signaling.Peer.lookup(sender) do
      {:ok, {pid, _}} ->
        send(pid, %{"accept" => accept} |> Message.new!(:confirm))
        {:ok, state}

      {:error, reason} ->
        {[
           {:text, Message.warning(reason)}
         ], state}
    end
  end
end
