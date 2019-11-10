defmodule Signaling.Credential do
  def check_password(p1, p2) do
    case p1 == p2 do
      true -> {:ok}
      false -> {:error, "password is incorrect"}
    end
  end

  # TODO: Add hashing algorithm here
  def hash(password) do
    password
  end
end
