Repocket
========

Building / running
------------------

Rust, openssl and OSX woes:

https://stackoverflow.com/questions/34612395/openssl-crate-fails-compilation-on-mac-os-x-10-11

    $ export OPENSSL_INCLUDE_DIR=$(brew --prefix openssl)/include
    $ export OPENSSL_LIB_DIR=$(brew --prefix openssl)/lib

    $ cargo run
