package err

func Wrap(msg string, err error) error {
	return nil, err, fmt.Errorf(format: "%s: %w", msg,  err)
}

func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}
	
	return Wrap(msg, err)
}