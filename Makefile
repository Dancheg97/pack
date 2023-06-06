
serve:
	sudo go run . s

test:
	go run . p localhost:4572/pack

certs:
	openssl req -x509 -newkey rsa:4096 -keyout docker/key.pem -out docker/cert.pem -sha256 -days 3650 -nodes -subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname" -addext "subjectAltName = DNS:localhost"
