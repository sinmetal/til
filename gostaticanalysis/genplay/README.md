# genplay

https://github.com/gostaticanalysis/skeleton を利用してCode生成を行う

## Skeleton install

```
go get -u github.com/gostaticanalysis/skeleton
```

## skeleton 生成

以下のコマンドで生成したものを改修している

```
skeleton -type=codegen genplay
```

## Develop & Update

以下のコマンドで genpayl.golden を再生成して、作られたものを見れるので、こいつを頻繁に実行しながら、ちびちび書いていく
templete は knife の go doc を見ながら、作っていく。 struct だと https://pkg.go.dev/github.com/gostaticanalysis/knife#Struct とか。

```
go test -update && cat testdata/src/a/genplay.golden
```