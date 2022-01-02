package core

import (
	"hade/app/route"
	"hade/framework"
	"hade/framework/gin"
)

func RunHttpEngine(container framework.Container) (*gin.Engine, error) {
	engine := gin.New()

	engine.SetContainer(container)

	engine.Use(gin.Recovery())

	route.Route(engine)

	return engine, nil
}
