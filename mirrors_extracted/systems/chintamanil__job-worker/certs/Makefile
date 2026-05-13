# Certificate generation using cfssl
# Usage: make all

CLIENTS := client1 client2 admin

.PHONY: all clean verify verify-chain

all: ca.pem server.pem $(addsuffix .pem,$(CLIENTS))

ca.pem: ca-csr.json ca-config.json
	cfssl gencert -initca ca-csr.json | cfssljson -bare ca

server.pem: ca.pem server-csr.json
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server-csr.json | cfssljson -bare server

%.pem: ca.pem %-csr.json
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client $*-csr.json | cfssljson -bare $*

verify:
	@for f in ca server $(CLIENTS); do echo "=== $$f ===" && openssl x509 -in $$f.pem -noout -subject -dates; done

verify-chain:
	@for f in server $(CLIENTS); do openssl verify -CAfile ca.pem $$f.pem; done

clean:
	rm -f *.pem *.csr
