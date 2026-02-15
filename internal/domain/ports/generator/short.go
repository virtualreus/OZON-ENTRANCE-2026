package generator

type ShortLinkGenerator interface {
	GenerateShortLink() (string, error)
}
