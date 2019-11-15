The Rather Large Multipart Test Suite
=====================================

The purpose of this test suite is to generate a corpus of data that can be used
to validate http servers / applications ability to parse multipart/form-data.
The present goal is not to provde/disprove "accuracy to specification", as no
specification for this protocol is ratified. The goal is to collect samples for
as many edge cases as possible, such that a robust parser can be validated to
a reasonable level of confidence.

Only a very minimal effort has been taken to try to keep this test suite small.
The goal is to attempt to collect data reliably, not (yet) to minimize the work
that implementors need to do. It is very likely that many of the resultant
samples (for example most versions of a browser from a single vendor) will be
protocol-identical. Before doing these tests, this cannot be easily proven,
particularly for closed source browsers.

Background
----------

I've been maintaining Rack for some years, and over those years I have found
that multipart is a bit of a crapshoot. Some Good Folk at w3c are trying to do
similar work (e.g. https://github.com/masinter/multipart-form-data), but not
with the same level of coverage to date. I plan to provide them with this data,
and/or resultant summaries.

Method
------

There are a lot of cases that are known to cause troubles in multipart forms,
each case will be tested in combination with other cases:
 * Different page character sets
 * Different form character sets
 * Single non-file fields
 * Single file fields
 * Multiple non-file fields
 * Multiple file fields
 * Single & multiple "multiple" file fields
 * Fields both filled in and not filled in
 * Content-Types: use known and unknown extensions
 * File names: plain names, reserved chars, 8 bit chars
 * Field names: plain names, reserved chars, 8 bit chars

There are several components (likely more, as this README may fall out of date):

 * dumper.go: this is a simple httpd that will be used for the tests. It serves
   up index.html, files from files/, files from tests/. When it recieves a POST
   to /dump/..., it will record the HTTP request data into a file named after
   the User-Agent header (URI escaped).
 * maketests.rb: this is a naive script that generates test html files in
   specific encodings, with particular form fields, etc. It also generates the
   files to be uploaded in files/.

Directories:
 * /dump/... : contains the results of browser uploads for each use case.
 * /tests/... : contains html files that should be filled out for each browser
   and submitted to produce dumps.
 * /files/... : contains files for uploading.

Results
-------

...


Copyright & License
-------------------

2014 (c) Google, Inc.
Apache 2. See LICENSE.txt.
This is NOT an official Google product.
