package generator

//go:generate mockery --name=ShortLinkGenerator
type ShortLinkGenerator interface {
	GenerateShortLink() (string, error)
}
