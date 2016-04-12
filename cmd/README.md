# Profiling

Add this to main before Analyze call:

	f, err := os.Create("./cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

    $ go build
    $ ./colors ~/Desktop/*.jpg
    $ go tool pprof ./colors cpu2.prof
    Entering interactive mode (type "help" for commands)
    (pprof) web
