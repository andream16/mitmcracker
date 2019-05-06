(*** SPN-DT - Substitution-Permutation network Decryption Tool - Copyright (C) 2014 Sebastian Podda - University of Cagliari ***)

(* Note: code comments are in italian. For help: sebastianpodda@unica.it *)

(* PARAMETRI DELLA SPN *)
let sbox_size = 4;;
let sbox_number = 16;;
let round_number = 2;;

(** SUBSTITUTION' e SUBSTITUTION **)
(*  Effettua una sostituzione (inversa) fissata su stringhe 4 bit (da modificare se si varia la sbox_size) *)
let substitution' str = 
	if (String.length str <> 4) then failwith "Invalid sbox size!" else
	match str with
		"0101" -> "0000"
	| "1100" -> "0001"
	| "0110" -> "0010" 
	| "0000" -> "0011"
	| "1000" -> "0100" 
	| "1111" -> "0101" 
	| "0111" -> "0110"
	| "0100" -> "0111"
  | "0001" -> "1000" 
	| "1010" -> "1001" 
	| "1101" -> "1010" 
	| "0010" -> "1011"
	| "1110" -> "1100"
	| "0011" -> "1101" 
	| "1001" -> "1110"
	| "1011" -> "1111"
	| _ -> "Error with substitution cipher!"
;;

let rec substitution str =
	match str with
		"" -> ""
	| _ -> let len = String.length str - 4 in 
		substitution' (String.sub str 0 4) ^ substitution (String.sub str 4 len)
;;

(** PERMUTATION' e PERMUTATION **)
(*	Restituisce una permutazione inversa (prefissata) di una stringa binaria lunga 64 caratteri *)
let rec permutation' str l  = 
	match l with
		[] -> ""
	| a::l' -> Char.escaped (String.get str a) ^ permutation' str l'
;;

let permutation str =
	if (String.length str == sbox_size * sbox_number) then 
		permutation' str [62; 10; 42; 18; 39; 2; 28; 13; 56; 61; 60; 45; 53; 51; 24; 38; 27; 49; 50;58; 17; 30; 48; 4; 63; 59; 44; 16; 21; 31; 19; 57; 20; 23; 12; 26; 37; 5; 8;35; 3; 34; 41; 32; 29; 25; 6; 52; 22; 40; 43; 1; 55; 14; 36; 9; 47; 46; 11;15; 33; 0; 54; 7]
	else
		failwith "Invalid pbox size!"
;;

(** BSTR_XOR **)
(*  Effettua lo xor tra due stringhe binarie della stessa dimensione *)
let rec bstr_xor s1 s2 =
	let len1 = String.length s1 in
	let len2 = String.length s2 in
	if (len1 <> len2) then failwith "Cannot xor binary strings with different size!"
	else
		if (len1 == 0) then "" else
		let c1 = String.get s1 0 in
		let c2 = String.get s2 0 in
		let sub1 = String.sub s1 1 (len1 -1) in
		let sub2 = String.sub s2 1 (len2 -1) in
		string_of_int ((Char.code c1 + Char.code c2) mod 2) ^ bstr_xor sub1 sub2
;;

(** ZERO_PADDING **)
(*	Costruisce una stringa composta da n zeri *)
let rec zero_padding n = 
	match n with
		0 -> ""
	| _ -> "0" ^ zero_padding (n-1)
;;

(** CHAR_TO_BINARY **)
(*  Converte un carattere in una stringa binaria di 8 bit *)
let char_to_binary c =
  let rec strip_bits i s =
    match i with
      0 -> s
    | _ -> strip_bits (i lsr 1) ((string_of_int (i land 0x01)) ^ s) in
  let res = strip_bits (Char.code c) "" in
	let len = String.length res in
	if (len < 8) then
		zero_padding (8-len) ^ res
	else res
;;

(** BINARY_TO_CHAR **)
(*  Converte una stringa binaria di 8 bit in un carattere ASCII esteso *)
let binary_to_char s =
  Char.chr (int_of_string ("0b" ^ s));;

(** BINARY_TO_STRING **)
(*  Data una stringa, esegue iterativamente binary_to_char su ogni carattere e restituisce il risultato *)
let rec binary_to_string str =
	match str with
		"" -> ""
	| _ -> let len = String.length str - 8 in 
		Char.escaped (binary_to_char (String.sub str 0 8)) ^ binary_to_string (String.sub str 8 len)
;;

(** STRING_TO_BINARY **)
(*  Come sopra, ma in direzione opposta *)
let rec string_to_binary str =
	match str with
		"" -> ""
	| str -> let c = String.get str 0 in
		char_to_binary c ^ (string_to_binary (String.sub str 1 (String.length str - 1)))
;;

(** REVERSE **)
(*  Restituisce la stringa passata, con i caratteri in ordine inverso *)
let rec reverse str =
	match str with
		"" -> ""
	| _ -> Char.escaped (String.get str (String.length str - 1)) ^ reverse (String.sub str 0 (String.length str - 1))
;;

(** GET_KEY **)
(*  Genera una chiave univoca a partire dalla chiave base e dal seed passato (tra 1 e 256) *)
let get_key key seed = 
	if (seed > 256 || seed < 1) then failwith "Invalid seed for generating keychain!" else
		let key_seed =
  	bstr_xor ((String.sub key 19 6) ^ "01")  (char_to_binary (Char.chr ((255 - seed)*246 mod 256))) ^
  	bstr_xor (String.sub key 12 8) (char_to_binary (Char.chr ((seed + 25) mod 256))) ^
  	bstr_xor ("00" ^ (String.sub key 14 5) ^ "1") (char_to_binary (Char.chr ((255 - seed)*37 mod 256))) ^
  	bstr_xor (String.sub key 7 8) (char_to_binary (Char.chr ((255 - seed + 22)*455 mod 256))) ^
  	bstr_xor (String.sub key 20 8) (char_to_binary (Char.chr (((seed + 154)*3) mod 256))) ^
  	bstr_xor ((String.sub key 5 4) ^ (String.sub key 25 4)) (char_to_binary (Char.chr ((seed + 97) mod 256))) ^
  	bstr_xor ((String.sub key 26 4) ^ (String.sub key 21 4)) (char_to_binary (Char.chr ((255 - seed + 19)*344 mod 256))) ^
  	bstr_xor (String.sub key 3 8) (char_to_binary (Char.chr ((seed*19 / 7) mod 256)))
		in
		bstr_xor ((bstr_xor key (reverse key)) ^ key) key_seed
		(* NB: di fatto, per il seed che va da 1 a 256, restituisce una keychain;
					 inoltre, piccole variazioni nel seed causano elevate variazioni nella chiave generata *)
;;

(** SPN_ALG' e SPN_ALG **)
(*	Esegue un algoritmo di decifratura SPN, espresso ricorsivamente *)
let rec spn_alg' w r k = 
	match r with
		0 -> w
	| _ -> spn_alg' (bstr_xor (substitution (permutation  w)) (get_key k r)) (r-1) k (* Eseguito Nr - 1 volte *)
;;

let spn_alg w k = spn_alg' (bstr_xor (substitution (bstr_xor w (get_key k (round_number+1)))) (get_key k round_number)) (round_number-1) k;;

(** HEX_TO_BINARY **)
(* 	Converte una stringa esadecimale in una stringa binaria *)
let rec hex_to_binary k =
	let len = String.length k in
	let s' = if (len>1) then hex_to_binary (String.sub k 1 (String.length k - 1)) else "" in
	match k.[0] with
	  '0' -> "0000" ^ s'
	| '1' -> "0001" ^ s'
	| '2' -> "0010" ^ s'
	| '3' -> "0011" ^ s'
	| '4' -> "0100" ^ s'
	| '5' -> "0101" ^ s'
	| '6' -> "0110" ^ s'
	| '7' -> "0111" ^ s'
	| '8' -> "1000" ^ s'
	| '9' -> "1001" ^ s'
	| 'A' -> "1010" ^ s'
	| 'B' -> "1011" ^ s'
	| 'C' -> "1100" ^ s'
	| 'D' -> "1101" ^ s'
	| 'E' -> "1110" ^ s'
	| 'F' -> "1111" ^ s'
	| 'a' -> "1010" ^ s'
	| 'b' -> "1011" ^ s'
	| 'c' -> "1100" ^ s'
	| 'd' -> "1101" ^ s'
	| 'e' -> "1110" ^ s'
	| 'f' -> "1111" ^ s'
	| _ -> failwith "An invalid hexadecimal string was passed to hex_to_binary!"
;;

(** BINARY_TO_HEX **)
(* 	Converte una stringa binaria in una stringa esadecimale *)
let rec binary_to_hex k =
	let len = String.length k in
	if (len mod 4 <> 0) then failwith "Binary string must have a length that is a multiple of 4!" else
	let s' = if (len>4) then binary_to_hex (String.sub k 4 (String.length k - 4)) else "" in
	match String.sub k 0 4 with
	  "0000" -> "0" ^ s'
	| "0001" -> "1" ^ s'
	| "0010" -> "2" ^ s'
	| "0011" -> "3" ^ s'
	| "0100" -> "4" ^ s'
	| "0101" -> "5" ^ s'
	| "0110" -> "6" ^ s'
	| "0111" -> "7" ^ s'
	| "1000" -> "8" ^ s'
	| "1001" -> "9" ^ s'
	| "1010" -> "A" ^ s'
	| "1011" -> "B" ^ s'
	| "1100" -> "C" ^ s'
	| "1101" -> "D" ^ s'
	| "1110" -> "E" ^ s'
	| "1111" -> "F" ^ s'
	| _ -> failwith "An invalid binary string was passed to binary_to_hex !"
;;


(** UNESCAPE **)
(*  Effettua l'unescape dei caratteri speciali sulla stringa passata *)
let unescape =
	let lexer = lazy (Genlex.make_lexer []) in
	fun s ->
		let tok_stream = Lazy.force lexer (Stream.of_string ("\"" ^ s ^ "\"")) in
		match Stream.peek tok_stream with
		| Some (Genlex.String s) -> s
		| _ -> assert false
;;

(** DECIPHER **)
(*  Decifra l'input utilizzando l'algoritmo SPN (al contrario) *)
let decipher str key = spn_alg str (hex_to_binary key);;

(** UNPAD' e UNPAD - DISABLED **)
(*  Rimuove il padding 100...000 finale alla stringa. NB: se la stringa non e' decifrata con la chiave giusta,
    la rimozione del padding elimina i bit finali (che non sono di padding) alterando la lunghezza della
		stringa stessa. Se pertanto tale lunghezza modulo 8 e' diversa da 0, il padding viene ripristinato *)
let rec unpad' str = 
	let len = String.length str - 1 in
	if (String.get str len  == '1') then
		(String.sub str 0 len)
	else
		unpad' (String.sub str 0 len)
;;

let unpad str = 
	let up = unpad' str in
	let remainder = (String.length up) mod 8 in
	if (remainder == 0) then up
	else up ^ "1" ^ zero_padding (7 - remainder)
;;	

let unpad str = str;; (* Dummy function *)

(** CBC' e CBC **)
(*  Decifra l'input a blocchi con mode of operation CBC *)
let rec cbc' str key vec =
	match str with
		"" -> ""
	| _ -> 
		let block = String.sub str 0 64  in
		bstr_xor vec (decipher block key) ^ cbc' (String.sub str 64 (String.length str - 64)) key block
;;

let cbc str key =
	if (String.length key == 8) then 
	(
  	try
  		unescape (binary_to_hex (unpad (cbc' (hex_to_binary str) key (zero_padding 64))))
  	with
  		x -> "Cannot unescape the decrypted text! Unescaped version here:\n\n" ^ 
  		binary_to_hex (unpad (cbc' (hex_to_binary str) key (zero_padding 64)))
	)
	else
		failwith "Wrong input! Only 32-bit keys are supported: please use 4 alphanumeric chars as a key!"
;;

(** READ_INPUT **)
(*  Legge dal canale di input e restituisce tutte le righe in una lista di stringhe *)
let read_input chan = 
let lines = ref [] in
try
  while true;
	do
    	lines := input_line chan:: !lines
  done; []
	with End_of_file -> 
		close_in chan; 
		List.rev !lines
;;

(** JOIN_STRINGS **)
(*  Data una lista di stringhe, le concatena e restituisce una stringa unica *)
let rec join_strings l = match l with
		[] -> ""
	| str::l' -> str ^ (join_strings l')
;;

let rec is_hex str = 
	if (String.length str < 1) then true
	else
		let s' = String.sub str 1 (String.length str - 1) in 
		match str.[0] with
		| '0' -> true && is_hex s'
		| '1' -> true && is_hex s'
		| '2' -> true && is_hex s'
		| '3' -> true && is_hex s'
		| '4' -> true && is_hex s'
		| '5' -> true && is_hex s'
		| '6' -> true && is_hex s'
		| '7' -> true && is_hex s'
		| '8' -> true && is_hex s'
		| '9' -> true && is_hex s'
		| 'A' -> true && is_hex s'
		| 'B' -> true && is_hex s'
		| 'C' -> true && is_hex s'
		| 'D' -> true && is_hex s'
		| 'E' -> true && is_hex s'
		| 'F' -> true && is_hex s'
		| _ -> false
;;

(* Il programma nasce per gestire chiavi di 32 bit ed estenderle a 64 bit per adattarla al blocco.
   L'informazione per raddoppiare le chiavi è comunque contenuta nei 32 bit scelti dall'utente.
	 La funzione qua sotto consente di usare anche chiavi da 20, 24 o 28 bit (per il challenge), 
	 semplicemente estendendole con zeri per arrivare a 32. Si noti che questo non comporta nessuna 
	 variazione significativa nella sicurezza della chiave (né in positivo, né in negativo) *)
let extend_key str len =
	if (len == 5) then "000" ^ str else
	if (len == 6) then "00" ^ str else
	if (len == 7) then "0" ^ str else
	str
;;

let check_hex str = 
	let len = String.length str in
	if (len > 4 && len < 9) then 
		(
		let s = String.uppercase (extend_key str len) in 
		if is_hex s then s else
			failwith "The plaintext must be a valid hexadecimal string!"
		)
	else failwith "The encryption key must be of 20, 24, 28 or 32 bits (5-8 hexadecimal chars)."
;;

let check_input str =
	let len = String.length str in
	if (len == 16) then 
		(
		let s = String.uppercase (str) in 
		if is_hex s then s else
			failwith "The plaintext must be a valid hexadecimal string!"
		)
	else failwith "This software version works only on plaintexts of 16 hexadecimal chars (64-bit)."
;;

(** MAIN **)
(*  Riceve gli argomenti del programma da riga di comando ed esegue la cifratura *)
let main = let argn = (Array.length Sys.argv) in
	match argn with
		  2 -> print_string (cbc (check_input (join_strings (read_input stdin))) (check_hex (Sys.argv.(1))) ^ "\n")
		| 3 -> print_string (cbc (check_input (join_strings (read_input (open_in (Sys.argv.(2)))))) (check_hex (Sys.argv.(1))) ^ "\n")
    | 4 ->
    (match Sys.argv.(1) with
      | "-s" -> print_string (cbc (check_input (Sys.argv.(3))) (check_hex (Sys.argv.(2))) ^ "\n")
      | _ -> print_string ("Wrong input!\nPlease use: $ encrypt [-s] <key> <text/filename>\n")
		)
		| argn -> print_string ("Wrong input!\nPlease use: $ decrypt [-s] <key> <text/filename>\n")
;;

