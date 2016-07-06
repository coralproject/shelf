// Package anvil provides support for validating an Anvil JWT and extracting
// the claims for authorization.
//
// [godoc](https://godoc.org/github.com/anvilresearch/go-anvil)
//
// An application needs to call the `/signin` and then the `/token` API calls.
// These calls authenticate the user for the application and provide the token
// required to make future calls into any webservice you are building that
// requires authentication/authorization.
//
// SignIn
//		You will need these values to call the signin API:
//		HOST         = 10.0.1.26:3000
//		CLIENTID     = 6b6efaae-0ab8-4152-8f92-a87c17921800
//		REDIRECT_URL = https://anvil.coralproject.net
//		EMAIL        = bill@thekennedyclan.net
//		PASSWORD     = Qfe^bJ9uD6cgnD-8
//		REFERRER     = https://anvil.coralproject.net/signin
//
// 		curl -X POST https://HOST/signin -d 'max_age=315569260&response_type=code&client_id=CLIENTID&redirect_uri=REDIRECT_URL&scope=openid%20profile%20email%20realm&provider=password&email=EMAIL&password=PASSWORD -H "referrer: REFERRER"
//
// Response
// 		Redirecting to https://anvil.coralproject.net?code=c9ce6c03ea6ad8dd3f0a%
//
// Token
//		You will need these values to call the token API:
//		HOST         = 10.0.1.26:3000
//		CLIENTID     = 6b6efaae-0ab8-4152-8f92-a87c17921800
//		REDIRECT_URL = https://anvil.coralproject.net
//		REFERRER     = https://anvil.coralproject.net/signin
//		CODE         = 6dafd2b59d6954849a6c  // From the response of the signin call
//
// 		curl -X POST https://CLIENTID:CODE@HOST/token -d 'grant_type=authorization_code&client_id=CLIENTID&code=CODE&redirect_uri=REDIRECT_URL' -H "referrer: REFERRER"
//
// Example
// 		// Create an Anvil value for the host we are using. Do this during
// 		// initialization.
// 		a, err := anvil.New("https://HOST")
// 		if err != nil {
// 		    // Log error and probably shutdown the service.
// 		    return
// 		}
//
// 		// This is an example handler that shows you how to use the Anvil value.
// 		handler := func(rw http.ResponseWriter, r *http.Request) {
//
// 		    // Have access to the Anvil value and use it to validate
// 		    // the request.
// 		    claims, err := a.ValidateFromRequest(r)
// 		    if err != nil {
//
// 		        // The token is not value so return an error.
// 		        rw.Header().Set("Content-Type", "application/json")
// 		        rw.WriteHeader(http.StatusUnauthorized)
// 		        json.NewEncoder(rw).Encode(struct{ Error string }{err.Error()})
// 		        return
// 		    }
//
// 		    // Everything is validated so move forward. The claims has what is
// 		    // need for authorization using the Scope field.
// 		    log.Println(claims.Scope)
// 		}
package anvil
