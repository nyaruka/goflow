package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net/http"

// 	"github.com/nyaruka/goflow/flows/flow"
// )

// type Error struct {
// 	Text string `json:"error"`
// }

// func WriteError(w http.ResponseWriter, err error) {
// 	w.WriteHeader(http.StatusBadRequest)
// 	json.NewEncoder(w).Encode(Error{err.Error()})
// }

// type ExpressionResponse struct {
// 	Result string   `json:"result"`
// 	Errors []string `json:"errors"`
// }

// type ExpressionRequest struct {
// 	Expression string          `json:"expression"`
// 	Context    json.RawMessage `json:"context"`
// }

// func Expression(w http.ResponseWriter, r *http.Request) {
// 	expression := ExpressionRequest{}

// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

// 	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
// 	if err != nil {
// 		WriteError(w, err)
// 		return
// 	}
// 	if err := r.Body.Close(); err != nil {
// 		WriteError(w, err)
// 		return
// 	}
// 	if err := json.Unmarshal(body, &expression); err != nil {
// 		WriteError(w, err)
// 		return
// 	}

// 	if expression.Context == nil || expression.Expression == "" {
// 		WriteError(w, fmt.Errorf("Missing context or expression element"))
// 		return
// 	}

// 	// build up our context
// 	//context, err := flow.ReadContext(expression.Context)
// 	if err != nil {
// 		WriteError(w, err)
// 		return
// 	}

// 	// evaluate it
// 	// result, err := excellent.EvaluateTemplateAsString(context.Environment(), context, expression.Expression)

// 	// w.WriteHeader(http.StatusOK)
// 	// response := ExpressionResponse{Result: result}
// 	// if err != nil {
// 	// 	switch err.(type) {
// 	// 	case excellent.TemplateErrors:
// 	// 		templateErrs := err.(excellent.TemplateErrors)
// 	// 		errs := make([]string, len(templateErrs))
// 	// 		for i := range errs {
// 	// 			errs[i] = templateErrs[i].Error()
// 	// 		}
// 	// 		response.Errors = errs
// 	// 	default:
// 	// 		response.Errors = []string{err.Error()}
// 	// 	}
// 	// }
// 	// json.NewEncoder(w).Encode(response)
// }

// type ExecuteRequest struct {
// 	Flow    json.RawMessage `json:"flow"`
// 	Context json.RawMessage `json:"context"`
// }

// func Execute(w http.ResponseWriter, r *http.Request) {
// 	execute := ExecuteRequest{}

// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

// 	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
// 	if err != nil {
// 		WriteError(w, err)
// 		return
// 	}
// 	if err := r.Body.Close(); err != nil {
// 		WriteError(w, err)
// 		return
// 	}
// 	if err := json.Unmarshal(body, &execute); err != nil {
// 		WriteError(w, err)
// 		return
// 	}

// 	if execute.Flow == nil || execute.Context == nil {
// 		WriteError(w, fmt.Errorf("Missing context or flow element"))
// 		return
// 	}

// 	// build up our context
// 	context, err := flow.ReadContext(execute.Context)
// 	if err != nil {
// 		WriteError(w, err)
// 		return
// 	}

// 	// and our flow
// 	// f, err := flow.ReadFlow(execute.Flow)
// 	// if err != nil {
// 	// 	WriteError(w, err)
// 	// 	return
// 	// }

// 	// // ok, off we go
// 	// if len(context.Path()) == 0 {
// 	// 	err = engine.StartFlow(f, context)
// 	// } else {
// 	// 	err = engine.ResumeFlow(f, context)
// 	// }
// 	// if err != nil {
// 	// 	WriteError(w, err)
// 	// 	return
// 	// }

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(context)
// }

func main() {
	// 	http.HandleFunc("/execute", Execute)
	// 	http.HandleFunc("/expression", Expression)
	// 	fmt.Println()
	// 	fmt.Println("Server running on port 8080")
	// 	fmt.Println("Endpoints:")
	// 	fmt.Println("  /execute    - run a flow. requires flow and context")
	// 	fmt.Println("  /expression - evaluate an expression. requires flow and expression")
	// 	fmt.Println()
	// 	log.Fatal(http.ListenAndServe(":8080", nil))
}
