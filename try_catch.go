package main

// func someFunc() error {
//   var err error
//   out := DataItem{}
//   TryCatchBlock {
//     Try: func() {
//       ...
//       out.Fill()
//       Throw(makeError(err, fileline))
//       ...
//     },
//     Catch: func(e Exception) {
//       log.Printf("%v\n", e)
//       err = fmt.Errorf("%v", e)
//     },
//   }.Do()
//   return err
// }

type TryCatchBlock struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func (tcb TryCatchBlock) Do() {
	if tcb.Finally != nil {
		defer tcb.Finally()
	}
	if tcb.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcb.Catch(r)
			}
		}()
	}
	tcb.Try()
}
