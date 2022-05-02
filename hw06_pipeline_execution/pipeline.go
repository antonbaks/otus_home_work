package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		stageCh := readStageChWithDoneCh(in, done)
		in = stage(stageCh)
	}

	return in
}

func readStageChWithDoneCh(in In, done In) Out {
	stageCh := make(Bi)

	go func() {
		defer close(stageCh)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}

				stageCh <- v
			}
		}
	}()

	return stageCh
}
