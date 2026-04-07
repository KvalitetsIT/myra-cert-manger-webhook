package solvers

import (
	restclient "k8s.io/client-go/rest"
)

type Solver[T any] interface {
	Name() string
	Present(ch T) error
	CleanUp(ch T) error
	Initialize(kubeClientConfig *restclient.Config, stopCh <-chan struct{}) error
}
