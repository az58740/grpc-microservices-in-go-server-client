
echo "####START####"


mkdir -p cert && cd cert
echo "Generating private key and self-signed certificate for CA..."
openssl req -x509 \
    -sha256 \
    -newkey rsa:4096 \
    -days 365 \
    -keyout ca.key \
    -out ca.crt \
    -subj "/C=IR/ST=AZ.Sharghi/L=Tabriz/O=Software/OU=Microservices/CN=*.microservices.dev/emailAddress=alizeinalzadeh@microservices.dev" \
    -nodes

    
echo "Generate private key and certificate signing request for server"
openssl req \
    -sha256 \
    -newkey rsa:4096 \
    -keyout server.key \
    -out server-req.pem \
    -subj "/C=IR/ST=AZ.Sharghi/L=Tabriz/O=Software/OU=Microservices/CN=*.microservices.dev/emailAddress=alizeinalzadeh@microservices.dev" \
    -nodes

echo "Sign certificate signing request for server by using private key of CA"
rm server-ext.cnf || true && echo "subjectAltName=DNS:*.microservices.dev,DNS:*.microservices.dev,IP:0.0.0.0" >> server-ext.cnf
openssl x509 \
    -req -in server-req.pem \
    -sha256 \
    -days 60 \
    -CA ca.crt \
    -CAkey ca.key \
    -CAcreateserial \
    -out server.crt \
    -extfile server-ext.cnf
    
echo "Generate private key and certificate signing request for client"
openssl req \
    -sha256 \
    -newkey rsa:4096 \
    -keyout client.key \
    -out client-req.pem \
    -subj "/C=IR/ST=AZ.Sharghi/L=Tabriz/O=Software/OU=Microservices/CN=*.microservices.dev/emailAddress=alizeinalzadeh@microservices.dev" \
    -nodes

    echo "Sign certificate signing request for client by using private key of CA"
rm client-ext.cnf || true && echo "subjectAltName=DNS:*.microservices.dev,DNS:*.microservices.dev,IP:0.0.0.0" >> client-ext.cnf
openssl x509 \
    -req -in client-req.pem \
    -sha256 \
    -days 60 \
    -CA ca.crt \
    -CAkey ca.key \
    -CAcreateserial \
    -out client.crt \
    -extfile client-ext.cnf
    