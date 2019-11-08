defmodule Signaling.Peers do
    def observers do
        match_pattern = {:"$1", :"$2", :"$3"}
        guards = [{:==, :"$3", :observer}]
        body = [{match_pattern}]
        Registry.select(Registry.Peers, [{match_pattern, guards, body}])
    end

    def broadcast_update do
        observers = Signaling.Peers.observers |> IO.inspect
        observers_encoded = observers 
            |> Enum.map(&(elem(&1, 0)))
            |> Jason.encode!
        
        observers
            |> Enum.map(&(elem(&1, 1)))
            |> Enum.each(&(send(&1, {:update_peers, observers_encoded})))
    end
end