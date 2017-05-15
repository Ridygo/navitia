package navitia

import (
  "context"
  "strings"

  "github.com/aabizri/navitia/types"
)

// Results stores results common informations
type Results struct {
  BaseUrl     string
  Scope       *Scope
  Session     *Session
}

type resultsobject interface {
  IDString() string
}

// creating stores creation time
func (r *Results) baseinfos(url string, session *Session) {
	r.BaseUrl = strings.Split(url, "?")[0]
  r.Session = session
}

// Single value results case
type SingleValueResults struct {
  Results
}

// creating stores creation time
func (r *SingleValueResults) Explore(ctx context.Context, selector string, opts ExploreRequest) (*ExploreResults, error) {
	// Create the URL
	url := r.BaseUrl + "/" + selector

	// Call
	return r.Session.explore(ctx, url, opts)
}

type MultiValuesResults struct {
  Results
}

// func (s *Session) RelatedContainers(obj []types.Container) (res []SingleValueResults) {
//   res = []SingleValueResults{}
//   for _, o := range obj {
//     result := SingleValueResults{}
//     id := string(o.ID)
//     result.BaseUrl = r.apiURL + "/" + id
//     result.Session = r.Session
//     res = append(res, result)
//   }
//   return
// }
func (s *Scope) RelatedContainers(obj []types.Container) (res []SingleValueResults) {
  res = []SingleValueResults{}
  for _, o := range obj {
    resType, _ := o.ID.Type()
  	typeSelector, _ := resourceTypeToSelector[resType]

    result := SingleValueResults{}
    id := string(o.ID)

    result.BaseUrl = s.baseURL + "/" + typeSelector + "/" + string(id)
    result.Session = s.session
    res = append(res, result)
  }
  return
}

//TODO other types

//
// // creating stores creation time
// func (r *MultiValuesResults) Explore(ctx context.Context, obj resultsobject, selector string, opts ExploreRequest) (*ExploreResults, error) {
//   var id string
//   if reflect.TypeOf(obj.ID).Kind() == reflect.Func {
//     id = obj.ID()
//   } else {
//     id = obj.ID
//   }
// 	// Create the URL
// 	url := r.BaseUrl + "/" + id + "/" + selector
//
// 	// Call
// 	return r.Session.explore(ctx, url, opts)
// }
