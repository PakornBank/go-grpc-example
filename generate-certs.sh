#!/bin/bash

ROOT_CERT_DIR="certs"
AUTH_CERT_DIR="auth/certs"
USER_CERT_DIR="user/certs"
GATEWAY_CERT_DIR="gateway/certs"

# Create directories
mkdir -p $ROOT_CERT_DIR $AUTH_CERT_DIR $USER_CERT_DIR $GATEWAY_CERT_DIR

echo "🚀 Generating Root Certificate Authority (CA)..."
openssl genrsa -out $ROOT_CERT_DIR/ca.key 4096
openssl req -x509 -new -key $ROOT_CERT_DIR/ca.key -sha256 -days 3650 -out $ROOT_CERT_DIR/ca.crt -subj "/CN=MyRootCA"

# Create OpenSSL config file for SAN (Subject Alternative Name)
OPENSSL_CNF="$ROOT_CERT_DIR/openssl.cnf"
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
    local dir=$2
    echo "🔹 Generating certificate for $name in $dir..."

    mkdir -p $dir

    # Generate private key
    openssl genrsa -out $dir/$name.key 4096

    # Generate Certificate Signing Request (CSR)
    openssl req -new -key $dir/$name.key -out $dir/$name.csr -config $OPENSSL_CNF

    # Sign the certificate with the CA
    openssl x509 -req -in $dir/$name.csr -CA $ROOT_CERT_DIR/ca.crt -CAkey $ROOT_CERT_DIR/ca.key -CAcreateserial \
        -out $dir/$name.crt -days 3650 -sha256 -extfile $OPENSSL_CNF -extensions v3_req

    # Copy CA certificate to the service's cert directory
    cp $ROOT_CERT_DIR/ca.crt $dir/

    # Remove CSR file
    rm -f $dir/$name.csr
}

# Generate certificates for gRPC servers
generate_cert "auth-grpc" "$AUTH_CERT_DIR"
generate_cert "user-grpc" "$USER_CERT_DIR"

generate_cert "api-gateway" "$GATEWAY_CERT_DIR"

echo "✅ All certificates generated successfully!"
echo "📂 Certificates stored in respective directories."
