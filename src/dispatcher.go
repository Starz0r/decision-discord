package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/spidernest-go/logger"
)

// Dispatcher will run the function fn after the duration d
// Optionally, a non-zero duration chkpnt will write it to disk
// to be restored later.
func Dispatch(d time.Duration, chkpnt time.Duration, fn func(), name string) {
	if chkpnt != 0 {
		// open a dispatcher checkpoint file for writing
		f, err := os.OpenFile("dispatcher_"+name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			logger.Fatal().
				Err(err).
				Msg("Dispatcher checkpoint file could not be opened for writing.")
		}

		// sleep and loop until it's past our time
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		go func(td *time.Duration, dd *time.Duration) {
			for *dd >= 0 {
				time.Sleep(*td)
				*dd -= *td
			}
		}(&chkpnt, &d)

		// if our program isn't killed we can write to disk on exit
		for d >= 0 {
			select {
			case <-ch:
				logger.Debug().Msg("Interrupt Recieved.")
				defer os.Exit(0)
				defer func(f *os.File, d time.Duration) {
					defer f.Close()
					f.WriteString(d.String())
				}(f, d)
				return
			}
		}
		fn()
	}
}
