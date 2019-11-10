defmodule Signaling.Message do
  @type type ::
          :register
          | :observe
          | :create_session
          | :update
          | :connect
          | :confirm
          | :webrtc
          | :error
          | :warning
  @type payload :: any()
  @derive Jason.Encoder
  @type t :: %__MODULE__{
          type: type(),
          payload: payload()
        }
  @enforce_keys [:type, :payload]
  defstruct [:type, :payload]

  @spec new(payload(), type()) :: t()
  def new(payload, type)
      when type in [
             :register,
             :observe,
             :create_session,
             :update,
             :connect,
             :confirm,
             :webrtc,
             :error,
             :warning
           ] do
    %Signaling.Message{type: type, payload: payload}
  end

  @spec new(payload(), type(), :encode) ::
          {:ok, String.t()} | {:error, Jason.EncodeError.t() | Exception.t()}
  def new(payload, type, :encode) do
    new(payload, type) |> Jason.encode()
  end

  @spec new!(payload(), type()) :: t()
  def new!(payload, type) do
    new(payload, type)
  end

  @spec new!(payload(), type(), :encode) :: String.t()
  def new!(payload, type, :encode) do
    case new(payload, type, :encode) do
      {:ok, result} -> result
      {:error, _} -> ""
    end
  end

  @spec error(String.t()) :: String.t()
  def error(reason) do
    new!(%{reason: reason}, :error, :encode)
  end

  @spec warning(String.t()) :: String.t()
  def warning(reason) do
    new!(%{reason: reason}, :warning, :encode)
  end
end
