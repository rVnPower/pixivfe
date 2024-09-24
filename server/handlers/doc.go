/*
Package handlers provides HTTP request handling functionality for PixivFE.
It defines various middleware functions, request handlers, and routing logic to manage
incoming HTTP requests and produce appropriate responses.

The package includes middleware for security headers (SetPrivacyHeaders), error catching
and handling (CatchError, HandleError), rate limiting (IPRateLimiter, RateLimitRequest), logging (LogRequest),
panic recovery (RecoverFromPanic), and user context injection (ProvideUserContext). These
middlewares can be applied to routes to add cross-cutting functionality across multiple endpoints.

Route definitions are centralized in the DefineRoutes function, which sets up all paths
and their corresponding handlers using the gorilla/mux router.
*/
package handlers
