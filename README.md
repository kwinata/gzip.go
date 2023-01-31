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

