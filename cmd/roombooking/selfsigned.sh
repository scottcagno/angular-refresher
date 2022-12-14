#!/usr/bin/env bash

# Setup: Create our base directory structure
mkdir -p cert/CA/localhost

# Step 1: Generate a CA certificate
#
cd cert/CA

# 1.1: Generate a private key
openssl genrsa -out CA.key -des3 2048

# 1.2: Generate a root CA certificate using the private key
openssl req -x509 -sha256 -new -nodes -days 3650 -key CA.key -out CA.pem

# Step 2: Generate a certificate
#
cd localhost
touch localhost.txt
echo "authorityKeyIdentifier = keyid,issuer" >> localhost.txt
echo "basicConstraints = CA:FALSE" >> localhost.txt
echo "keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment" >> localhost.txt
echo "subjectAltName = @alt_names" >> localhost.txt
echo "" >> localhost.txt
echo "[alt_names]" >> localhost.txt
echo "DNS.1 = localhost" >> localhost.txt
echo "IP.1 = 127.0.0.1" >> localhost.txt

# 2.1: Generate the localhost private key
openssl genrsa -out localhost.key -des3 2048

# 2.2: Generate the CSR using the private key
openssl req -new -key localhost.key -out localhost.csr

# 2.3: Now with this CSR, we can request the CA to sign a certificate
openssl x509 -req -in localhost.csr -CA ../CA.pem -CAkey ../CA.key -CAcreateserial -days 3650 -sha256 -extfile localhost.ext -out localhost.crt

# 2.4: We will also need to decrypt the localhost.key and store that file too
openssl rsa -in localhost.key -out localhost.decrypted.key

cd ../..