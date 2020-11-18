# jq sample

## command sample

### json 比較

```
diff <(cat a.json | jq . --sort-keys) <(cat b.json | jq . --sort-keys)
```