```sh
PROJECT_NAME=go-superinit
npx bmad-method install -f -i claude-code -d ./
go install github.com/spf13/cobra-cli@latest
go mod init github.com/joescharf/$PROJECT_NAME
cobra-cli init --viper --author "Joe Scharf joe@joescharf.com" --config $HOME/.config/$PROJECT_NAME
cobra-cli add version

mockery init $PROJECT_NAME
go mod tidy
git init
# add .gitignore for Go projects
echo "bin/" >> .gitignore
echo "*.exe" >> .gitignore
echo "*.dll" >> .gitignore
echo "*.so" >> .gitignore
echo "*.dylib" >> .gitignore
echo "vendor/" >> .gitignore
echo "*.test" >> .gitignore
echo "coverage.out" >> .gitignore
echo ".vscode/" >> .gitignore
echo ".idea/" >> .gitignore
echo "*.log" >> .gitignore
git add .
git commit -m "Initial commit"
```
