# Sprawozdanie: Kryptoanaliza szyfru afinicznego

**Autor:** Adrian Jabłoński  
**Data:** 11 października 2025  
**Przedmiot:** Zaawansowane techniki kryptografii i kryptoanalizy

---

## 1. Wstęp

Celem zadania była kryptoanaliza szyfru afinicznego z wykorzystaniem **ataku ze znanym szyfrogramem**. Szyfr afiniczny to szyfr podstawieniowy monoalfabetyczny definiowany wzorem:

$$e(x) = (ax + b) \mod 26$$

gdzie:

- $x$ - pozycja litery tekstu jawnego (A=0, B=1, ..., Z=25)
- $a, b$ - klucz szyfrujący
- $\gcd(a, 26) = 1$ (warunek konieczny odwracalności)

Funkcja deszyfrująca ma postać:

$$d(y) = a^{-1}(y - b) \mod 26$$

gdzie $a^{-1}$ to odwrotność modularna liczby $a$ modulo 26.

---

## 2. Analiza szyfrogramu

### 2.1. Dane wejściowe

**Szyfrogram:**

```
YAE VCCX ZA FWSG ZRC I IVX WZ QEUZ NC SAFPWQC OWZR ZRC LCVMZR AH ZRC ILFRINCZ
```

### 2.2. Histogram częstości liter

Pierwszym krokiem była analiza częstości występowania poszczególnych liter w szyfrogramie:

| Litera | Liczność | Procent | Wizualizacja |
|--------|----------|---------|--------------|
| Z | 9 | 14.8% | █████████ |
| C | 9 | 14.8% | █████████ |
| R | 6 | 9.8% | ██████ |
| W | 4 | 6.6% | ████ |
| I | 4 | 6.6% | ████ |
| A | 4 | 6.6% | ████ |
| V | 3 | 4.9% | ███ |
| F | 3 | 4.9% | ███ |
| Q | 2 | 3.3% | ██ |
| N | 2 | 3.3% | ██ |
| E | 2 | 3.3% | ██ |
| L | 2 | 3.3% | ██ |
| X | 2 | 3.3% | ██ |
| S | 2 | 3.3% | ██ |

---

## 3. Metoda kryptoanalizy - układ równań liniowych

### 3.1. Próby złamania klucza

#### **Próba 1:** Założenie E→Z, T→C

Dane:

- $p_1 = E = 4$, $c_1 = Z = 25$
- $p_2 = T = 19$, $c_2 = C = 2$

Obliczenia:

- $\Delta x = (4 - 19) \mod 26 = -15 \mod 26 = 11$
- $\Delta y = (25 - 2) \mod 26 = 23$
- $\Delta x^{-1} = 11^{-1} \mod 26 = ?$

Sprawdzenie: $\gcd(11, 26) = 1$ ✓

Szukamy $11 \cdot x \equiv 1 \pmod{26}$:

- $11 \cdot 19 = 209 = 8 \cdot 26 + 1$ ✓
- Zatem $\Delta x^{-1} = 19$

Obliczamy $a$:

- $a = (23 \cdot 19) \mod 26 = 437 \mod 26 = 21$

Sprawdzenie: $\gcd(21, 26) = 1$ ✓

Obliczamy $b$:

- $b = (25 - 21 \cdot 4) \mod 26 = (25 - 84) \mod 26 = -59 \mod 26 = 19$

**Wynik:** $a = 21$, $b = 19$

**Deszyfrowanie:**
Potrzebujemy $a^{-1} = 21^{-1} \mod 26$:

- $21 \cdot 5 = 105 = 4 \cdot 26 + 1$ ✓
- Zatem $a^{-1} = 5$

Funkcja deszyfrująca: $d(y) = 5(y - 19) \mod 26$

Odszyfrowany tekst:

```
ZJD KTTU EJ IPVN EQT X XKU PE LDFE WT VJIGPLT BPEQ EQT MTKREQ JS EQT XMIQXWTE
```

---

#### **Próba 2:** Założenie E→C, T→Z

Dane:

- $p_1 = E = 4$, $c_1 = C = 2$
- $p_2 = T = 19$, $c_2 = Z = 25$

Obliczenia:

- $\Delta x = (4 - 19) \mod 26 = -15 \mod 26 = 11$
- $\Delta y = (2 - 25) \mod 26 = -23 \mod 26 = 3$
- $\Delta x^{-1} = 11^{-1} \mod 26 = 19$ (jak wyżej)

Obliczamy $a$:

- $a = (3 \cdot 19) \mod 26 = 57 \mod 26 = 5$

Sprawdzenie: $\gcd(5, 26) = 1$ ✓

Obliczamy $b$:

- $b = (2 - 5 \cdot 4) \mod 26 = (2 - 20) \mod 26 = -18 \mod 26 = 8$

**Wynik:** $a = 5$, $b = 8$

**Deszyfrowanie:**
Potrzebujemy $a^{-1} = 5^{-1} \mod 26$:

- $5 \cdot 21 = 105 = 4 \cdot 26 + 1$ ✓
- Zatem $a^{-1} = 21$

Funkcja deszyfrująca: $d(y) = 21(y - 8) \mod 26$

Odszyfrowany tekst:

```
YOU NEED TO PICK THE A AND IT MUST BE COPRIME WITH THE LENGTH OF THE ALPHABET
```

---

### Analiza innych możliwości

Gdybyśmy spróbowali inne kombinacje z najczęstszymi literami (Z, C, R):

### Próba: R→E, Z→T

Dane:

- $p_1 = E = 4$, $c_1 = R = 17$
- $p_2 = T = 19$, $c_2 = Z = 25$

Obliczenia:

- $\Delta x = 11$, $\Delta y = 8$
- $a = (8 \cdot 19) \mod 26 = 22$

Sprawdzenie: $\gcd(22, 26) = 2 \neq 1$

**Klucz odrzucony** - funkcja nie jest odwracalna.

---

## 4. Dowód matematyczny - odwracalność funkcji

### Twierdzenie

Funkcja $e(x) = (ax + b) \mod 26$ jest różnowartościowa wtedy i tylko wtedy, gdy $\gcd(a, 26) = 1$.

### Dowód

Aby szyfrowanie za pomocą funkcji afinicznej $e(x) = (ax + b) \pmod{26}$ było odwracalne, kluczowe jest, aby funkcja ta była **różnowartościowa**. Oznacza to, że dwóm różnym literom tekstu jawnego (reprezentowanym przez liczby $x_1$ i $x_2$) muszą odpowiadać dwie różne litery szyfrogramu (reprezentowane przez $y_1$ i $y_2$).

#### Krok 1: Założenie różnowartościowości

Załóżmy, że nasza funkcja szyfrująca jest różnowartościowa. Oznacza to, że jeśli weźmiemy dwie różne wartości wejściowe $x_1$ i $x_2$ (gdzie $x_1 \neq x_2$), to wartości wyjściowe również muszą być różne:

$$e(x_1) \neq e(x_2)$$

Podstawiając do wzoru funkcji, otrzymujemy:

$$(ax_1 + b) \pmod{26} \neq (ax_2 + b) \pmod{26}$$

Aby uprościć, możemy pominąć `mod 26` na chwilę, pamiętając, że wszystkie operacje odbywają się w pierścieniu reszt modulo 26.

$$ax_1 + b \neq ax_2 + b$$

Redukując wyraz `b` po obu stronach, dochodzimy do:

$$ax_1 \neq ax_2$$

Co można zapisać jako:

$$a(x_1 - x_2) \neq 0$$

Wracając do notacji modulo:

$$a(x_1 - x_2) \not\equiv 0 \pmod{26}$$

Ponieważ założyliśmy, że $x_1 \neq x_2$, różnica $(x_1 - x_2)$ jest liczbą całkowitą niezerową z przedziału $[-25, 25]$. Powyższa nierówność oznacza, że iloczyn `a` i dowolnej niezerowej wartości $(x_1 - x_2)$ nie może być wielokrotnością 26.

#### Krok 2: Powiązanie z NWD(a, 26) = 1

Rozważmy, co by się stało, gdyby największy wspólny dzielnik `a` i 26 był większy niż 1, czyli $NWD(a, 26) = d > 1$.

Skoro $d$ jest dzielnikiem obu liczb, możemy zapisać:

- $a = d \cdot k_1$
- $26 = d \cdot k_2$

Gdzie $k_1$ i $k_2$ są liczbami całkowitymi, a $k_2 < 26$.

Teraz, jeśli za $(x_1 - x_2)$ podstawimy wartość $k_2$, otrzymamy:

$$a \cdot k_2 = (d \cdot k_1) \cdot k_2 = k_1 \cdot (d \cdot k_2) = k_1 \cdot 26$$

Oznacza to, że:

$$a \cdot k_2 \equiv 0 \pmod{26}$$

Znaleźliśmy więc niezerową wartość $(x_1 - x_2) = k_2$, dla której $a(x_1 - x_2)$ jest równe zeru modulo 26. To jest sprzeczne z naszym warunkiem różnowartościowości $a(x_1 - x_2) \not\equiv 0 \pmod{26}$.

Dlatego, aby funkcja była różnowartościowa, **NWD(a, 26) nie może być większe niż 1**. Ponieważ NWD jest zawsze liczbą dodatnią, jedyną możliwością jest:

$$NWD(a, 26) = 1$$

Oznacza to, że `a` i 26 muszą być liczbami **względnie pierwszymi**. Tylko wtedy istnieje element odwrotny do `a` modulo 26, co gwarantuje możliwość odwrócenia szyfrowania (deszyfracji).

## 5. Wpływ długości szyfrogramu na kryptoanalizę

### Dyskusja teoretyczna

#### Krótki szyfrogram

Statystyki dotyczące występowania liter są niewiarygodne i zaburzone przez losowość. Najczęściej występująca litera w krótkim szyfrogramie wcale nie musi odpowiadać najczęstszej literze w języku angielskim (E). To zmusza do testowania wielu błędnych hipotez (np. zakładania, że najczęstsza litera to T, A, itd.), co znacznie wydłuża proces łamania szyfru.

#### Długi szyfrogram

Zgodnie z prawem wielkich liczb, rozkład częstotliwości liter w długim tekście stabilizuje się i coraz dokładniej odzwierciedla statystyczny model języka. Najczęstsza litera w szyfrogramie z ogromnym prawdopodobieństwem odpowiada literze 'E' w tekście jawnym, a druga najczęstsza – literze 'T'. Daje to szybkie obliczenie klucza (a, b), co drastycznie skraca czas kryptoanalizy.
