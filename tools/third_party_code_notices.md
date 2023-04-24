# License Compliance

Generate the `Third_Party_Code` directory by running the following:

```shell
go run github.com/google/go-licenses save --save_path=./Third_Party_Code --force ./...
go run github.com/google/go-licenses report --ignore github.com/Juniper/apstra-go-sdk ./... --template .notices.tpl > Third_Party_Code/NOTICES.md
```
