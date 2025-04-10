package scheme

const (
	File       = "file"
	Data       = "data"
	Http       = "http"
	Https      = "https"
	HttpUnix   = Http + "+" + Unix
	Websocket  = "wss"
	Unix       = "unix"
	Unixgram   = "unixgram"
	Unixpacket = "unixpacket"
)

var (
	HttpsData = []string{Https, Data}
)
