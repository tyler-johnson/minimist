# minimist

Simple CLI args parser.

Port of [minimist](https://github.com/substack/minimist) to golang.

This is a fork of <https://github.com/mgutz/minimist> that removes the type coersion methods and adds deep flag support.

## options

`--a            // a == true`
`--a=foo        // a == "foo"`
`--a foo        // a == "foo"`

`--no-a         // a == false`
`-a             // a == true`
`-ab            // a == true, b == true`
`-ab foo        // a == true, b == "foo"`

`--a.b          // a == { b: true }`
`--a.b foo      // a == { b: "foo" }`

# license

MIT
