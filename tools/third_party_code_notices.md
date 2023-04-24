# License Compliance

Generate the `Third_Party_Code` directory by running the following:

```shell
go run github.com/chrismarget-j/go-licenses save --save_path=./Third_Party_Code --force ./...
go run github.com/chrismarget-j/go-licenses report --ignore github.com/Juniper/apstra-go-sdk ./... --template .notices.tpl > Third_Party_Code/NOTICES.md
```
