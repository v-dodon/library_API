context
=======
[![Build Status](https://travis-ci.org/github.com.github.com.github.com.gorilla/context.png?branch=master)](https://travis-ci.org/github.com.github.com.github.com.gorilla/context)

github.com.github.com.github.com.gorilla/context is a general purpose registry for global request variables.

> Note: github.com.github.com.github.com.gorilla/context, having been born well before `context.Context` existed, does not play well
> with the shallow copying of the request that [`http.Request.WithContext`](https://golang.org/pkg/net/http/#Request.WithContext) (added to net/http Go 1.7 onwards) performs. You should either use *just* github.com.github.com.github.com.gorilla/context, or moving forward, the new `http.Request.Context()`.

Read the full documentation here: http://www.gorillatoolkit.org/pkg/context
