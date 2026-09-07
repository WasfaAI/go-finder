package finder

import "context"

func rawTotal(c *configFinder) error {
	results := *c.Results
	if c.GroupBys != nil {
		if len(*c.GroupBys) > 0 {
			results = results.
				GroupBy(*c.GroupBys...)
		}
	}
	// getting total count of unfiltered result -------------------------------
	var total uint64
	totalQuery, args, err := c.QB.
		Select(`COUNT(*)`).
		FromSelect(results, "results").
		ToSql()
	if err != nil {
		return err
	}
	row := c.DB.QueryRowContext(context.Background(), totalQuery, args...)
	if err := row.Scan(&total); err != nil {
		total = 0
	}
	c.Meta.RawTotal = total
	return nil
}
