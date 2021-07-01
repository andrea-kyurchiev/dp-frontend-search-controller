package steps

import "github.com/cucumber/godog"

// RegisterSteps registers the specific steps needed to do component tests for the search controller
func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^all the downstream services are healthy$`, c.allTheDownstreamServicesAreHealthy)
	ctx.Step(`^one the downstream services is warning$`, c.oneOfTheDownstreamServicesIsWarning)
	ctx.Step(`^one the downstream services is failing$`, c.oneOfTheDownstreamServicesIsFailing)
}

func (c *Component) allTheDownstreamServicesAreHealthy() {
	c.FakeAPIRouter.setJSONResponseForGet("/health", 200)
	c.FakeRendererApp.setJSONResponseForGet("/health", 200)
}

func (c *Component) oneOfTheDownstreamServicesIsWarning() {
	c.FakeAPIRouter.setJSONResponseForGet("/health", 429)
	c.FakeRendererApp.setJSONResponseForGet("/health", 200)
}

func (c *Component) oneOfTheDownstreamServicesIsFailing() {
	c.FakeAPIRouter.setJSONResponseForGet("/health", 500)
	c.FakeRendererApp.setJSONResponseForGet("/health", 200)
}
