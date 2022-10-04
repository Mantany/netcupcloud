package server

type serverType struct {
}

type shopProduct struct {
	name   string
	shopID int
}

type location struct {
	city    string
	country string
	long    string
	lang    string
}

type server struct {
	rDNS     string
	vpsID    string
	location location
	product  shopProduct
}
