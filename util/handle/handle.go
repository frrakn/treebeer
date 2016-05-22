package handle

import (
	"github.com/golang/glog"
	"github.com/juju/errors"
)

func Fatal(err error) {
	glog.Fatal("%s", errors.ErrorStack(err))
}

func Error(err error) {
	glog.Error("%s", errors.ErrorStack(err))
}
