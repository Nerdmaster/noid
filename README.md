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

If folks are gonna do this, they should have an option that doesn't force a
specific language on them.  i.e., they should have something that can just run
as a standalone binary, and is an **actual** micro-service (writing a Ruby
*library* and calling it a micro-service doesn't count, guys).
