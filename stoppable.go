package main

import(
  "fmt"
  "time"
  "log"
)

// IStoppable is basic interface for every Stoppable-derived struct
type IStoppable interface {
  Construct() error
  Do() error
  Destruct() error
}

// Stoppable is used to stop loops gracefully
type Stoppable struct {
  Name        string
  SleepMs     int   // time to sleep beetween working iterations (ms)
  Stopped     chan struct{}
  Stopping    chan struct{}
  self        IStoppable
}

// Start initialized channels needed for stopping
func (s *Stoppable) Start(obj IStoppable, name string, period ...int) error {
  s.self = obj
  s.Name = name
  err := s.self.Construct()
  if err != nil {
    err = makeError(fmt.Errorf("%s CONSTRUCT => %s", s.Name, err),
      FileLine())
    return err
  }

  s.SleepMs = 10
  if len(period) > 0 {
    s.SleepMs = period[0]
  }
  log.Printf("%s SLEEP = %d ms", s.Name, s.SleepMs)
  s.Stopping = make(chan struct{})
  s.Stopped = make(chan struct{})
  go s.Loop()
  log.Println(s.Name + " started")
  return nil
}

// Stop must be called outside loop goroutine
func (s *Stoppable) Stop() error {
  close(s.Stopping)
  <-s.Stopped
  log.Println(s.Name + " stopped")

  err := s.self.Destruct()
  if err != nil {
    err = makeError(fmt.Errorf("%s DESTRUCT => %s", s.Name, err),
      FileLine())
    log.Println(err.Error())
    return err
  }
  return nil
}

// Loop must start asynchronously to run each routine() as loop iteration
func (s *Stoppable) Loop() error {
  defer close(s.Stopped)
  for {
    select {
    case <-s.Stopping:
      log.Println(s.Name + " stopping...")
      return nil
    default:
      err := s.self.Do()
      if err != nil {
        log.Println(makeError(fmt.Errorf("%s causes %s", s.Name, err),
          FileLine()).Error())
      }
      if s.SleepMs > 0 {
        time.Sleep(time.Duration(s.SleepMs) * time.Millisecond)
      }
    } // select
  }
  return nil
}
