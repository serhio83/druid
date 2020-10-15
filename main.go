package main

import (
	s "github.com/serhio83/druid/pkg/service"
)

func main() {
	service := s.Service{}
	service.Initialize()
	service.Run()
}
