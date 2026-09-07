# Finder - an index builder for REST APIs

built on [squirrel](https://github.com/Masterminds/squirrel) Query Builder
check the [finder-example](https://github.com/m-row/finder-example) for details
on usage

- [x] searchable
- [x] filterable
- [x] sortable
- [x] paginated

`IndexBuilder` groups by the model id by default for backwards compatibility.
Set `ConfigIndex.DisableDefaultGroupBy` to `true` when every configured join is
known to return at most one row per model and index-friendly pagination is
required. Supplying `GroupBys` explicitly continues to take precedence.

## CRUD - common operations

additionally there are postgres queries built to operate on one row and return
it in the same query

- [x] show
- [x] store
- [x] update
- [x] destroy
