# go-gitapp 

Description:
	This is a containerized Go app to authenticate with github using oauth2.
	It allows you to clone selected repository and then creates a PR with latest branch.

TODO:
	Not all corner cases are handled in this version. This is mainly to do POC.
	Unit tests.

Usage: 
	1. Go inside cloned directory.
	2. Run: docker build -t go_gitapp:1.0 .
	3. Run: docker run --publish 80:80 --detach --name go-gitapp go_gitapp:1.0
	4. Go to localhost OR 127.0.0.1 using web browser.
