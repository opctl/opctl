package plans

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/fakes"
)

var _ = Describe("describe handler", func() {
	It("returns the interface and forwards ref, credentials, and cancellation", func() {
		fake := new(fakes.FakeNode)
		fake.PlanOpReturns(&model.OpPlan{OpSpec: &model.OpSpec{Name: "hello"}}, nil)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"op":{"ref":"/ops/hello","pullCreds":{"Username":"user","Password":"pass"}}}`)).WithContext(ctx)
		response := httptest.NewRecorder()
		NewHandler(fake).Handle(response, request)
		Expect(response.Code).To(Equal(http.StatusOK))
		Expect(response.Body.String()).To(MatchJSON(`{"name":"hello","description":""}`))
		actualCtx, req := fake.PlanOpArgsForCall(0)
		Expect(actualCtx).To(Equal(ctx))
		Expect(req.Op.Ref).To(Equal("/ops/hello"))
		Expect(req.Op.PullCreds.Username).To(Equal("user"))
	})
	It("maps resolution and declaration failures to HTTP status codes", func() {
		for _, tc := range []struct {
			err    error
			status int
		}{
			{model.ErrDataProviderAuthentication{}, 401}, {model.ErrDataProviderAuthorization{}, 403}, {model.ErrDataRefResolution{}, 404}, {model.ErrDataUnableToResolve{}, 404}, {model.ErrDataMissingVersion{}, 400}, {model.ErrDataGitInvalidRef{}, 400}, {&node.OpDeclarationError{Err: errors.New("compile failed")}, 400}, {errors.New("internal failure"), 500}, {context.Canceled, 500}, {context.DeadlineExceeded, 500},
		} {
			fake := new(fakes.FakeNode)
			fake.PlanOpReturns(nil, tc.err)
			response := httptest.NewRecorder()
			NewHandler(fake).Handle(response, httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"op":{"ref":"/ops/hello"}}`)))
			Expect(response.Code).To(Equal(tc.status))
		}
	})
	It("rejects malformed requests and unsupported methods before describing", func() {
		fake := new(fakes.FakeNode)
		for _, tc := range []struct {
			method, body string
			status       int
		}{{"POST", "{", 400}, {"GET", "{}", 405}} {
			response := httptest.NewRecorder()
			NewHandler(fake).Handle(response, httptest.NewRequest(tc.method, "/", strings.NewReader(tc.body)))
			Expect(response.Code).To(Equal(tc.status))
		}
		Expect(fake.PlanOpCallCount()).To(Equal(0))
	})
})
