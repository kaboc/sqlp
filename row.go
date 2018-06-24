package sqlp

type Row struct {
	rows *Rows
	err  error
}

func (r *Row) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	defer r.rows.Close()

	if !r.rows.Next() {
		if err := r.rows.Err(); err != nil {
			return err
		}
		return ErrNoRows
	}

	err := r.rows.Scan(dest...)
	if err != nil {
		return err
	}

	return r.rows.Close()
}

func (r *Row) ScanToStruct(structPtr interface{}) error {
	if r.err != nil {
		return r.err
	}
	defer r.rows.Close()

	if !r.rows.Next() {
		if err := r.rows.Err(); err != nil {
			return err
		}
		return ErrNoRows
	}

	err := r.rows.ScanToStruct(structPtr)
	if err != nil {
		return err
	}

	return r.rows.Close()
}

func (r *Row) ScanToMap() (map[string]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	defer r.rows.Close()

	if !r.rows.Next() {
		if err := r.rows.Err(); err != nil {
			return nil, err
		}
		return nil, ErrNoRows
	}

	result, err := r.rows.ScanToMap()
	if err != nil {
		return nil, err
	}

	return result, r.rows.Close()
}

func (r *Row) ScanToSlice() ([]string, error) {
	mp, err := r.ScanToMap()
	if err != nil {
		return nil, err
	}

	result := make([]string, len(r.rows.columnNames))
	for i, v := range r.rows.columnNames {
		result[i] = mp[v]
	}

	return result, nil
}
