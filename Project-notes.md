Ran this command:
go mod init social-app

need to change this to the actual github repo

setting-up development patterns:
1. made bin folder - contains the compiled application
2. cmd folder - contains the executables or entry points
3. cmd/api folder 
4. cmd/migrate folder 
5. cmd/migrate/migrations folder 
6. make internal folder - interacts with postgres, data validations, rate limiting etc
7. docs folder
8. Scripts folder