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

Set `ConfigIndex.RawTotal` when the caller already computed the unfiltered
total. Finder uses that value in pagination metadata and skips its duplicate
raw-total query.

## Base scopes and post-raw-total filters

Use `Wheres` for conditions that define the base dataset. These conditions
affect `raw_total`, `total`, and the returned rows. Tenant, business,
authorization, ownership, and other mandatory scope conditions belong here.

Use `FilterWheres` for optional filters that should be applied only after the
unfiltered count is calculated. These conditions affect `total` and the
returned rows, but do not affect `raw_total`.

```go
stateFilters := []squirrel.Sqlizer{
    squirrel.Eq{
        "COALESCE(item_exception.is_hidden, FALSE)": true,
    },
}

baseWheres := []squirrel.Sqlizer{
    squirrel.Eq{"products.business_id": businessID},
}

config := &finder.ConfigIndex{
    DB:           db,
    QB:           qb,
    Wheres:       &baseWheres,
    FilterWheres: &stateFilters,
    Selects:      &selects,
}
```

For example, if the base dataset contains 102 products and one product matches
the state filter, the pagination metadata is:

```json
{
  "raw_total": 102,
  "total": 1
}
```

`FilterWheres` is useful for filters that are difficult or inappropriate to
express through Finder's standard URL filter parser, including:

- computed SQL expressions such as `COALESCE(...)`;
- state derived from a joined table;
- aggregate or custom conditions built by the caller;
- endpoint-specific filters that must preserve the base `raw_total`.

Do not put access-control or tenant-scope conditions in `FilterWheres`, because
that would make `raw_total` include rows the caller should not see. Keep those
conditions in `Wheres`.

### Count-query behavior

Finder avoids the filtered count query when there are no active index filters.
In that case, `total` reuses `raw_total`. An index is considered filtered when
at least one of the following is present:

- `FilterWheres` contains a condition;
- `IsPublic` is enabled;
- the URL includes a non-empty `q` or `filters` value, or contains `from`.

When a filter is active, Finder calculates the base `raw_total`, applies all
post-raw-total and URL filters, and then calculates the filtered `total`.
`ConfigIndex.RawTotal` can still be supplied to reuse a count already computed
by the caller.

## CRUD - common operations

additionally there are postgres queries built to operate on one row and return
it in the same query

- [x] show
- [x] store
- [x] update
- [x] destroy
