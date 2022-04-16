package lighthouse

// Request specifies the format of a request to the lighthouse server
type Request struct {
	REID interface{}
	AUTH map[string]interface{}
	VERB string
	PATH []string
	META map[interface{}]interface{}
	PAYL interface{}
}

// Response specifies the format of a response from the lighthouse server
type Response struct {
	REID     interface{}
	RNUM     int
	RESPONSE string
	META     map[interface{}]interface{}
	PAYL     interface{}
	WARNINGS []interface{}
}

// Builder for Request
func NewRequest() *Request {
	return &Request{
		AUTH: map[string]interface{}{},
		META: map[interface{}]interface{}{},
	}
}

func (r *Request) Reid(reid interface{}) *Request {
	r.REID = reid
	return r
}

func (r *Request) Auth(user, token string) *Request {
	r.AUTH["USER"] = user
	r.AUTH["TOKEN"] = token
	return r
}

func (r *Request) Verb(verb string) *Request {
	r.VERB = verb
	return r
}

func (r *Request) Path(path ...string) *Request {
	r.PATH = path
	return r
}

func (r *Request) Meta(key interface{}, value interface{}) *Request {
	r.META[key] = value
	return r
}

func (r *Request) Payl(payl interface{}) *Request {
	r.PAYL = payl
	return r
}
