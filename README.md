# Weaver

1) Build with just `make`

2) Run with:

`$~ sudo ./dist/weaver <target_binary> <target_symbol>`

Now whenever the binary is run and the symbol is executed, weaver will dump the first 50 bytes of the stack, and all registers.

This will only work on x86_64!

For example:

`sudo ./dist/weaver ./dist/tester main.test_combined_int`

Running `./dist/tester` should give you:

```
{
 "MemoryStack": [
  224,
  128,
  19,
  0,
  192,
  0,
  0,
  0,
  50,
  0,
  0,
  0,
  0,
  0,
  0,
  0,
  224,
  0,
  0,
  0,
  0,
  0,
  0,
  0,
  224,
  128,
  19,
  0,
  192,
  0,
  0,
  0,
  224,
  0,
  0,
  0,
  0,
  0,
  0,
  0,
  224,
  0,
  0,
  0,
  0,
  0,
  0,
  0,
  208,
  31
 ],
 "Registers": {
  "R15": -1,
  "R14": 824633729440,
  "R13": 0,
  "R12": 824634025680,
  "Bp": 824634025840,
  "Bx": 3,
  "R11": 824633827424,
  "R10": 139676345013048,
  "R9": 1,
  "R8": 1,
  "Ax": 2,
  "Cx": 3,
  "Dx": 4616376,
  "Si": 1,
  "Di": 824633827440,
  "Orig_ax": -1,
  "Ip": 4543680,
  "Cs": 51,
  "Flags": 518,
  "Sp": 824634025768,
  "Ss": 43
 }
}
```
