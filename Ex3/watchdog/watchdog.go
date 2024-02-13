package watchdog

import (
	"time"
)


type Watchdog struct {
    period time.Duration
    callback func()
    feed_ch chan bool
}

func New(period time.Duration, feed_ch chan bool, callback func()) (*Watchdog) {
    return &Watchdog {
        period: period,
        callback: callback,
        feed_ch: feed_ch,
    }
}


func Feed(w *Watchdog) {
    w.feed_ch <- true 
}

func Start(w *Watchdog) {
    ticker := time.NewTicker(w.period)

    for {
        select {
            case <- ticker.C:
                w.callback()
            case <- w.feed_ch:
                ticker.Reset(w.period)
        }
    }
}

