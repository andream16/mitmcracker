.PHONY: clean

all:
	ocamlopt -o encrypt encrypt.ml
	ocamlopt -o decrypt decrypt.ml

clean:
	@'rm' *.cmo *.cmi *.cmx *.o encrypt decrypt 2>/dev/null || true
