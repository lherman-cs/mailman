defmodule Signaling.Handler.Sender do
  @behaviour Signaling.Handler
  alias Signaling.Message

  def handle(%Message{:type => type, :payload => payload}, state) do
    case type do
      :confirm -> receive_confirmation(payload, state)
      _ -> {:ok, state}
    end
  end

  defp receive_confirmation(payload = %{"accept" => _accept}, state) do
    {[
       {:text, payload |> Message.new!(:confirm, :encode)}
     ], state}
  end
end
