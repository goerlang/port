-module(pingpong).

-export([ start/1
        , start/2
        , ping/1
        ]).

start(line) ->
  spawn(fun() -> init(line, 0) end).

start(packet, N) when N == 1;
                      N == 2;
                      N == 4 ->
  spawn(fun() -> init(packet, N) end).

init(Mode, Size) ->
  register(pingpong, self()),
  process_flag(trap_exit, true),
  Port = erlang:open_port({spawn_executable, "/usr/bin/env"},
                          case Mode of
                            line ->
                              [{line, 65536}, {args, ["pingpong", "-l"]}];
                            packet ->
                              [{packet, Size}, {args, [ "pingpong"
                                                      , "-p"
                                                      , "-psize=" ++ integer_to_list(Size)
                                                      ]
                                               }
                              ]
                          end),
  loop(Port).

ping(Data) ->
  pingpong ! {ping, self(), Data},
  receive Answer -> Answer end.

loop(Port) ->
  io:format("loop~n"),
  receive
    {ping, Caller, Msg} ->
      true = erlang:port_command(Port, Msg),
      receive
        {Port, {_, Data}} ->
          Caller ! {pong, Data};
        {'EXIT', Port, Reason} ->
          io:format("port exited: ~p~n", [Reason]),
          exit(port_terminated);
        Other ->
          Caller ! {other, Other}
      end,
      loop(Port);
    stop ->
      Port ! {self(), close},
      receive
        {Port, closed} ->
          io:format("port closed~n"),
          exit(normal)
      end;
    {'EXIT', Port, Reason} ->
      io:format("port exited: ~p~n", [Reason]),
      exit(port_terminated)
  end.
