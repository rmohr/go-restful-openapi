package restfulspec

import (
	"testing"

	restful "github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
)

func TestRouteToPath(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/tests/{v}")
	ws.Param(ws.PathParameter("v", "value of v").DefaultValue("default-v"))
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_XML)
	ws.Route(ws.GET("/a/{b}").To(dummy).
		Doc("get the a b test").
		Param(ws.PathParameter("b", "value of b").DefaultValue("default-b")).
		Param(ws.QueryParameter("q", "value of q").DefaultValue("default-q")).
		Returns(200, "list of a b tests", []Sample{}).
		Writes([]Sample{}))
	ws.Route(ws.GET("/a/{b}/{c:[a-z]+}").To(dummy).
		Doc("get the a b test").
		Param(ws.PathParameter("b", "value of b").DefaultValue("default-b")).
		Param(ws.PathParameter("c", "with regex").DefaultValue("abc")).
		Param(ws.QueryParameter("q", "value of q").DefaultValue("default-q")).
		Returns(200, "list of a b tests", []Sample{}).
		Writes([]Sample{}))

	p := BuildPaths(ws)
	t.Log(asJSON(p))

	if p.Paths["/tests/{v}/a/{b}"].Get.Parameters[0].Type != "string" {
		t.Error("Parameter type is not set.")
	}
	if _, exists := p.Paths["/tests/{v}/a/{b}/{c}"]; !exists {
		t.Error("Expected path to exist after it was sanitized.")
	}
}

func TestMultipleMethodsRouteToPath(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/tests/a")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_XML)
	ws.Route(ws.GET("/a/b").To(dummy).
		Doc("get a b test").
		Returns(200, "list of a b tests", []Sample{}).
		Writes([]Sample{}))
	ws.Route(ws.POST("/a/b").To(dummy).
		Doc("post a b test").
		Returns(200, "list of a b tests", []Sample{}).
		Returns(500, "internal server error", []Sample{}).
		Reads(Sample{}).
		Writes([]Sample{}))

	p := BuildPaths(ws)
	t.Log(asJSON(p))

	if p.Paths["/tests/a/a/b"].Get.Description != "get a b test" {
		t.Errorf("GET description incorrect")
	}
	if p.Paths["/tests/a/a/b"].Post.Description != "post a b test" {
		t.Errorf("POST description incorrect")
	}
	if _, exists := p.Paths["/tests/a/a/b"].Post.Responses.StatusCodeResponses[500]; !exists {
		t.Errorf("Response code 500 not added to spec.")
	}

	expectedRef := spec.MustCreateRef("#/definitions/restfulspec.Sample")
	postBodyparam := p.Paths["/tests/a/a/b"].Post.Parameters[0]
	postBodyRef := postBodyparam.Schema.Ref
	if postBodyRef.String() != expectedRef.String() {
		t.Errorf("Expected: %s, Got: %s", expectedRef.String(), postBodyRef.String())
	}

	if postBodyparam.Format != "" || postBodyparam.Type != "" || postBodyparam.Default != nil {
		t.Errorf("Invalid parameter property is set on body property")
	}
}
