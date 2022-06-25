[![Go Report Card](https://goreportcard.com/badge/github.com/madzohan/tgpl)](https://goreportcard.com/report/github.com/madzohan/tgpl)
[![codecov](https://codecov.io/gh/madzohan/tgpl/branch/master/graph/badge.svg)](https://codecov.io/gh/madzohan/tgpl)
## "The Go Programming Language" book's exercises solving with test coverage
### by @madzohan

FAQ:
Q) How I can run tests manually?
A) Look at the job named "Run Unit tests" in `.github/workflows/go.yml`. Append `&& go tool cover -html=coverage.txt` to the cmd and you'll get nice html report opened in your default browser.
