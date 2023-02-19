# gzip.go

Largely derived from http://www.infinitepartitions.com/art001.html

## Building blocks / techniques

- Huffman encoding
- RLE
- un-ordered lay out of alphabets, putting rarely-to-use alphabet at the back
- Dynamic length integer values (used for encoding back-pointer copy length & distance)
- Values need not to start from zero if vertain minimum value is guaranteed

Useful ideas:

- laid out the process flow as a state diagram

Presentation flow:
- build a schema flow diagram. But instead of explaining everything directly, explain parts of it first, then slowly
build up to get the final result
  
## Technique 1: Huffman Encoding

It's one instance of optimal encoding system.

What is encoding? Encoding is mapping from **symbols** to another **symbols**. 
In this case we will map the symbol to binary values, which we will simply call as **codes**.
It is important that the encoded result must be decode-able. Some (invalid) mapping can cause the encoded data to be undecode-able because of ambiguity.

### Basic example

One example of encoding, is let say, encode the 4 shapes of a card set, as such:

| Shape     | Code |
|-----------|------|
| `club`    | 00   |
| `heart`   | 01   |
| `spade`   | 10   |
| `diamond` | 11   |

`club`, `heart`, `spade`, and `diamond` into some binary values, for example, `00`, `01`, `10`, `11`. In this way, a string of codes like `01101111` can be decoded as `heart`, `spade`, `diamond`, `diamond`.

In this particular encoding, we are using `2 bits` for each symbol, hence the **average length** of the codes are `2 bits` as well.

### Encoding Performance

In most cases, the main goal of an encoding system is to **reduce space**. That is equivalent to minimizing the **average length**. An example of a very redundant encoding system is to encode the symbols as ASCII characters, simply like mapping:

| Shape     | Code      |
|-----------|-----------|
| `club`    | `club`    |
| `heart`   | `heart`   |
| `spade`   | `spade`   |
| `diamond` | `diamond` |

In this case, the **average length** is `(4 + 5 + 5 + 7)bytes/4 = 5.25 bytes = 42 bits`. This is 21 times longer than our previous 2-bits encoding.

Can we do better than 2-bits? How about this following encoding?

### Invalid encoding and lower-bound of information

| Shape     | Code |
|-----------|------|
| `club`    | 0    |
| `heart`   | 1    |
| `spade`   | 00   |
| `diamond` | 01   |

Assuming that all shape has equal frequency of appearance, this encoding has an average length of 1.5 bits.

However, there is an issue with this encoding. When we encounter the code `00`, we don't know if it should mean 2 `club`s or 1 `spade`. So we call this invalid encoding.

So seems like the `2 bits` encoding is the best one for now. But how do we know or how do we prove it? We can have another valid encoding that is not of equal length, or which we can call "variable length encoding", for example:

| Shape     | Code |
|-----------|------|
| `club`    | 1    |
| `heart`   | 01   |
| `spade`   | 001  |
| `diamond` | 0001 |

This doesn't have ambiguity in the decoding state, but the average bit length is 2.5.

But all this while we only talk in the assumption of the frequency being equal. Assume we have a different card sets that has different frequency of shapes:

| Shape     | Frequency |
|-----------|-----------|
| `club`    | 0.5       |
| `heart`   | 0.3       |
| `spade`   | 0.1       |
| `diamond` | 0.1       |

In this case, to measure the encoding performance, we cannot simply take the average bit length, but rather we need to weight it against the frequency.

| Shape     | Frequency | Fixed Length Enc | Contribution |
|-----------|-----------|------------------|--------------|
| `club`    | 0.5       | `00`             | 1 bit        |
| `heart`   | 0.3       | `01`             | 0.6 bit      |
| `spade`   | 0.1       | `10`             | 0.2 bit      |
| `diamond` | 0.1       | `11`             | 0.2 bit      |

*Because we are talking in terms of expectation / frequency, we can use non-integer value for "bit count" despite bit being a discrete measure. 

Here, the fixed-length 2-bit encoding uses about (1 + 0.6 + 0.2 + 0.2) = 2 bits per symbol (following the symbol frequency/distribution). This is not very surprising because regardless of the frequency distribution, each symbol contributes the same amount.

Let's compare it with our variable length encoding:

| Shape     | Frequency | Variable Length Enc | Contribution |
|-----------|-----------|---------------------|--------------|
| `club`    | 0.5       | `1`                 | 0.5 bit      |
| `heart`   | 0.3       | `01`                | 0.6 bit      |
| `spade`   | 0.1       | `001`               | 0.3 bit      |
| `diamond` | 0.1       | `0001`              | 0.4 bit      |

Here, the sum of contribution is 1.8 bits.

In fact, we can make the encoding shorter:

| Shape     | Frequency | Variable Length Enc | Contribution |
|-----------|-----------|---------------------|--------------|
| `club`    | 0.5       | `1`                 | 0.5 bit      |
| `heart`   | 0.3       | `01`                | 0.6 bit      |
| `spade`   | 0.1       | `001`               | 0.3 bit      |
| `diamond` | 0.1       | `000`               | 0.3 bit      |

This encoding might be slightly less intuitive, but consider streams like `000010010001` can be decoded easily as:

- `000`
- `01`
- `001`
- `000`
- `1`

This encoding has the total contribution of 1.7 bits per symbol. It can be proven that this encoding is the most optimal one for such frequency distribution. 

### Huffman Algorithm

Huffman Algorithm is an algorithm to generate huffman encoding. Huffman Algorithm works on the distribution of the symbols and will generate an optimal* encoding for the set of symbols. 

*this method is only optimal for if we want to encode one symbol at a time (can be proven). If we remove that constraint and we can encode strings of symbols instead, we can do better compression (lower average bits per symbol).

Algorithm:
0. Given `n` symbols with probability/frequency of `p_0`, `p_1`, ..., `p_n`
1. Build subtree using 2 symbols with the lowest `p_i`.
2. At each step, choose 2 symbols/subtrees with the lowest total probability/frequency, combine to form new subtree
3. Result: optimal tree built with bottom-up
4. Apply `0` to left edge and `1` to right edge (or vice versa, doesn't matter)

Example:

| Shape | Frequency |
|-------|-----------|
| `A`   | 1/3       |
| `B`   | 1/2       |
| `C`   | 1/12      |
| `D`   | 1/12      |

![img_1.png](img_1.png)

Thus, we will get encoding of:

| Shape | Encoding |
|-------|----------|
| `A`   | `11`     |
| `B`   | `0`      |
| `C`   | `100`    |
| `D`   | `101`    |

------------

Read more: https://computationstructures.org/lectures/info/info.html (Entropy, Huffman Algorithm)

## Technique 2: Run-length Encoding

RLE is a fairly simple compression mechanism for data that have a lot of consecutive values, for example:

```
WWWWWWWWWWWWBWWWWWWWWWWWWBBBWWWWWWWWWWWWWWWWWWWWWWWWBWWWWWWWWWWWWWW
```

Can be encoded with RLE as:

```
12W1B12W3B24W1B14W
```

More care will be needed in the actual encoding system to differentiate the "counter" from the "symbol"

One example of such mechanism could be:

(this also optimize for if we have runs of literals, we don't waste on writing the header overhead)

```
buf = []
c = reader.next()
while c != EOF:
    if MSB(c) == 0:  # Most Significant Bit
        literal_count = LSB(c, 7)  # Least-Significant Bit
        for i in range(literal_count):
            buf.append(reader.next())
    if MSB(c) == 0:
        repeat_count = LSB(c, 7)
        val = reader.next()
        for i in range(repeat_count):
            buf.append(reader.next())
    c = reader.next()
```

Such that strings like  `AAAABCDEFGGGGGGGGGGGGGGGGGGGG` can be encoded as:

```
[
    0x84 (set the MSB to be 1, and the 7-bit LSB is 4),
    'A',
    0x05 (MSB is 0, there will be 5 literals),
    'B',
    'C',
    'D',
    'E',
    'F',
    0x94 (set the MSB to be 1, add 0x14 (20 repetitions) of the following:),
    'G',
]
```

with a total of 10 bytes.

The actual encoding system can differ, this is just one example implementation of RLE.

## LZ77

LZ77 is a data compression algorithm published by Lempel and Ziv in 1988 (not to be confused with LZ78, which is similar but slightly different). In principle, LZ77 maintains a sliding window where we can reference some earlier data in the sliding window.

One implementation could be to have triplets of data: `(back pointer, copy-length, next_byte)` as such:

![img_2.png](img_2.png)

Notice the triplet like `(3, 4, b)`, the string `aaca` is a self-reference. (the last `a` in the blue box refers to first `a` in the blue box). This always happens when `copy-length` is bigger than `back-pointer`.

The actual LZ77 implementation used in GZIP slightly more complicated for more efficient use of space.

