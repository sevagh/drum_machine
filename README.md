Pipe [harmonixset](https://github.com/urinieto/harmonixset) beat txt files through drum_machine to listen to a clicktrack.

The file format drum_machine supports are the beat annotations here: [harmonixset/datasets/beats_and_downbeats](https://github.com/urinieto/harmonixset/tree/master/dataset/beats_and_downbeats)

The clicktrack is generated using [libmetro](https://github.com/sevagh/libmetro). The beat txt files are parsed using a handwritten tokenizer/lexer and parser, based on [this tutorial](https://blog.gopheracademy.com/advent-2014/parsers-lexers/). drum_machine is fuzz tested with the entire harmonixset dataset as a corpus, and fuzzing is done with [go-fuzz](https://github.com/dvyukov/go-fuzz) ([tutorial](https://dgraph.io/blog/post/continuous-fuzzing-with-go/)). Run fuzz test with `./fuzz_test.sh`.

### Mapping harmonix to libmetro

libmetro provides the following constructs:

* _beats_ are represented with _note_ objects, which are wrappers around a vector of floats generated from Stk. libmetro is more flexible but the C wrapper and use in drum_machine are limited to the MIDI Drum instrument, and 2 'timbres' of beats
* _bars_ are represented with _measure_ objects, which are collections of a sequence of _notes_ (can be padded with silence)
* the _tempo_ of the song is configured on the _metronome_ object, which can register a sequence of _measures_ to loop through

Although the result is not mind-blowingly exciting given the tempo and beat structure of pop songs is stable, drum_machine still a fun combination of several projects.

### Compile

Getting `drum_machine` built is not straightforward. You have to statically compile and install several C++ libraries from scratch. They must be compiled statically Go + cgo doesn't care about your configured dynamic library paths.

* https://github.com/andrewrk/libsoundio

After ensuring you have the correct dependencies (I've tested it only on Linux with ALSA, pulseaudio, and JACK libraries installed):

```
libsoundio $ mkdir -p build && cd build && cmake .. && make && sudo make install
```

* https://github.com/thestk/stk

I suggest you copy `rawwaves` to `~/rawwaves` and compile this way:

```
stk $ autoconf
stk $ ./configure
stk $ ./configure RAWWAVE_PATH='$(HOME)/rawwaves/' BUILD_STATIC='yes'
stk $ cd src && make && sudo make install

# manually install the .a file, it's not in the makefile i believe
stk $ sudo cp libstk.a /usr/local/lib/libstk.a
```

* https://github.com/sevagh/libmetro

Libmetro is a C++ library but it has a minimal C wrapper to work with cgo.

```
libmetro $ mkdir -p build && cd build && cmake .. -DBUILD_STATIC=ON && make && sudo make install
```

Finally, compile this project:

```
drum_machine $ go build .
```

cgo configuration is:

```
// #cgo LDFLAGS: /usr/local/lib/libmetro.a /usr/local/lib64/libsoundio.a /usr/local/lib/libstk.a -lstdc++ -lasound -lpulse -ljack -lm
// #include <libmetro/cmetro.h>
import "C"
```
