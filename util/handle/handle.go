package handle

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/juju/errors"
)

func Fatal(err error) {
	fmt.Println(err)
	fmt.Println("\n\n\n\n")
	glog.Fatal("%s", errors.ErrorStack(err))
}

func Error(err error) {
	glog.Error("%s", errors.ErrorStack(err))
}
