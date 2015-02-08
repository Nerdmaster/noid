C, noid, Go!  Go, noid, Go!
=====

**NOTE**: I have no idea if I'll ever actually use C in this project, but I
couldn't just *not* use that title.  It would be criminal.

License
-----

This software is licensed under the [Creative Commons Attribution 4.0 International](./LICENSE.txt)
license, which seems to be the most public-domain-like license I can find
that's not on shaky ground legally.

About
-----

This is my attempt to make an awful system at least a little more sensible.
[NOID](https://wiki.ucop.edu/display/Curation/NOID)s certainly have their place
in very specific situations, but overall...  I'm underwhelmed.

The idea of having a meaningless, long-term, portable identifier is certainly
well-intentioned, but I feel the NOID approach, and how I've seen it used, adds
unnecessary overhead and confusion:

- Managing such an identifier should be done at a "redirection layer" or
  something - the app should know that the data exists, and be able to find
  related data, but it shouldn't have to do a ton of work to use the identifier
  as a primary key
- I can't see that much need for such a complex templating system - simply
  minting sequential numeric values should be fine in the vast majority of
  situations, and becomes a much simpler problem to solve
  - If the values need to truly *look* unpredictable (though I don't see this
    mattering myself), a simple xor works well enough
  - If the values need to *be* unpredictable, then the NOID spec actually
    doesn't work, as its purpose is to mint the same list of values for a given
    template.  Apps end up forcing random starting points to have
    "unpredictability", but that just suggests to me that the
    sequential-plus-xor approach, with a random xor, would have been better.
    (Actually, I may just do this for the "random" ordering in this app....)
- The full spec defines all kinds of operations on noids, which, to me, seem
  oddly specific.  I'm okay with the core concept of meaningless identifiers,
  but the spec ends up defining operations on a database, blending a simple
  concept with very implementation-specific details.
  - My go-based service won't have all those database-type features, as it's a
    huge waste of effort to try and accommodate every use case for opaque IDs.
  - As an aside, I think this general "make the spec / tool do everything for
    every possible use case we can ever imagine" philosophy is exactly why we
    (libraries) end up being so incompatible with the rest of the world
    technologically.
- I'm not sold on needing to prevent transcription errors.  Removing the letter
  "L", and having an optional check-digit?  Waste of thought, unnecessary code
  complexity.  Transcription errors will still fail, just less often.  Better
  to avoid transcription by hand in the first place.
- I'm wrong about a lot of things, so maybe I'll update this as I learn more

Why
-----

*Please note*: at the time I started building this, I had *not* seen [the ndlib
noids library](https://github.com/ndlib/noids).  I'm not sure how I missed it,
but I did.  Clearly the idea of doing noids as a Go service wasn't as original
as I thought, and so the existence of my repo looks pretty ****ing lame.

I swear it isn't stolen - when I was building mine all I could find was
horrible options that required either a bunch of dependencies (Ruby libraries),
were too complex (the Perl version), or weren't finished.

I may or may not continue to improve on mine, but as I've just noticed the
ndlib app, I'm not longer thinking I'm filling a void in the library world
(while also mocking that community).  So if you need something for production,
chances are you should use theirs!

I imagine our approaches will differ in pretty significant ways, but the ndlib
github group is made up of people who are, in general, going to understand
these types of problems far better than I.  So if I do something cooler (likely
by accident), it's still probably not worth using my code over theirs.

I feel kind of stupid.

Variances
-----

This library deviates from the official spec in many ways, of course, but as
much as possible I will try to detail the changes that are significant or which
could cause major confusion.

### Digits and extended digits

In order to more easily construct random noids and perform arithmetic while
generating a noid, all digits are of a very specific bit size.  The standard
"digit" is 3 bits, while the extended digit is 5 bits.

This means a digit has only 8 possible values instead of 10, whereas an
extended digit has 32 possible values instead of the 29 the spec details.

For cases where the regular digit is used more frequently than extended digits,
this greatly diminishes the number of noids available.  In cases where the
extended digit is used twice for every regular digit ("eedeed", for instance),
the number of mintable noids is very close, but my system still has around 5%
fewer available noids.

Knowing exactly how many bits will be in use has little practical value, but is
useful for some of the internals of the system, particularly creating the
"random" noids without having to hold a huge pool of used / unused noids.  By
storing a simple sequence value and shuffling around bits based on the template
in use, noids can appear to have no obvious sequence, but there's no actual
need to store anything beyond the internal sequential value.

As an example, let's look at the template "rede".  In binary, the minimum value
would be `00000 000 00000` (0) and the maximum would be `11111 111 11111`
(8191).  The spacing is there to make it more clear that each "e" is 5 bits,
each "d" is 3.

At sequence 0, a non-random noid value would just fill all those bits with zero
and convert each chunk of bits to the appropriate digit.  In the case of zero,
the digit is always "0", giving us "000".  As the internal sequence counter
increments, so does the noid, and in a fairly predictable way: 001, 002, 003,
004, 005, 006, etc.

For a random noid, however, we can use our knowledge that there are precisely
13 bits available in order to xor and do bit swapping so that each sequential
value *appears* to be non-sequential.  If, for instance, we had a 13-bit xor
value of `1010101010101`, the sequence would look like this: p2p, p2n, p2r,
p2q, p2j, p2h, p2m, etc.  That's still not terrribly random, but it's better
than nothing.

However, by knowing that we have a set number of bits, we can define a
bit-shuffling algorithm as follows (for simplicity in following the internals,
I'm calling bit 0 the *most* significant digit, but the code doesn't actually
operate in this manner):

- Swap 12 and 3
- Swap 11 and 6
- Swap 10 and 8
- Swap 7 and 1

If we shuffle after the xor, our sequence becomes: r0q, p0q, r2q, p2q, r06,
p06, r26, p26, r0y, etc.  Though this still has some kind of pattern, it's not
immediately obvious what comes next at any given index.  And yet, we never need
anything other than a preset xor, bit swap list, and the sequence index.  By
basing the xor and bit swap on a template, we guarantee that a sequence index
plus a template are all we need to generate a given "random" noid consistently.
