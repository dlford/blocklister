package blocklist

type BlockList struct {
	Title    string
	MaxElem  int
	Chains   []string
	IPs      []string
	Schedule string
}
