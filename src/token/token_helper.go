package token

func (t *Token) String() string {
	return string(t.Lit)
}
