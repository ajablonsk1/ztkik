# Sprawozdanie - Lab 2

## 1. Środowisko testowe


- **Procesor**: Apple M2 Max (12-core: 8 Performance @ 3.5-3.7 GHz, 4 Efficiency @ 2.4 GHz)
- **Pamięć RAM**: 32 GB LPDDR5-6400
- **System operacyjny**: macOS
  
- **Język programowania**: Go (Golang)
- **Biblioteki kryptograficzne**: `crypto/rsa`, `crypto/aes`, `crypto/des`, `crypto/cipher`
- **Biblioteka do wykresów**: `gonum.org/v1/plot`

## 2. Wyniki i wykresy

### Wyniki

#### Generowanie kluczy

![](results/keygen/aes128.jpg)

![](results/keygen/aes256.jpg)

![](results/keygen/des192.jpg)

![](results/keygen/rsa2048.jpg)

![](results/keygen/rsa3072.jpg)

#### Szyfrowanie

![](results/encryption/aes128.jpg)

![](results/encryption/aes256.jpg)

![](results/encryption/3des192.jpg)

![](results/encryption/rsa2048.jpg)

#### Deszyfracja

![](results/decryption/aes128.jpg)

![](results/decryption/aes256.jpg)

![](results/decryption/3des192.jpg)

![](results/decryption/rsa2048.jpg)

### Wykresy

#### Generowanie kluczy

![](results/plots/keygen_asymmetric.png)

![](results/plots/keygen_symmetric.png)

![](results/plots/keygen_all.png)

#### Szyfrowanie

![](results/plots/encryption_all.png)

![](results/plots/encryption_all_4points.png)

![](results/plots/encryption_all_throughput.png)

#### Deszyfracja

![](results/plots/decryption_all.png)

![](results/plots/decryption_all_4points.png)

![](results/plots/decryption_all_throughput.png)

## 3. Dyskusja wyników

### Różnice i podobieństwa między algorytmami

Algorytmy symetryczne (AES, 3DES) są rzędy wielkości szybsze od RSA – zarówno w generowaniu kluczy, jak i w samym szyfrowaniu. AES-128 i AES-256 mają praktycznie identyczną wydajność, natomiast 3DES jest znacznie wolniejszy i osiąga przepustowość około 10 razy niższą niż AES. RSA z kolei ogranicza się do szyfrowania małych porcji danych – maksymalnie około 190 bajtów na operację, co w praktyce go dyskwalifikuje przy większych plikach.

### Wpływ długości klucza

W przypadku RSA przejście z klucza 2048-bitowego na 3072-bitowy wydłuża czas generowania kluczy około 4-5 razy. Dla AES różnica między wersją 128-bitową a 256-bitową jest minimalna – poniżej 10%. Dłuższe klucze dają lepsze bezpieczeństwo przy relatywnie niewielkim spadku wydajności.

### Wpływ rozmiaru danych

Dla małych danych (poniżej 1 KB) przepustowość jest niestabilna i szybko rośnie. Przy dużych plikach (powyżej 1 MB) przepustowość się stabilizuje – AES osiąga około 6000 MB/s, a 3DES około 450 MB/s. RSA przy dużych danych wypada bardzo słabo – ograniczenie do ~190 bajtów na operację wymusza dzielenie danych na kawałki, co dodatkowo wydłuża czas szyfrowania i czyni go niepraktycznym.

### Wady i zalety

**RSA** zapewnia bezpieczeństwo asymetryczne potrzebne do wymiany kluczy i podpisów cyfrowych, ale jest bardzo wolny. Najlepiej sprawdza się właśnie w tych dwóch scenariuszach – wymiana kluczy i podpisy.

**AES** ma bardzo wysoką przepustowość i działa wydajnie niezależnie od rozmiaru danych. Główna wada to konieczność bezpiecznej wymiany klucza między stronami. Idealny do szyfrowania dużych plików i komunikacji w czasie rzeczywistym.

**3DES** ma niską wydajność jak na algorytm symetryczny. Używany głównie w starszych systemach legacy. NIST zaleca jego wycofanie – obecnie powinien być stosowany tylko tam, gdzie konieczna jest kompatybilność wsteczna.

## 4. Rekomendacje zastosowania algorytmów

| Scenariusz                  | Rekomendacja  | Wytyczne NIST/ENISA/ISO       | Zgodność | Komentarz                                                                  |
| --------------------------- | ------------- | ----------------------------- | -------- | -------------------------------------------------------------------------- |
| Szyfrowanie dużych plików   | AES-256-GCM   | NIST SP 800-38D, ENISA 2024+  | Tak      | Bardzo szybki, plus sam sprawdza integralność                              |
| Komunikacja real-time       | AES-128-GCM   | NIST SP 800-57, ENISA         | Tak      | Minimalny lag, wystarczająco bezpieczny do 2030                            |
| Wymiana kluczy              | ECDH/X25519   | NIST SP 800-56A Rev. 3        | Tak      | Forward secrecy out of the box, RSA już przeżytek                          |
| Podpisy cyfrne              | RSA-3072/4096 | NIST FIPS 186-5, SP 800-57    | Tak      | Wszędzie działa, podpisujesz raz na jakiś czas więc wydajność nie gra roli |
| Szyfrowanie baz danych      | AES-256-GCM   | ISO/IEC 18033-3, NIST 800-38D | Tak      | Praktycznie zero overhead'u, integralność w pakiecie                       |
| Systemy legacy              | 3DES          | NIST: zakaz od 2024           | Nie      | Tylko jeśli musisz, planuj migrację ASAP                                   |
| Archiwizacja długoterminowa | AES-256-GCM   | NIST SP 800-57 Part 1 Rev. 5  | Tak      | Ochrona na następne 30+ lat, na razie jest quantum-resistant               |
| IoT                         | AES-128-GCM   | NIST SP 800-38D, ENISA        | Tak      | Lekki, często wsparty sprzętowo w MCU                                      |
