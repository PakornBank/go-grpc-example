#!/bin/bash

CERT_DIR="certs"
mkdir -p $CERT_DIR

echo "🚀 Generating Root Certificate Authority (CA)..."
openssl genrsa -out $CERT_DIR/ca.key 4096
openssl req -x509 -new -key $CERT_DIR/ca.key -sha256 -days 3650 -out $CERT_DIR/ca.crt -subj "/CN=MyRootCA"

# Create OpenSSL config file for SAN (Subject Alternative Name)
OPENSSL_CNF="$CERT_DIR/openssl.cnf"
cat > $OPENSSL_CNF <<EOL
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_req
prompt = no

[req_distinguished_name]
CN = localhost

[v3_req]
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth, clientAuth
subjectAltName = @alt_names

[alt_names]
IP.1 = 127.0.0.1
DNS.1 = localhost
EOL

generate_cert() {
    local name=$1
    echo "🔹 Generating certificate for $name..."

    # Generate private key
    openssl genrsa -out $CERT_DIR/$name.key 4096

    # Generate Certificate Signing Request (CSR)
    openssl req -new -key $CERT_DIR/$name.key -out $CERT_DIR/$name.csr -config $OPENSSL_CNF

    # Sign the certificate with the CA
    openssl x509 -req -in $CERT_DIR/$name.csr -CA $CERT_DIR/ca.crt -CAkey $CERT_DIR/ca.key -CAcreateserial \
        -out $CERT_DIR/$name.crt -days 3650 -sha256 -extfile $OPENSSL_CNF -extensions v3_req
}

# Generate certificates for gRPC servers
generate_cert "auth-grpc"
generate_cert "user-grpc"

# Generate certificate for API Gateway (client)
generate_cert "api-gateway"

# Remove unnecessary CSR files
rm -f $CERT_DIR/*.csr

echo "✅ All certificates generated successfully!"
echo "📂 Certificates stored in: $CERT_DIR"