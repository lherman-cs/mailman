defmodule Signaling.Peer do
  alias Signaling.Message
  @type type :: :waiter | :observer | :sender | :receiver

  @type t :: %__MODULE__{
          type: type,
          name: String.t(),
          password: String.t()
        }
  @enforce_keys [:type, :name]
  defstruct [:type, :name, :password]

  @spec new(type(), String.t(), String.t()) :: t()
  def new(type, name, password \\ "")
      when type in [:waiter, :observer, :sender, :receiver] and is_binary(name) and
             is_binary(password) do
    %Signaling.Peer{type: type, name: name, password: password}
  end

  @spec register(t()) :: {:ok} | {:error}
  def register(peer) do
    case Registry.register(__MODULE__, peer.name, peer) do
      {:ok, _} -> {:ok}
      {:error, _} -> {:error}
    end
  end

  def update(peer) do
    case Registry.update_value(__MODULE__, peer.name, fn _ -> peer end) do
      :error -> {:error}
      _ -> {:ok, peer}
    end
  end

  @spec find(type()) :: {:ok, [term()]} | {:error}
  def find(type) do
    match_pattern = {:"$1", :"$2", :"$3"}
    guards = [{:==, {:map_get, :type, :"$3"}, type}]
    body = [{match_pattern}]
    spec = [{match_pattern, guards, body}]

    try do
      {:ok, Registry.select(__MODULE__, spec)}
    catch
      {:error, _} -> {:error}
    end
  end

  @spec lookup(String.t()) :: {:ok, {pid(), t()}} | {:error, String.t()}
  def lookup(key) when is_binary(key) do
    case Registry.lookup(__MODULE__, key) do
      [first | _] -> {:ok, first}
      [] -> {:error, "#{key} is not registered"}
    end
  end

  def broadcast_update do
    new_message = &%Signaling.Message{type: :update, payload: &1}

    with {:ok, receivers} <- find(:receiver),
         {:ok, observers} <- find(:observer) do
      message =
        receivers
        |> Enum.map(&elem(&1, 0))
        |> Message.new!(:update, :encode)
        |> new_message.()

      observers
      |> Enum.map(&elem(&1, 1))
      |> Enum.each(&send(&1, message))
    else
      _ -> {:error}
    end
  end
end
