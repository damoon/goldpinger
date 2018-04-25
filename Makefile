run:
	elm-reactor --address=0.0.0.0

live-elm:
	elm-live Main.elm

live-go:
	CompileDaemon -command="./goldpinger"
