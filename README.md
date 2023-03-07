# gzip.go

Largely derived from http://www.infinitepartitions.com/art001.html
  
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

## Technique 2: LZ77

LZ77 is a data compression algorithm published by Lempel and Ziv in 1977 (not to be confused with LZ78, which is similar but slightly different). In principle, LZ77 maintains a sliding window where we can reference some earlier data in the sliding window.

One implementation could be to have triplets of data: `(back pointer, copy-length, next_byte)` as such:

![img_2.png](img_2.png)

Notice the triplet like `(3, 4, b)`, the string `aaca` is a self-reference. (the last `a` in the blue box refers to first `a` in the blue box). This always happens when `copy-length` is bigger than `back-pointer`.

The actual LZ77 implementation used in GZIP slightly more complicated for more efficient use of space. But we can discuss it later.



## Step-by-step

Consider we have this body of text (char count 1408 chars)

```
So shaken as we are, so wan with care,
Find we a time for frighted peace to pant,
And breathe short-winded accents of new broils
To be commenced in strands afar remote.
No more the thirsty entrance of this soil
Shall daub her lips with her own children's blood;
Nor more shall trenching war channel her fields,
Nor bruise her flowerets with the armed hoofs
Of hostile paces: those opposed eyes,
Which, like the meteors of a troubled heaven,
All of one nature, of one substance bred,
Did lately meet in the intestine shock
And furious close of civil butchery
Shall now, in mutual well-beseeming ranks,
March all one way and be no more opposed
Against acquaintance, kindred and allies:
The edge of war, like an ill-sheathed knife,
No more shall cut his master. Therefore, friends,
As far as to the sepulchre of Christ,
Whose soldier now, under whose blessed cross
We are impressed and engaged to fight,
Forthwith a power of English shall we levy;
Whose arms were moulded in their mothers' womb
To chase these pagans in those holy fields
Over whose acres walk'd those blessed feet
Which fourteen hundred years ago were nail'd
For our advantage on the bitter cross.
But this our purpose now is twelve month old,
And bootless 'tis to tell you we will go:
Therefore we meet not now. Then let me hear
Of you, my gentle cousin Westmoreland,
What yesternight our council did decree
In forwarding this dear expedience.
```

After applying LZ77 (1988 chars): // todo add stats for back of the envelop maths
```
So shaken as we are, so wan with c<16,4>
Find<12,5> time for frighted peace to pant,
A<41,3>breathe<2,3>ort-w<40,3><64,3>accents of new<85,3>oils
To be commenc<104,3>in strands afar remote.
No mor<71,3>h<175,4>irsty <110,3><150,3><70,3><115,3><181,3>s<20,3>il
Shall daub her lips<27,6><222,4>own children's blood;<168,3>r<171,6>s<212,5>t<249,3><244,3>ng<23,3>r<243,3>annel<235,5>fields,<261,5>bruise<298,6>loweret<229,7><177,4>arm<142,3>hoofs
Of<350,3>stile<75,3>ces:<340,3>o<319,3>opp<377,3>d eye<308,3>Which,<225,3>k<175,6>meteor<113,5><47,3>roubl<348,4>eaven<80,3><274,3><419,3>one natu<17,4><445,7>subst<193,5><86,3>d,
Did lately<410,3>et<144,4><407,4>inte<362,3>n<92,5>ck<81,5>furious cl<377,5>f civil butc<322,3>y<210,7>now,<498,4>mutual<43,3>ll-beseem<283,4><192,3>k<392,3>March <560,4><463,4>way a<83,4><450,3><170,7><381,7>
Against<106,3>qu<644,3><471,5>, k<101,3>r<104,4><620,3><607,3>i<371,3>
T<503,3>edg<538,5><287,3><400,7><25,3>i<581,3>sh<88,5>d knife,<168,9><271,6>cut <202,4>master. <684,3>re<54,3><661,3><58,3>en<307,4>As <157,4><10,3><73,3><90,5>epulchr<691,5>Christ<393,4><536,4>soldi<323,3><564,5>u<102,3>r w<818,5><428,3>s<385,4>cross
W<14,5> impr<850,6><672,4>engag<876,3><789,3>f<60,4><37,3><96,3>h<336,5>a p<328,4><805,4>English<736,7><44,3>levy;<816,7><345,3><11,4><866,3>mould<142,6><792,3>i<264,4><972,3>rs' womb<128,4><291,3>s<405,5>s<366,4>gans<968,6><947,4>ho<491,3><303,6>
Ov<839,9>ac<872,3><695,3>lk'd<1016,7><848,8>f<495,3><394,6><53,3>urte<7,3>hu<666,6>yea<416,3>ago<955,6>nail'd<900,4> <1085,3> adv<77,3>a<690,4><500,6>bit<754,3><855,6>.
B<744,3><201,5><1127,4>pur<636,4><830,4> <1168,3>t<579,3>v<959,4>n<908,3><824,3><80,7>oot<1066,4> 'ti<787,6><580,3> you<935,4>w<709,3> go<682,5><762,6><1237,4><494,5>no<1266,4>w<757,5>n<938,3>t<1262,3><432,4>r<356,4><1234,3>, my g<189,3>l<133,4>us<1014,3>W<509,3><732,4>l<879,3><815,4>at<1100,3><753,4>n<895,4><1170,5><1312,3>nc<546,3>d<484,3>de<1047,3>e
In<53,4><696,3>d<590,4><1166,5>d<1290,3> expe<826,3><659,3>.
```

Note that this encoding of `<back_pointer,length>` is not the most performant one and also not valid if we have the literal `<` (or at least need to escape)