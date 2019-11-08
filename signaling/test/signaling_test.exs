defmodule SignalingTest do
  use ExUnit.Case
  doctest Signaling

  test "greets the world" do
    assert Signaling.hello() == :world
  end
end
