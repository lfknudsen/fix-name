# FixName

Utility for fixing the names of files downloaded in Firefox (maybe the issue
exists in Chromium-based browsers as well, but I wouldn't know) where spaces are
sometimes replaced by `+`.

This CLI programme corrects that:\
`++` is replaced with ` +`.\
`+` is replaced with ` `.

It will (as the final step) trim spaces from the start and end.\
Works correctly with multi-character symbols.

## Usage

```
fixname [-r] [-n]
```

Supports the following flags:\
`-r`: Perform the action in subdirectories as well.\
`-n`: Get rid of _all_ pluses, including those that may have been intentional.

## Example outcome

Without `-n`:

`hel+lo.png` becomes `hel lo.png`

`Salt+++Pepper.mkv` becomes `Salt + Pepper.mkv`

`Remember+Google++.dat` becomes `Remember Google +.dat`

With `-n`:

`Salt+++Pepper.mkv` becomes `Salt Pepper.mkv`.

`Remember+Google++.dat` becomes `Remember Google.dat`